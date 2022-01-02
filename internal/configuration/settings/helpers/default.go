package helpers

import (
	"net"
	"time"

	"github.com/qdm12/golibs/logging"
)

func DefaultInt(existing *int, defaultValue int) (
	result *int) {
	if existing != nil {
		return existing
	}
	result = new(int)
	*result = defaultValue
	return result
}

func DefaultUint8(existing *uint8, defaultValue uint8) (
	result *uint8) {
	if existing != nil {
		return existing
	}
	result = new(uint8)
	*result = defaultValue
	return result
}

func DefaultUint16(existing *uint16, defaultValue uint16) (
	result *uint16) {
	if existing != nil {
		return existing
	}
	result = new(uint16)
	*result = defaultValue
	return result
}

func DefaultBool(existing *bool, defaultValue bool) (
	result *bool) {
	if existing != nil {
		return existing
	}
	result = new(bool)
	*result = defaultValue
	return result
}

func DefaultString(existing string, defaultValue string) (
	result string) {
	if existing != "" {
		return existing
	}
	return defaultValue
}

func DefaultStringPtr(existing *string, defaultValue string) (result *string) {
	if existing != nil {
		return existing
	}
	result = new(string)
	*result = defaultValue
	return result
}

func DefaultDuration(existing *time.Duration,
	defaultValue time.Duration) (result *time.Duration) {
	if existing != nil {
		return existing
	}
	result = new(time.Duration)
	*result = defaultValue
	return result
}

func DefaultLogLevel(existing *logging.Level,
	defaultValue logging.Level) (result *logging.Level) {
	if existing != nil {
		return existing
	}
	result = new(logging.Level)
	*result = defaultValue
	return result
}

func DefaultIP(existing net.IP, defaultValue net.IP) (
	result net.IP) {
	if existing != nil {
		return existing
	}
	return defaultValue
}
