package dns

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/golibs/logging/mock_logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_generateUnboundConf(t *testing.T) {
	t.Parallel()
	settings := settings.DNS{
		Providers:          []models.DNSProvider{constants.Cloudflare, constants.Quad9},
		AllowedHostnames:   []string{"a"},
		PrivateAddresses:   []string{"9.9.9.9"},
		BlockMalicious:     true,
		BlockSurveillance:  false,
		BlockAds:           false,
		VerbosityLevel:     2,
		ValidationLogLevel: 3,
		Caching:            true,
		IPv6:               true,
	}
	mockCtrl := gomock.NewController(t)
	ctx := context.Background()

	clientCalls := map[models.URL]int{
		constants.MaliciousBlockListIPsURL:       0,
		constants.MaliciousBlockListHostnamesURL: 0,
	}
	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			url := models.URL(r.URL.String())
			if _, ok := clientCalls[url]; !ok {
				t.Errorf("unknown URL %q", url)
				return nil, nil
			}
			clientCalls[url]++
			var body string
			switch url {
			case constants.MaliciousBlockListIPsURL:
				body = "c\nd"
			case constants.MaliciousBlockListHostnamesURL:
				body = "b\na\nc"
			default:
				t.Errorf("unknown URL %q", url)
				return nil, nil
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader(body)),
			}, nil
		}),
	}

	logger := mock_logging.NewMockLogger(mockCtrl)
	logger.EXPECT().Info("%d hostnames blocked overall", 2)
	logger.EXPECT().Info("%d IP addresses blocked overall", 3)
	lines, warnings := generateUnboundConf(ctx, settings, "nonrootuser", client, logger)
	require.Len(t, warnings, 0)
	for url, count := range clientCalls {
		assert.Equalf(t, 1, count, "for url %q", url)
	}
	const expected = `
server:
  cache-max-ttl: 9000
  cache-min-ttl: 3600
  do-ip4: yes
  do-ip6: yes
  harden-algo-downgrade: yes
  harden-below-nxdomain: yes
  harden-referral-path: yes
  hide-identity: yes
  hide-version: yes
  interface: 0.0.0.0
  key-cache-size: 16m
  key-cache-slabs: 4
  msg-cache-size: 4m
  msg-cache-slabs: 4
  num-threads: 1
  port: 53
  prefetch-key: yes
  prefetch: yes
  root-hints: "/etc/unbound/root.hints"
  rrset-cache-size: 4m
  rrset-cache-slabs: 4
  rrset-roundrobin: yes
  tls-cert-bundle: "/etc/ssl/certs/ca-certificates.crt"
  trust-anchor-file: "/etc/unbound/root.key"
  use-syslog: no
  username: "nonrootuser"
  val-log-level: 3
  verbosity: 2
  local-zone: "b" static
  local-zone: "c" static
  private-address: 9.9.9.9
  private-address: c
  private-address: d
forward-zone:
  forward-no-cache: no
  forward-tls-upstream: yes
  name: "."
  forward-addr: 1.1.1.1@853#cloudflare-dns.com
  forward-addr: 1.0.0.1@853#cloudflare-dns.com
  forward-addr: 2606:4700:4700::1111@853#cloudflare-dns.com
  forward-addr: 2606:4700:4700::1001@853#cloudflare-dns.com
  forward-addr: 9.9.9.9@853#dns.quad9.net
  forward-addr: 149.112.112.112@853#dns.quad9.net
  forward-addr: 2620:fe::fe@853#dns.quad9.net
  forward-addr: 2620:fe::9@853#dns.quad9.net`
	assert.Equal(t, expected, "\n"+strings.Join(lines, "\n"))
}

