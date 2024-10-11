package updater

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type httpExchange struct {
	requestURL     string
	responseStatus int
	responseBody   io.ReadCloser
}

func Test_addServersFromAPI(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		hts       hostToServers
		exchanges []httpExchange
		expected  hostToServers
		err       error
	}{
		"fetch API error": {
			exchanges: []httpExchange{{
				requestURL:     "https://api.surfshark.com/v4/server/clusters/generic",
				responseStatus: http.StatusNoContent,
			}},
			err: errors.New("HTTP status code not OK: 204 No Content"),
		},
		"success": {
			hts: hostToServers{
				"existinghost": []models.Server{{Hostname: "existinghost"}},
			},
			exchanges: []httpExchange{{
				requestURL:     "https://api.surfshark.com/v4/server/clusters/generic",
				responseStatus: http.StatusOK,
				responseBody: io.NopCloser(strings.NewReader(`[
				{"connectionName":"host1","region":"region1","country":"country1","location":"location1"},
				{"connectionName":"host1","region":"region1","country":"country1","location":"location1","pubkey":"pubKeyValue"},
				{"connectionName":"host2","region":"region2","country":"country1","location":"location2"}
			]`)),
			}, {
				requestURL:     "https://api.surfshark.com/v4/server/clusters/double",
				responseStatus: http.StatusOK,
				responseBody:   io.NopCloser(strings.NewReader(`[]`)),
			}, {
				requestURL:     "https://api.surfshark.com/v4/server/clusters/static",
				responseStatus: http.StatusOK,
				responseBody:   io.NopCloser(strings.NewReader(`[]`)),
			}, {
				requestURL:     "https://api.surfshark.com/v4/server/clusters/obfuscated",
				responseStatus: http.StatusOK,
				responseBody:   io.NopCloser(strings.NewReader(`[]`)),
			}},
			expected: map[string][]models.Server{
				"existinghost": {{Hostname: "existinghost"}},
				"host1": {{
					VPN:      vpn.OpenVPN,
					Region:   "region1",
					Country:  "country1",
					City:     "location1",
					Hostname: "host1",
					TCP:      true,
					UDP:      true,
				}, {
					VPN:      vpn.Wireguard,
					Region:   "region1",
					Country:  "country1",
					City:     "location1",
					Hostname: "host1",
					WgPubKey: "pubKeyValue",
				}},
				"host2": {{
					VPN:      vpn.OpenVPN,
					Region:   "region2",
					Country:  "country1",
					City:     "location2",
					Hostname: "host2",
					TCP:      true,
					UDP:      true,
				}},
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			currentExchangeIndex := 0

			client := &http.Client{
				Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					assert.Equal(t, http.MethodGet, r.Method)
					exchange := testCase.exchanges[currentExchangeIndex]
					currentExchangeIndex++
					assert.Equal(t, exchange.requestURL, r.URL.String())
					return &http.Response{
						StatusCode: exchange.responseStatus,
						Status:     http.StatusText(exchange.responseStatus),
						Body:       exchange.responseBody,
					}, nil
				}),
			}

			err := addServersFromAPI(ctx, client, testCase.hts)

			assert.Equal(t, testCase.expected, testCase.hts)
			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_fetchAPI(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		exchanges []httpExchange
		data      []serverData
		err       error
	}{
		"http response status not ok": {
			exchanges: []httpExchange{{
				requestURL:     "https://api.surfshark.com/v4/server/clusters/generic",
				responseStatus: http.StatusNoContent,
			}},
			err: errors.New("HTTP status code not OK: 204 No Content"),
		},
		"nil body": {
			exchanges: []httpExchange{{
				requestURL:     "https://api.surfshark.com/v4/server/clusters/generic",
				responseStatus: http.StatusOK,
			}},
			err: errors.New("decoding response body: EOF"),
		},
		"no server": {
			exchanges: []httpExchange{{
				requestURL:     "https://api.surfshark.com/v4/server/clusters/generic",
				responseStatus: http.StatusOK,
				responseBody:   io.NopCloser(strings.NewReader(`[]`)),
			}, {
				requestURL:     "https://api.surfshark.com/v4/server/clusters/double",
				responseStatus: http.StatusOK,
				responseBody:   io.NopCloser(strings.NewReader(`[]`)),
			}, {
				requestURL:     "https://api.surfshark.com/v4/server/clusters/static",
				responseStatus: http.StatusOK,
				responseBody:   io.NopCloser(strings.NewReader(`[]`)),
			}, {
				requestURL:     "https://api.surfshark.com/v4/server/clusters/obfuscated",
				responseStatus: http.StatusOK,
				responseBody:   io.NopCloser(strings.NewReader(`[]`)),
			}},
		},
		"success": {
			exchanges: []httpExchange{{
				requestURL:     "https://api.surfshark.com/v4/server/clusters/generic",
				responseStatus: http.StatusOK,
				responseBody: io.NopCloser(strings.NewReader(`[
					{"connectionName":"host1","region":"region1","country":"country1","location":"location1"},
					{"connectionName":"host2","region":"region2","country":"country1","location":"location2"}
				]`)),
			}, {
				requestURL:     "https://api.surfshark.com/v4/server/clusters/double",
				responseStatus: http.StatusOK,
				responseBody:   io.NopCloser(strings.NewReader(`[]`)),
			}, {
				requestURL:     "https://api.surfshark.com/v4/server/clusters/static",
				responseStatus: http.StatusOK,
				responseBody:   io.NopCloser(strings.NewReader(`[]`)),
			}, {
				requestURL:     "https://api.surfshark.com/v4/server/clusters/obfuscated",
				responseStatus: http.StatusOK,
				responseBody:   io.NopCloser(strings.NewReader(`[]`)),
			}},
			data: []serverData{
				{
					Region:   "region1",
					Country:  "country1",
					Location: "location1",
					Host:     "host1",
				},
				{
					Region:   "region2",
					Country:  "country1",
					Location: "location2",
					Host:     "host2",
				},
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			currentExchangeIndex := 0

			client := &http.Client{
				Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					assert.Equal(t, http.MethodGet, r.Method)
					exchange := testCase.exchanges[currentExchangeIndex]
					currentExchangeIndex++
					assert.Equal(t, exchange.requestURL, r.URL.String())
					return &http.Response{
						StatusCode: exchange.responseStatus,
						Status:     http.StatusText(exchange.responseStatus),
						Body:       exchange.responseBody,
					}, nil
				}),
			}

			data, err := fetchAPI(ctx, client)

			assert.Equal(t, testCase.data, data)
			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
