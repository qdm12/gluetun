package natpmp

import (
	"context"
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Client_AddPortMapping(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		ctx                       context.Context
		gateway                   netip.Addr
		protocol                  string
		internalPort              uint16
		requestedExternalPort     uint16
		lifetime                  time.Duration
		initialRetry              time.Duration
		exchanges                 []udpExchange
		durationSinceStartOfEpoch time.Duration
		assignedInternalPort      uint16
		assignedExternalPort      uint16
		assignedLifetime          time.Duration
		err                       error
		errMessage                string
	}{
		"lifetime_too_long": {
			lifetime:   time.Duration(uint64(^uint32(0))+1) * time.Second,
			err:        ErrLifetimeTooLong,
			errMessage: "lifetime is too long: 4294967296 seconds must at most 4294967295 seconds",
		},
		"protocol_unknown": {
			lifetime:   time.Second,
			protocol:   "xyz",
			err:        ErrNetworkProtocolUnknown,
			errMessage: "network protocol is unknown: xyz",
		},
		"rpc_error": {
			ctx:                   context.Background(),
			gateway:               netip.AddrFrom4([4]byte{127, 0, 0, 1}),
			protocol:              "udp",
			internalPort:          123,
			requestedExternalPort: 456,
			lifetime:              1200 * time.Second,
			initialRetry:          time.Millisecond,
			exchanges:             []udpExchange{{close: true}},
			err:                   ErrConnectionTimeout,
			errMessage:            "executing remote procedure call: connection timeout: after 1ms",
		},
		"add_udp": {
			ctx:                   context.Background(),
			gateway:               netip.AddrFrom4([4]byte{127, 0, 0, 1}),
			protocol:              "udp",
			internalPort:          123,
			requestedExternalPort: 456,
			lifetime:              1200 * time.Second,
			initialRetry:          time.Second,
			exchanges: []udpExchange{{
				request:  []byte{0x0, 0x1, 0x0, 0x0, 0x0, 0x7b, 0x1, 0xc8, 0x0, 0x0, 0x4, 0xb0},
				response: []byte{0x0, 0x81, 0x0, 0x0, 0x0, 0x13, 0xfe, 0xff, 0x0, 0x7b, 0x1, 0xc8, 0x0, 0x0, 0x4, 0xb0},
			}},
			durationSinceStartOfEpoch: 0x13feff * time.Second,
			assignedInternalPort:      0x7b,
			assignedExternalPort:      0x1c8,
			assignedLifetime:          0x4b0 * time.Second,
		},
		"add_tcp": {
			ctx:                   context.Background(),
			gateway:               netip.AddrFrom4([4]byte{127, 0, 0, 1}),
			protocol:              "tcp",
			internalPort:          123,
			requestedExternalPort: 456,
			lifetime:              1200 * time.Second,
			initialRetry:          time.Second,
			exchanges: []udpExchange{{
				request:  []byte{0x0, 0x2, 0x0, 0x0, 0x0, 0x7b, 0x1, 0xc8, 0x0, 0x0, 0x4, 0xb0},
				response: []byte{0x0, 0x82, 0x0, 0x0, 0x0, 0x14, 0x3, 0x21, 0x0, 0x7b, 0x1, 0xc8, 0x0, 0x0, 0x4, 0xb0},
			}},
			durationSinceStartOfEpoch: 0x140321 * time.Second,
			assignedInternalPort:      0x7b,
			assignedExternalPort:      0x1c8,
			assignedLifetime:          0x4b0 * time.Second,
		},
		"remove_udp": {
			ctx:          context.Background(),
			gateway:      netip.AddrFrom4([4]byte{127, 0, 0, 1}),
			protocol:     "udp",
			internalPort: 123,
			initialRetry: time.Second,
			exchanges: []udpExchange{{
				request:  []byte{0x0, 0x1, 0x0, 0x0, 0x0, 0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
				response: []byte{0x0, 0x81, 0x0, 0x0, 0x0, 0x14, 0x3, 0xd5, 0x0, 0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			}},
			durationSinceStartOfEpoch: 0x1403d5 * time.Second,
			assignedInternalPort:      0x7b,
		},
		"remove_tcp": {
			ctx:          context.Background(),
			gateway:      netip.AddrFrom4([4]byte{127, 0, 0, 1}),
			protocol:     "tcp",
			internalPort: 123,
			initialRetry: time.Second,
			exchanges: []udpExchange{{
				request:  []byte{0x0, 0x2, 0x0, 0x0, 0x0, 0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
				response: []byte{0x0, 0x82, 0x0, 0x0, 0x0, 0x14, 0x4, 0x96, 0x0, 0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			}},
			durationSinceStartOfEpoch: 0x140496 * time.Second,
			assignedInternalPort:      0x7b,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			remoteAddress := launchUDPServer(t, testCase.exchanges)

			client := Client{
				serverPort:   uint16(remoteAddress.Port),
				initialRetry: testCase.initialRetry,
				maxRetries:   1,
			}

			durationSinceStartOfEpoch, assignedInternalPort,
				assignedExternalPort, assignedLifetime, err :=
				client.AddPortMapping(testCase.ctx, testCase.gateway,
					testCase.protocol, testCase.internalPort,
					testCase.requestedExternalPort, testCase.lifetime)

			assert.Equal(t, testCase.durationSinceStartOfEpoch, durationSinceStartOfEpoch)
			assert.Equal(t, testCase.assignedInternalPort, assignedInternalPort)
			assert.Equal(t, testCase.assignedExternalPort, assignedExternalPort)
			assert.Equal(t, testCase.assignedLifetime, assignedLifetime)
			assert.ErrorIs(t, err, testCase.err)
			if testCase.err != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
		})
	}
}
