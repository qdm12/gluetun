package natpmp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_checkRequest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		request    []byte
		err        error
		errMessage string
	}{
		"too_short": {
			request:    []byte{1},
			err:        ErrRequestSizeTooSmall,
			errMessage: "message size is too small: need at least 2 bytes and got 1 byte(s)",
		},
		"success": {
			request: []byte{0, 0},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := checkRequest(testCase.request)

			assert.ErrorIs(t, err, testCase.err)
			if testCase.err != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
		})
	}
}

func Test_checkResponse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		response              []byte
		expectedOperationCode byte
		expectedResponseSize  uint
		err                   error
		errMessage            string
	}{
		"too_short": {
			response:   []byte{1},
			err:        ErrResponseSizeTooSmall,
			errMessage: "response size is too small: need at least 4 bytes and got 1 byte(s)",
		},
		"size_mismatch": {
			response:             []byte{0, 0, 0, 0},
			expectedResponseSize: 5,
			err:                  ErrResponseSizeUnexpected,
			errMessage:           "response size is unexpected: expected 5 bytes and got 4 byte(s)",
		},
		"protocol_unknown": {
			response:             []byte{1, 0, 0, 0},
			expectedResponseSize: 4,
			err:                  ErrProtocolVersionUnknown,
			errMessage:           "protocol version is unknown: 1",
		},
		"operation_code_unexpected": {
			response:              []byte{0, 2, 0, 0},
			expectedOperationCode: 1,
			expectedResponseSize:  4,
			err:                   ErrOperationCodeUnexpected,
			errMessage:            "operation code is unexpected: expected 0x1 and got 0x2",
		},
		"result_code_failure": {
			response:              []byte{0, 1, 0, 1},
			expectedOperationCode: 1,
			expectedResponseSize:  4,
			err:                   ErrVersionNotSupported,
			errMessage:            "result code: version is not supported",
		},
		"success": {
			response:              []byte{0, 1, 0, 0},
			expectedOperationCode: 1,
			expectedResponseSize:  4,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := checkResponse(testCase.response,
				testCase.expectedOperationCode,
				testCase.expectedResponseSize)

			assert.ErrorIs(t, err, testCase.err)
			if testCase.err != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
		})
	}
}

func Test_checkResultCode(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		resultCode uint16
		err        error
		errMessage string
	}{
		"success": {},
		"version_unsupported": {
			resultCode: 1,
			err:        ErrVersionNotSupported,
			errMessage: "version is not supported",
		},
		"not_authorized": {
			resultCode: 2,
			err:        ErrNotAuthorized,
			errMessage: "not authorized",
		},
		"network_failure": {
			resultCode: 3,
			err:        ErrNetworkFailure,
			errMessage: "network failure",
		},
		"out_of_resources": {
			resultCode: 4,
			err:        ErrOutOfResources,
			errMessage: "out of resources",
		},
		"unsupported_operation_code": {
			resultCode: 5,
			err:        ErrOperationCodeNotSupported,
			errMessage: "operation code is not supported",
		},
		"unknown": {
			resultCode: 6,
			err:        ErrResultCodeUnknown,
			errMessage: "result code is unknown: 6",
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := checkResultCode(testCase.resultCode)

			assert.ErrorIs(t, err, testCase.err)
			if testCase.err != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
		})
	}
}
