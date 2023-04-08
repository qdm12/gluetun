package natpmp

import (
	"encoding/binary"
	"errors"
	"fmt"
)

var (
	ErrRequestSizeTooSmall = errors.New("message size is too small")
)

func checkRequest(request []byte) (err error) {
	const minMessageSize = 2 // version number + operation code
	if len(request) < minMessageSize {
		return fmt.Errorf("%w: need at least %d bytes and got %d byte(s)",
			ErrRequestSizeTooSmall, minMessageSize, len(request))
	}

	return nil
}

var (
	ErrResponseSizeTooSmall    = errors.New("response size is too small")
	ErrResponseSizeUnexpected  = errors.New("response size is unexpected")
	ErrProtocolVersionUnknown  = errors.New("protocol version is unknown")
	ErrOperationCodeUnexpected = errors.New("operation code is unexpected")
)

func checkResponse(response []byte, expectedOperationCode byte,
	expectedResponseSize uint) (err error) {
	const minResponseSize = 4
	if len(response) < minResponseSize {
		return fmt.Errorf("%w: need at least %d bytes and got %d byte(s)",
			ErrResponseSizeTooSmall, minResponseSize, len(response))
	}

	if len(response) != int(expectedResponseSize) {
		return fmt.Errorf("%w: expected %d bytes and got %d byte(s)",
			ErrResponseSizeUnexpected, expectedResponseSize, len(response))
	}

	protocolVersion := response[0]
	if protocolVersion != 0 {
		return fmt.Errorf("%w: %d", ErrProtocolVersionUnknown, protocolVersion)
	}

	operationCode := response[1]
	if operationCode != expectedOperationCode {
		return fmt.Errorf("%w: expected 0x%x and got 0x%x",
			ErrOperationCodeUnexpected, expectedOperationCode, operationCode)
	}

	resultCode := binary.BigEndian.Uint16(response[2:4])
	err = checkResultCode(resultCode)
	if err != nil {
		return fmt.Errorf("result code: %w", err)
	}

	return nil
}

var (
	ErrVersionNotSupported       = errors.New("version is not supported")
	ErrNotAuthorized             = errors.New("not authorized")
	ErrNetworkFailure            = errors.New("network failure")
	ErrOutOfResources            = errors.New("out of resources")
	ErrOperationCodeNotSupported = errors.New("operation code is not supported")
	ErrResultCodeUnknown         = errors.New("result code is unknown")
)

// checkResultCode checks the result code and returns an error
// if the result code is not a success (0).
// See https://www.ietf.org/rfc/rfc6886.html#section-3.5
//
//nolint:gomnd
func checkResultCode(resultCode uint16) (err error) {
	switch resultCode {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("%w", ErrVersionNotSupported)
	case 2:
		return fmt.Errorf("%w", ErrNotAuthorized)
	case 3:
		return fmt.Errorf("%w", ErrNetworkFailure)
	case 4:
		return fmt.Errorf("%w", ErrOutOfResources)
	case 5:
		return fmt.Errorf("%w", ErrOperationCodeNotSupported)
	default:
		return fmt.Errorf("%w: %d", ErrResultCodeUnknown, resultCode)
	}
}
