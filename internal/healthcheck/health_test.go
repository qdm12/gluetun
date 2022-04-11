package healthcheck

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Server_healthCheck(t *testing.T) {
	t.Parallel()

	t.Run("canceled real dialer", func(t *testing.T) {
		t.Parallel()

		dialer := &net.Dialer{}
		const address = "cloudflare.com:443"

		server := &Server{
			dialer: dialer,
			config: settings.Health{
				TargetAddress: address,
			},
		}

		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		err := server.healthCheck(canceledCtx)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "operation was canceled")
	})

	t.Run("dial localhost:0", func(t *testing.T) {
		t.Parallel()

		listener, err := net.Listen("tcp4", "localhost:0")
		require.NoError(t, err)
		t.Cleanup(func() {
			err = listener.Close()
			assert.NoError(t, err)
		})

		listeningAddress := listener.Addr()

		dialer := &net.Dialer{}
		server := &Server{
			dialer: dialer,
			config: settings.Health{
				TargetAddress: listeningAddress.String(),
			},
		}

		const timeout = 100 * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err = server.healthCheck(ctx)

		assert.NoError(t, err)
	})
}

func Test_makeAddressToDial(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		address       string
		addressToDial string
		err           error
	}{
		"host without port": {
			address:       "test.com",
			addressToDial: "test.com:443",
		},
		"host with port": {
			address:       "test.com:80",
			addressToDial: "test.com:80",
		},
		"bad address": {
			address: "test.com::",
			err:     fmt.Errorf("cannot split host and port from address: address test.com::: too many colons in address"), //nolint:lll
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			addressToDial, err := makeAddressToDial(testCase.address)

			assert.Equal(t, testCase.addressToDial, addressToDial)
			if testCase.err != nil {
				assert.EqualError(t, err, testCase.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
