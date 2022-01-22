package env

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/qdm12/govalid/binary"
	"github.com/qdm12/govalid/integer"
)

func envToCSV(envKey string) (values []string) {
	csv := os.Getenv(envKey)
	if csv == "" {
		return nil
	}
	return lowerAndSplit(csv)
}

func envToInt(envKey string) (n int, err error) {
	s := os.Getenv(envKey)
	if s == "" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

func envToStringPtr(envKey string) (stringPtr *string) {
	s := os.Getenv(envKey)
	if s == "" {
		return nil
	}
	return &s
}

func envToBoolPtr(envKey string) (boolPtr *bool, err error) {
	s := os.Getenv(envKey)
	if s == "" {
		return nil, nil //nolint:nilnil
	}
	value, err := binary.Validate(s)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func envToIntPtr(envKey string) (intPtr *int, err error) {
	s := os.Getenv(envKey)
	if s == "" {
		return nil, nil //nolint:nilnil
	}
	value, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func envToUint8Ptr(envKey string) (uint8Ptr *uint8, err error) {
	s := os.Getenv(envKey)
	if s == "" {
		return nil, nil //nolint:nilnil
	}

	const min, max = 0, 255
	value, err := integer.Validate(s, integer.OptionRange(min, max))
	if err != nil {
		return nil, err
	}

	uint8Ptr = new(uint8)
	*uint8Ptr = uint8(value)
	return uint8Ptr, nil
}

func envToUint16Ptr(envKey string) (uint16Ptr *uint16, err error) {
	s := os.Getenv(envKey)
	if s == "" {
		return nil, nil //nolint:nilnil
	}

	const min, max = 0, 65535
	value, err := integer.Validate(s, integer.OptionRange(min, max))
	if err != nil {
		return nil, err
	}

	uint16Ptr = new(uint16)
	*uint16Ptr = uint16(value)
	return uint16Ptr, nil
}

func envToDurationPtr(envKey string) (durationPtr *time.Duration, err error) {
	s := os.Getenv(envKey)
	if s == "" {
		return nil, nil //nolint:nilnil
	}

	durationPtr = new(time.Duration)
	*durationPtr, err = time.ParseDuration(s)
	if err != nil {
		return nil, err
	}

	return durationPtr, nil
}

func lowerAndSplit(csv string) (values []string) {
	csv = strings.ToLower(csv)
	return strings.Split(csv, ",")
}

var ErrDecodeBase64 = errors.New("cannot decode base64 string")

func decodeBase64(b64String string) (decoded string, err error) {
	b, err := base64.StdEncoding.DecodeString(b64String)
	if err != nil {
		return "", fmt.Errorf("%w: %s: %s",
			ErrDecodeBase64, b64String, err)
	}
	return string(b), nil
}

func unsetEnvKeys(envKeys []string, err error) (newErr error) {
	newErr = err
	for _, envKey := range envKeys {
		unsetErr := os.Unsetenv(envKey)
		if unsetErr != nil && newErr == nil {
			newErr = fmt.Errorf("cannot unset environment variable %s: %w", envKey, unsetErr)
		}
	}
	return newErr
}

func stringPtr(s string) *string { return &s }
func uint16Ptr(n uint16) *uint16 { return &n }
func boolPtr(b bool) *bool       { return &b }
