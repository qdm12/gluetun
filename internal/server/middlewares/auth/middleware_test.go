package auth

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_authHandler_ServeHTTP(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		settings      Settings
		makeLogger    func(ctrl *gomock.Controller) *MockDebugLogger
		requestMethod string
		requestPath   string
		statusCode    int
		responseBody  string
	}{
		"route_has_no_role": {
			settings: Settings{
				Roles: []Role{
					{Name: "role1", Auth: AuthNone, Routes: []string{"GET /a"}},
				},
			},
			makeLogger: func(ctrl *gomock.Controller) *MockDebugLogger {
				logger := NewMockDebugLogger(ctrl)
				logger.EXPECT().Debugf("no authentication role defined for route %s", "GET /b")
				return logger
			},
			requestMethod: http.MethodGet,
			requestPath:   "/b",
			statusCode:    http.StatusUnauthorized,
			responseBody:  "Unauthorized\n",
		},
		"authorized_unprotected_by_default": {
			settings: Settings{
				Roles: []Role{
					{Name: "public", Auth: AuthNone, Routes: []string{"GET /v1/vpn/status"}},
				},
			},
			makeLogger: func(ctrl *gomock.Controller) *MockDebugLogger {
				logger := NewMockDebugLogger(ctrl)
				logger.EXPECT().Warnf("route %s is unprotected by default, "+
					"please set up authentication following the documentation at "+
					"https://github.com/qdm12/gluetun-wiki/blob/main/setup/advanced/control-server.md#authentication "+
					"since this will become no longer publicly accessible after release v3.40.",
					"GET /v1/vpn/status")
				logger.EXPECT().Debugf("access to route %s authorized for role %s",
					"GET /v1/vpn/status", "public")
				return logger
			},
			requestMethod: http.MethodGet,
			requestPath:   "/v1/vpn/status",
			statusCode:    http.StatusOK,
		},
		"authorized_none": {
			settings: Settings{
				Roles: []Role{
					{Name: "role1", Auth: AuthNone, Routes: []string{"GET /a"}},
				},
			},
			makeLogger: func(ctrl *gomock.Controller) *MockDebugLogger {
				logger := NewMockDebugLogger(ctrl)
				logger.EXPECT().Debugf("access to route %s authorized for role %s",
					"GET /a", "role1")
				return logger
			},
			requestMethod: http.MethodGet,
			requestPath:   "/a",
			statusCode:    http.StatusOK,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			var debugLogger DebugLogger
			if testCase.makeLogger != nil {
				debugLogger = testCase.makeLogger(ctrl)
			}
			middleware, err := New(testCase.settings, debugLogger)
			require.NoError(t, err)

			childHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			handler := middleware(childHandler)

			server := httptest.NewServer(handler)
			t.Cleanup(server.Close)

			client := server.Client()

			requestURL, err := url.JoinPath(server.URL, testCase.requestPath)
			require.NoError(t, err)
			request, err := http.NewRequestWithContext(context.Background(),
				testCase.requestMethod, requestURL, nil)
			require.NoError(t, err)

			response, err := client.Do(request)
			require.NoError(t, err)
			t.Cleanup(func() {
				err = response.Body.Close()
				assert.NoError(t, err)
			})

			assert.Equal(t, testCase.statusCode, response.StatusCode)
			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			assert.Equal(t, testCase.responseBody, string(body))
		})
	}
}