func Test_buildBlocked(t *testing.T) {
	t.Parallel()
	type blockParams struct {
		blocked   bool
		content   []byte
		clientErr error
	}
	tests := map[string]struct {
		malicious        blockParams
		ads              blockParams
		surveillance     blockParams
		allowedHostnames []string
		privateAddresses []string
		hostnamesLines   []string
		ipsLines         []string
		errsString       []string
	}{
		"none blocked": {},
		"all blocked without lists": {
			malicious: blockParams{
				blocked: true,
			},
			ads: blockParams{
				blocked: true,
			},
			surveillance: blockParams{
				blocked: true,
			},
		},
		"all blocked with lists": {
			malicious: blockParams{
				blocked: true,
				content: []byte("malicious"),
			},
			ads: blockParams{
				blocked: true,
				content: []byte("ads"),
			},
			surveillance: blockParams{
				blocked: true,
				content: []byte("surveillance"),
			},
			hostnamesLines: []string{
				"  local-zone: \"ads\" static",
				"  local-zone: \"malicious\" static",
				"  local-zone: \"surveillance\" static"},
			ipsLines: []string{
				"  private-address: ads",
				"  private-address: malicious",
				"  private-address: surveillance"},
		},
		"all blocked with allowed hostnames": {
			malicious: blockParams{
				blocked: true,
				content: []byte("malicious"),
			},
			ads: blockParams{
				blocked: true,
				content: []byte("ads"),
			},
			surveillance: blockParams{
				blocked: true,
				content: []byte("surveillance"),
			},
			allowedHostnames: []string{"ads"},
			hostnamesLines: []string{
				"  local-zone: \"malicious\" static",
				"  local-zone: \"surveillance\" static"},
			ipsLines: []string{
				"  private-address: ads",
				"  private-address: malicious",
				"  private-address: surveillance"},
		},
		"all blocked with private addresses": {
			malicious: blockParams{
				blocked: true,
				content: []byte("malicious"),
			},
			ads: blockParams{
				blocked: true,
				content: []byte("ads"),
			},
			surveillance: blockParams{
				blocked: true,
				content: []byte("surveillance"),
			},
			privateAddresses: []string{"ads", "192.100.1.5"},
			hostnamesLines: []string{
				"  local-zone: \"ads\" static",
				"  local-zone: \"malicious\" static",
				"  local-zone: \"surveillance\" static"},
			ipsLines: []string{
				"  private-address: 192.100.1.5",
				"  private-address: ads",
				"  private-address: malicious",
				"  private-address: surveillance"},
		},
		"all blocked with lists and one error": {
			malicious: blockParams{
				blocked: true,
				content: []byte("malicious"),
			},
			ads: blockParams{
				blocked:   true,
				content:   []byte("ads"),
				clientErr: fmt.Errorf("ads error"),
			},
			surveillance: blockParams{
				blocked: true,
				content: []byte("surveillance"),
			},
			hostnamesLines: []string{
				"  local-zone: \"malicious\" static",
				"  local-zone: \"surveillance\" static"},
			ipsLines: []string{
				"  private-address: malicious",
				"  private-address: surveillance"},
			errsString: []string{
				`Get "https://raw.githubusercontent.com/qdm12/files/master/ads-ips.updated": ads error`,
				`Get "https://raw.githubusercontent.com/qdm12/files/master/ads-hostnames.updated": ads error`,
			},
		},
		"all blocked with errors": {
			malicious: blockParams{
				blocked:   true,
				clientErr: fmt.Errorf("malicious"),
			},
			ads: blockParams{
				blocked:   true,
				clientErr: fmt.Errorf("ads"),
			},
			surveillance: blockParams{
				blocked:   true,
				clientErr: fmt.Errorf("surveillance"),
			},
			errsString: []string{
				`Get "https://raw.githubusercontent.com/qdm12/files/master/malicious-ips.updated": malicious`,
				`Get "https://raw.githubusercontent.com/qdm12/files/master/malicious-hostnames.updated": malicious`,
				`Get "https://raw.githubusercontent.com/qdm12/files/master/ads-ips.updated": ads`,
				`Get "https://raw.githubusercontent.com/qdm12/files/master/ads-hostnames.updated": ads`,
				`Get "https://raw.githubusercontent.com/qdm12/files/master/surveillance-ips.updated": surveillance`,
				`Get "https://raw.githubusercontent.com/qdm12/files/master/surveillance-hostnames.updated": surveillance`,
			},
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			clientCalls := map[models.URL]int{}
			if tc.malicious.blocked {
				clientCalls[constants.MaliciousBlockListIPsURL] = 0
				clientCalls[constants.MaliciousBlockListHostnamesURL] = 0
			}
			if tc.ads.blocked {
				clientCalls[constants.AdsBlockListIPsURL] = 0
				clientCalls[constants.AdsBlockListHostnamesURL] = 0
			}
			if tc.surveillance.blocked {
				clientCalls[constants.SurveillanceBlockListIPsURL] = 0
				clientCalls[constants.SurveillanceBlockListHostnamesURL] = 0
			}

			client := &http.Client{
				Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					url := models.URL(r.URL.String())
					if _, ok := clientCalls[url]; !ok {
						t.Errorf("unknown URL %q", url)
						return nil, nil
					}
					clientCalls[url]++
					var body []byte
					var err error
					switch url {
					case constants.MaliciousBlockListIPsURL, constants.MaliciousBlockListHostnamesURL:
						body = tc.malicious.content
						err = tc.malicious.clientErr
					case constants.AdsBlockListIPsURL, constants.AdsBlockListHostnamesURL:
						body = tc.ads.content
						err = tc.ads.clientErr
					case constants.SurveillanceBlockListIPsURL, constants.SurveillanceBlockListHostnamesURL:
						body = tc.surveillance.content
						err = tc.surveillance.clientErr
					default: // just in case if the test is badly written
						t.Errorf("unknown URL %q", url)
						return nil, nil
					}
					if err != nil {
						return nil, err
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(body)),
					}, nil
				}),
			}

			hostnamesLines, ipsLines, errs := buildBlocked(ctx, client,
				tc.malicious.blocked, tc.ads.blocked, tc.surveillance.blocked,
				tc.allowedHostnames, tc.privateAddresses)

			var errsString []string
			for _, err := range errs {
				errsString = append(errsString, err.Error())
			}
			assert.ElementsMatch(t, tc.errsString, errsString)
			assert.ElementsMatch(t, tc.hostnamesLines, hostnamesLines)
			assert.ElementsMatch(t, tc.ipsLines, ipsLines)

			for url, count := range clientCalls {
				assert.Equalf(t, 1, count, "for url %q", url)
			}
		})
	}
}

