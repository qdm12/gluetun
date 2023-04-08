package natpmp

import (
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Client_ExternalAddress(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		gateway                   netip.Addr
		initialRetry              time.Duration
		exchange                  udpExchange
		durationSinceStartOfEpoch time.Duration
		externalIPv4Address       netip.Addr
		err                       error
		errMessage                string
	}{
		"failure": {
			gateway:      netip.AddrFrom4([4]byte{127, 0, 0, 1}),
			initialRetry: time.Millisecond,
			exchange:     udpExchange{close: true},
			err:          ErrConnectionTimeout,
			errMessage:   "executing remote procedure call: connection timeout: after 1ms",
		},
		"success": {
			gateway:      netip.AddrFrom4([4]byte{127, 0, 0, 1}),
			initialRetry: time.Millisecond,
			exchange: udpExchange{
				request:  []byte{0, 0},
				response: []byte{0x0, 0x80, 0x0, 0x0, 0x0, 0x13, 0xf2, 0x4f, 0x49, 0x8c, 0x36, 0x9a},
			},
			durationSinceStartOfEpoch: time.Duration(0x13f24f) * time.Second,
			externalIPv4Address:       netip.AddrFrom4([4]byte{0x49, 0x8c, 0x36, 0x9a}),
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			remoteAddress := launchUDPServer(t, []udpExchange{testCase.exchange})

			client := Client{
				serverPort:   uint16(remoteAddress.Port),
				initialRetry: testCase.initialRetry,
				maxRetries:   1,
			}

			durationSinceStartOfEpoch, externalIPv4Address, err :=
				client.ExternalAddress(testCase.gateway)
			assert.ErrorIs(t, err, testCase.err)
			if testCase.err != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
			assert.Equal(t, testCase.durationSinceStartOfEpoch, durationSinceStartOfEpoch)
			assert.Equal(t, testCase.externalIPv4Address, externalIPv4Address)
		})
	}
}
