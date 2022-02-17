package updater

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
)

type roundTripFunc func(r *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func Test_fetchServers(t *testing.T) {
	t.Parallel()

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	testCases := map[string]struct {
		ctx            context.Context
		responseStatus int
		responseBody   io.ReadCloser
		servers        []models.Server
		errWrapped     error
		errMessage     string
	}{
		"context canceled": {
			ctx:        canceledCtx,
			errWrapped: context.Canceled,
			errMessage: `cannot fetch HTML code: Get "https://www.vpnsecure.me/vpn-locations/": context canceled`,
		},
		"success": {
			ctx:            context.Background(),
			responseStatus: http.StatusOK,
			responseBody: ioutil.NopCloser(strings.NewReader(`
			<dl class="grid__i">
				<dt>host<span>
				<div><svg></svg></div>
				<div><span>City:</span> <strong>City</strong></div>
				<div><span>Region:</span> <strong>Region</strong></div>
				<div><span>Premium:</span> <strong>YES</strong></div>
			</dl>
			`)),
			servers: []models.Server{
				{
					Hostname: "host.isponeder.com",
					City:     "City",
					Region:   "Region",
					Premium:  true,
				},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := &http.Client{
				Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					assert.Equal(t, http.MethodGet, r.Method)
					assert.Equal(t, r.URL.String(), "https://www.vpnsecure.me/vpn-locations/")

					ctxErr := r.Context().Err()
					if ctxErr != nil {
						return nil, ctxErr
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     http.StatusText(testCase.responseStatus),
						Body:       testCase.responseBody,
					}, nil
				}),
			}

			servers, err := fetchServers(testCase.ctx, client)

			assert.ErrorIs(t, err, testCase.errWrapped)
			if testCase.errWrapped != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
			assert.Equal(t, testCase.servers, servers)
		})
	}
}

func Test_fetchHTML(t *testing.T) {
	t.Parallel()

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	testCases := map[string]struct {
		ctx            context.Context
		responseStatus int
		responseBody   io.ReadCloser
		data           []byte
		errWrapped     error
		errMessage     string
	}{
		"context canceled": {
			ctx:        canceledCtx,
			errWrapped: context.Canceled,
			errMessage: `Get "https://www.vpnsecure.me/vpn-locations/": context canceled`,
		},
		"response status not ok": {
			ctx:            context.Background(),
			responseStatus: http.StatusNotFound,
			errWrapped:     ErrHTTPStatusCode,
			errMessage:     `HTTP status code is not OK: 404 Not Found`,
		},
		"success": {
			ctx:            context.Background(),
			responseStatus: http.StatusOK,
			responseBody:   ioutil.NopCloser(strings.NewReader("some body")),
			data:           []byte("some body"),
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := &http.Client{
				Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					assert.Equal(t, http.MethodGet, r.Method)
					assert.Equal(t, r.URL.String(), "https://www.vpnsecure.me/vpn-locations/")

					ctxErr := r.Context().Err()
					if ctxErr != nil {
						return nil, ctxErr
					}

					return &http.Response{
						StatusCode: testCase.responseStatus,
						Status:     http.StatusText(testCase.responseStatus),
						Body:       testCase.responseBody,
					}, nil
				}),
			}

			data, err := fetchHTML(testCase.ctx, client)

			assert.ErrorIs(t, err, testCase.errWrapped)
			if testCase.errWrapped != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
			assert.Equal(t, testCase.data, data)
		})
	}
}

func Test_parseHTML(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		html    string
		servers []models.Server
	}{
		"empty html": {},
		"html without blocks": {
			html: "some html",
		},
		"single block": {
			html: `
				<dl class="grid__i">
					<dt>host<span>
					<div><svg></svg></div>
					<div><span>City:</span> <strong>City</strong></div>
					<div><span>Region:</span> <strong>Region</strong></div>
					<div><span>Premium:</span> <strong>YES</strong></div>
				</dl>
				`,
			servers: []models.Server{
				{
					Hostname: "host.isponeder.com",
					City:     "City",
					Region:   "Region",
					Premium:  true,
				},
			},
		},
		"two block": {
			html: `

			<dl class="grid__i">
			<dt>host<span>
			<div><svg></svg></div>
			<div><span>City:</span> <strong>City</strong></div>
			<div><span>Region:</span> <strong>Region</strong></div>
			<div><span>Premium:</span> <strong>YES</strong></div>
		</dl>

			<dl class="grid__i">
		<dt>host2<span>
		<div><svg></svg></div>
		<div><span>City:</span> <strong>City 2</strong></div>
		<div><span>Region:</span> <strong>Region 2</strong></div>
		<div><span>Premium:</span> <strong>No</strong></div>
	</dl>
				`,
			servers: []models.Server{
				{
					Hostname: "host.isponeder.com",
					City:     "City",
					Region:   "Region",
					Premium:  true,
				},
				{
					Hostname: "host2.isponeder.com",
					City:     "City 2",
					Region:   "Region 2",
				},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			servers := parseHTML(testCase.html)

			assert.Equal(t, testCase.servers, servers)
		})
	}
}

func Test_extractFromHTMLBlock(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		html   string
		server models.Server
		ok     bool
	}{
		"empty html block": {},
		"host field absent": {
			html: `
				<dl class="grid__i">
					<dt><span>
					<div><svg></svg></div>
					<div><span>City:</span> <strong>City</strong></div>
					<div><span>Region:</span> <strong>Region</strong></div>
					<div><span>Premium:</span> <strong>YES</strong></div>
				</dl>
				`,
		},
		"all fields present": {
			html: `
				<dl class="grid__i">
					<dt>host<span>
					<div><svg></svg></div>
					<div><span>City:</span> <strong>City</strong></div>
					<div><span>Region:</span> <strong>Region</strong></div>
					<div><span>Premium:</span> <strong>YES</strong></div>
				</dl>
				`,
			server: models.Server{
				Hostname: "host.isponeder.com",
				City:     "City",
				Region:   "Region",
				Premium:  true,
			},
			ok: true,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			server, ok := extractFromHTMLBlock(testCase.html)

			assert.Equal(t, testCase.ok, ok)
			assert.Equal(t, testCase.server, server)
		})
	}
}