func Test_getList(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		content   []byte
		status    int
		clientErr error
		results   []string
		err       error
	}{
		"no result": {
			status: http.StatusOK,
		},
		"bad status": {
			status: http.StatusInternalServerError,
			err:    fmt.Errorf("bad HTTP status from irrelevant_url: Internal Server Error"),
		},
		"network error": {
			status:    http.StatusOK,
			clientErr: fmt.Errorf("error"),
			err:       fmt.Errorf(`Get "irrelevant_url": error`),
		},
		"results": {
			content: []byte("a\nb\nc\n"),
			status:  http.StatusOK,
			results: []string{"a", "b", "c"},
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			client := &http.Client{
				Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					assert.Equal(t, "irrelevant_url", r.URL.String())
					if tc.clientErr != nil {
						return nil, tc.clientErr
					}
					return &http.Response{
						StatusCode: tc.status,
						Status:     http.StatusText(tc.status),
						Body:       ioutil.NopCloser(bytes.NewReader(tc.content)),
					}, nil
				}),
			}

			results, err := getList(ctx, client, "irrelevant_url")
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.results, results)
		})
	}
}

func Test_buildBlockedHostnames(t *testing.T) {
	t.Parallel()
	type blockParams struct {
		blocked   bool
		content   []byte
		clientErr error
	}
	tests := map[string]struct {
		malicious        blockParams
		ads              blockParams
		surveillance     blockParams
		allowedHostnames []string
		lines            []string
		errsString       []string
	}{
		"nothing blocked": {},
		"only malicious blocked": {
			malicious: blockParams{
				blocked:   true,
				content:   []byte("site_a\nsite_b"),
				clientErr: nil,
			},
			lines: []string{
				"  local-zone: \"site_a\" static",
				"  local-zone: \"site_b\" static"},
		},
		"all blocked with some duplicates": {
			malicious: blockParams{
				blocked: true,
				content: []byte("site_a\nsite_b"),
			},
			ads: blockParams{
				blocked: true,
				content: []byte("site_a\nsite_c"),
			},
			surveillance: blockParams{
				blocked: true,
				content: []byte("site_c\nsite_a"),
			},
			lines: []string{
				"  local-zone: \"site_a\" static",
				"  local-zone: \"site_b\" static",
				"  local-zone: \"site_c\" static"},
		},
		"all blocked with one errored": {
			malicious: blockParams{
				blocked: true,
				content: []byte("site_a\nsite_b"),
			},
			ads: blockParams{
				blocked: true,
				content: []byte("site_a\nsite_c"),
			},
			surveillance: blockParams{
				blocked:   true,
				clientErr: fmt.Errorf("surveillance error"),
			},
			lines: []string{
				"  local-zone: \"site_a\" static",
				"  local-zone: \"site_b\" static",
				"  local-zone: \"site_c\" static"},
			errsString: []string{
				`Get "https://raw.githubusercontent.com/qdm12/files/master/surveillance-hostnames.updated": surveillance error`,
			},
		},
		"blocked with allowed hostnames": {
			malicious: blockParams{
				blocked: true,
				content: []byte("site_a\nsite_b"),
			},
			ads: blockParams{
				blocked: true,
				content: []byte("site_c\nsite_d"),
			},
			allowedHostnames: []string{"site_b", "site_c"},
			lines: []string{
				"  local-zone: \"site_a\" static",
				"  local-zone: \"site_d\" static"},
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			clientCalls := map[models.URL]int{}
			if tc.malicious.blocked {
				clientCalls[constants.MaliciousBlockListHostnamesURL] = 0
			}
			if tc.ads.blocked {
				clientCalls[constants.AdsBlockListHostnamesURL] = 0
			}
			if tc.surveillance.blocked {
				clientCalls[constants.SurveillanceBlockListHostnamesURL] = 0
			}

			client := &http.Client{
				Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					url := models.URL(r.URL.String())
					if _, ok := clientCalls[url]; !ok {
						t.Errorf("unknown URL %q", url)
						return nil, nil
					}
					clientCalls[url]++
					var body []byte
					var err error
					switch url {
					case constants.MaliciousBlockListHostnamesURL:
						body = tc.malicious.content
						err = tc.malicious.clientErr
					case constants.AdsBlockListHostnamesURL:
						body = tc.ads.content
						err = tc.ads.clientErr
					case constants.SurveillanceBlockListHostnamesURL:
						body = tc.surveillance.content
						err = tc.surveillance.clientErr
					default: // just in case if the test is badly written
						t.Errorf("unknown URL %q", url)
						return nil, nil
					}
					if err != nil {
						return nil, err
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(body)),
					}, nil
				}),
			}

			lines, errs := buildBlockedHostnames(ctx, client,
				tc.malicious.blocked, tc.ads.blocked,
				tc.surveillance.blocked, tc.allowedHostnames)

			var errsString []string
			for _, err := range errs {
				errsString = append(errsString, err.Error())
			}
			assert.ElementsMatch(t, tc.errsString, errsString)
			assert.ElementsMatch(t, tc.lines, lines)

			for url, count := range clientCalls {
				assert.Equalf(t, 1, count, "for url %q", url)
			}
		})
	}
}

