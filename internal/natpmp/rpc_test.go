package natpmp

import (
	"context"
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Client_rpc(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		ctx              context.Context
		gateway          netip.Addr
		request          []byte
		responseSize     uint
		initialRetry     time.Duration
		exchanges        []udpExchange
		expectedResponse []byte
		err              error
		errMessage       string
	}{
		"gateway_ip_unspecified": {
			gateway:    netip.IPv6Unspecified(),
			request:    []byte{0, 0},
			err:        ErrGatewayIPUnspecified,
			errMessage: "gateway IP is unspecified",
		},
		"request_too_small": {
			gateway:      netip.AddrFrom4([4]byte{127, 0, 0, 1}),
			request:      []byte{0},
			initialRetry: time.Second,
			err:          ErrRequestSizeTooSmall,
			errMessage: `checking request: message size is too small: ` +
				`need at least 2 bytes and got 1 byte\(s\)`,
		},
		"write_error": {
			ctx:     context.Background(),
			gateway: netip.AddrFrom4([4]byte{127, 0, 0, 1}),
			request: []byte{0, 0},
			errMessage: `writing to connection: write udp ` +
				`127.0.0.1:[1-9][0-9]{0,4}->127.0.0.1:[1-9][0-9]{0,4}: ` +
				`i/o timeout`,
		},
		"call_error": {
			ctx:          context.Background(),
			gateway:      netip.AddrFrom4([4]byte{127, 0, 0, 1}),
			request:      []byte{0, 1},
			initialRetry: time.Millisecond,
			exchanges:    []udpExchange{{close: true}},
			err:          ErrConnectionTimeout,
			errMessage:   "connection timeout: after 1ms",
		},
		"response_too_small": {
			ctx:          context.Background(),
			gateway:      netip.AddrFrom4([4]byte{127, 0, 0, 1}),
			request:      []byte{0, 0},
			initialRetry: time.Second,
			exchanges: []udpExchange{{
				request:  []byte{0, 0},
				response: []byte{1},
			}},
			err: ErrResponseSizeTooSmall,
			errMessage: `checking response: response size is too small: ` +
				`need at least 4 bytes and got 1 byte\(s\)`,
		},
		"unexpected_response_size": {
			ctx:          context.Background(),
			gateway:      netip.AddrFrom4([4]byte{127, 0, 0, 1}),
			request:      []byte{0x0, 0x2, 0x0, 0x0, 0x0, 0x7b, 0x1, 0xc8, 0x0, 0x0, 0x4, 0xb0},
			responseSize: 5,
			initialRetry: time.Second,
			exchanges: []udpExchange{{
				request:  []byte{0x0, 0x2, 0x0, 0x0, 0x0, 0x7b, 0x1, 0xc8, 0x0, 0x0, 0x4, 0xb0},
				response: []byte{0, 1, 2, 3}, // size 4
			}},
			err: ErrResponseSizeUnexpected,
			errMessage: `checking response: response size is unexpected: ` +
				`expected 5 bytes and got 4 byte\(s\)`,
		},
		"unknown_protocol_version": {
			ctx:          context.Background(),
			gateway:      netip.AddrFrom4([4]byte{127, 0, 0, 1}),
			request:      []byte{0x0, 0x2, 0x0, 0x0, 0x0, 0x7b, 0x1, 0xc8, 0x0, 0x0, 0x4, 0xb0},
			responseSize: 16,
			initialRetry: time.Second,
			exchanges: []udpExchange{{
				request:  []byte{0x0, 0x2, 0x0, 0x0, 0x0, 0x7b, 0x1, 0xc8, 0x0, 0x0, 0x4, 0xb0},
				response: []byte{0x1, 0x82, 0x0, 0x0, 0x0, 0x14, 0x4, 0x96, 0x0, 0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			}},
			err:        ErrProtocolVersionUnknown,
			errMessage: "checking response: protocol version is unknown: 1",
		},
		"unexpected_operation_code": {
			ctx:          context.Background(),
			gateway:      netip.AddrFrom4([4]byte{127, 0, 0, 1}),
			request:      []byte{0x0, 0x2, 0x0, 0x0, 0x0, 0x7b, 0x1, 0xc8, 0x0, 0x0, 0x4, 0xb0},
			responseSize: 16,
			initialRetry: time.Second,
			exchanges: []udpExchange{{
				request:  []byte{0x0, 0x2, 0x0, 0x0, 0x0, 0x7b, 0x1, 0xc8, 0x0, 0x0, 0x4, 0xb0},
				response: []byte{0x0, 0x88, 0x0, 0x0, 0x0, 0x14, 0x4, 0x96, 0x0, 0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			}},
			err:        ErrOperationCodeUnexpected,
			errMessage: "checking response: operation code is unexpected: expected 0x82 and got 0x88",
		},
		"failure_result_code": {
			ctx:          context.Background(),
			gateway:      netip.AddrFrom4([4]byte{127, 0, 0, 1}),
			request:      []byte{0x0, 0x2, 0x0, 0x0, 0x0, 0x7b, 0x1, 0xc8, 0x0, 0x0, 0x4, 0xb0},
			responseSize: 16,
			initialRetry: time.Second,
			exchanges: []udpExchange{{
				request:  []byte{0x0, 0x2, 0x0, 0x0, 0x0, 0x7b, 0x1, 0xc8, 0x0, 0x0, 0x4, 0xb0},
				response: []byte{0x0, 0x82, 0x0, 0x11, 0x0, 0x14, 0x4, 0x96, 0x0, 0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			}},
			err:        ErrResultCodeUnknown,
			errMessage: "checking response: result code: result code is unknown: 17",
		},
		"success": {
			ctx:          context.Background(),
			gateway:      netip.AddrFrom4([4]byte{127, 0, 0, 1}),
			request:      []byte{0x0, 0x2, 0x0, 0x0, 0x0, 0x7b, 0x1, 0xc8, 0x0, 0x0, 0x4, 0xb0},
			responseSize: 16,
			initialRetry: time.Second,
			exchanges: []udpExchange{{
				request:  []byte{0x0, 0x2, 0x0, 0x0, 0x0, 0x7b, 0x1, 0xc8, 0x0, 0x0, 0x4, 0xb0},
				response: []byte{0x0, 0x82, 0x0, 0x0, 0x0, 0x0, 0x4, 0x96, 0x0, 0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			}},
			expectedResponse: []byte{0x0, 0x82, 0x0, 0x0, 0x0, 0x0, 0x4, 0x96, 0x0, 0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
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

			response, err := client.rpc(testCase.ctx, testCase.gateway,
				testCase.request, testCase.responseSize)

			if testCase.errMessage != "" {
				require.Error(t, err)
				assert.Regexp(t, testCase.errMessage, err.Error())
			} else {
				assert.ErrorIs(t, err, testCase.err)
			}
			assert.Equal(t, testCase.expectedResponse, response)
		})
	}
}
