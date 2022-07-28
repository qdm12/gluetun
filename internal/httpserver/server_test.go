package httpserver

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:generate mockgen -destination=logger_mock_test.go -package $GOPACKAGE . Logger

func Test_New(t *testing.T) {
	t.Parallel()

	someHandler := http.NewServeMux()
	someLogger := &testLogger{}

	testCases := map[string]struct {
		settings   Settings
		expected   *Server
		errWrapped error
		errMessage string
	}{
		"empty settings": {
			errWrapped: ErrHandlerIsNotSet,
			errMessage: "http server settings validation failed: HTTP handler cannot be left unset",
		},
		"filled settings": {
			settings: Settings{
				Address:           ":8001",
				Handler:           someHandler,
				Logger:            someLogger,
				ReadHeaderTimeout: time.Second,
				ReadTimeout:       time.Second,
				ShutdownTimeout:   time.Second,
			},
			expected: &Server{
				address:           ":8001",
				handler:           someHandler,
				logger:            someLogger,
				readHeaderTimeout: time.Second,
				readTimeout:       time.Second,
				shutdownTimeout:   time.Second,
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			server, err := New(testCase.settings)
			assert.ErrorIs(t, err, testCase.errWrapped)
			if testCase.errWrapped != nil {
				require.EqualError(t, err, testCase.errMessage)
			}

			if server != nil {
				assert.NotNil(t, server.addressSet)
				server.addressSet = nil
			}

			assert.Equal(t, testCase.expected, server)
		})
	}
}