func Test_buildBlockedIPs(t *testing.T) {
	t.Parallel()
	type blockParams struct {
		blocked   bool
		content   []byte
		clientErr error
	}
	tests := map[string]struct {
		malicious        blockParams
		ads              blockParams
		surveillance     blockParams
		privateAddresses []string
		lines            []string
		errsString       []string
	}{
		"nothing blocked": {},
		"only malicious blocked": {
			malicious: blockParams{
				blocked:   true,
				content:   []byte("site_a\nsite_b"),
				clientErr: nil,
			},
			lines: []string{
				"  private-address: site_a",
				"  private-address: site_b"},
		},
		"all blocked with some duplicates": {
			malicious: blockParams{
				blocked: true,
				content: []byte("site_a\nsite_b"),
			},
			ads: blockParams{
				blocked: true,
				content: []byte("site_a\nsite_c"),
			},
			surveillance: blockParams{
				blocked: true,
				content: []byte("site_c\nsite_a"),
			},
			lines: []string{
				"  private-address: site_a",
				"  private-address: site_b",
				"  private-address: site_c"},
		},
		"all blocked with one errored": {
			malicious: blockParams{
				blocked: true,
				content: []byte("site_a\nsite_b"),
			},
			ads: blockParams{
				blocked: true,
				content: []byte("site_a\nsite_c"),
			},
			surveillance: blockParams{
				blocked:   true,
				clientErr: fmt.Errorf("surveillance error"),
			},
			lines: []string{
				"  private-address: site_a",
				"  private-address: site_b",
				"  private-address: site_c"},
			errsString: []string{
				`Get "https://raw.githubusercontent.com/qdm12/files/master/surveillance-ips.updated": surveillance error`,
			},
		},
		"blocked with private addresses": {
			malicious: blockParams{
				blocked: true,
				content: []byte("site_a\nsite_b"),
			},
			ads: blockParams{
				blocked: true,
				content: []byte("site_c"),
			},
			privateAddresses: []string{"site_c", "site_d"},
			lines: []string{
				"  private-address: site_a",
				"  private-address: site_b",
				"  private-address: site_c",
				"  private-address: site_d"},
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			clientCalls := map[models.URL]int{}
			if tc.malicious.blocked {
				clientCalls[constants.MaliciousBlockListIPsURL] = 0
			}
			if tc.ads.blocked {
				clientCalls[constants.AdsBlockListIPsURL] = 0
			}
			if tc.surveillance.blocked {
				clientCalls[constants.SurveillanceBlockListIPsURL] = 0
			}

			client := &http.Client{
				Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					url := models.URL(r.URL.String())
					if _, ok := clientCalls[url]; !ok {
						t.Errorf("unknown URL %q", url)
						return nil, nil
					}
					clientCalls[url]++
					var body []byte
					var err error
					switch url {
					case constants.MaliciousBlockListIPsURL:
						body = tc.malicious.content
						err = tc.malicious.clientErr
					case constants.AdsBlockListIPsURL:
						body = tc.ads.content
						err = tc.ads.clientErr
					case constants.SurveillanceBlockListIPsURL:
						body = tc.surveillance.content
						err = tc.surveillance.clientErr
					default: // just in case if the test is badly written
						t.Errorf("unknown URL %q", url)
						return nil, nil
					}
					if err != nil {
						return nil, err
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(body)),
					}, nil
				}),
			}

			lines, errs := buildBlockedIPs(ctx, client,
				tc.malicious.blocked, tc.ads.blocked,
				tc.surveillance.blocked, tc.privateAddresses)

			var errsString []string
			for _, err := range errs {
				errsString = append(errsString, err.Error())
			}
			assert.ElementsMatch(t, tc.errsString, errsString)
			assert.ElementsMatch(t, tc.lines, lines)

			for url, count := range clientCalls {
				assert.Equalf(t, 1, count, "for url %q", url)
			}
		})
	}
}
