package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/qdm12/gosettings/sources/env"
	"github.com/qdm12/govalid/binary"
	"github.com/qdm12/govalid/integer"
)

func envToCSV(envKey string) (values []string) {
	csv := env.Get(envKey)
	if csv == "" {
		return nil
	}
	return lowerAndSplit(csv)
}

func envToFloat64(envKey string) (f float64, err error) {
	s := env.Get(envKey)
	if s == "" {
		return 0, nil
	}
	const bits = 64
	return strconv.ParseFloat(s, bits)
}

func envToStringPtr(envKey string) (stringPtr *string) {
	s := env.Get(envKey)
	if s == "" {
		return nil
	}
	return &s
}

func envToBoolPtr(envKey string) (boolPtr *bool, err error) {
	s := env.Get(envKey)
	value, err := binary.Validate(s)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func envToIntPtr(envKey string) (intPtr *int, err error) {
	s := env.Get(envKey)
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
	s := env.Get(envKey)
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
	s := env.Get(envKey)
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
	s := env.Get(envKey)
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

func unsetEnvKeys(envKeys []string, err error) (newErr error) {
	newErr = err
	for _, envKey := range envKeys {
		unsetErr := os.Unsetenv(envKey)
		if unsetErr != nil && newErr == nil {
			newErr = fmt.Errorf("unsetting environment variable %s: %w", envKey, unsetErr)
		}
	}
	return newErr
}

func stringPtr(s string) *string { return &s }
func uint32Ptr(n uint32) *uint32 { return &n }
func boolPtr(b bool) *bool       { return &b }
