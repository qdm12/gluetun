package httpserver

import (
	"net/http"
	"testing"
	"time"

	"github.com/qdm12/govalid/address"
	"github.com/stretchr/testify/assert"
)

func Test_Settings_SetDefaults(t *testing.T) {
	t.Parallel()

	const defaultTimeout = 3 * time.Second

	testCases := map[string]struct {
		settings Settings
		expected Settings
	}{
		"empty settings": {
			settings: Settings{},
			expected: Settings{
				Address:           ":8000",
				ReadHeaderTimeout: defaultTimeout,
				ReadTimeout:       defaultTimeout,
				ShutdownTimeout:   defaultTimeout,
			},
		},
		"filled settings": {
			settings: Settings{
				Address:           ":8001",
				ReadHeaderTimeout: time.Second,
				ReadTimeout:       time.Second,
				ShutdownTimeout:   time.Second,
			},
			expected: Settings{
				Address:           ":8001",
				ReadHeaderTimeout: time.Second,
				ReadTimeout:       time.Second,
				ShutdownTimeout:   time.Second,
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			testCase.settings.SetDefaults()

			assert.Equal(t, testCase.expected, testCase.settings)
		})
	}
}

func Test_Settings_Copy(t *testing.T) {
	t.Parallel()

	someHandler := http.NewServeMux()
	someLogger := &testLogger{}

	testCases := map[string]struct {
		settings Settings
		expected Settings
	}{
		"empty settings": {},
		"filled settings": {
			settings: Settings{
				Address:           ":8001",
				Handler:           someHandler,
				Logger:            someLogger,
				ReadHeaderTimeout: time.Second,
				ReadTimeout:       time.Second,
				ShutdownTimeout:   time.Second,
			},
			expected: Settings{
				Address:           ":8001",
				Handler:           someHandler,
				Logger:            someLogger,
				ReadHeaderTimeout: time.Second,
				ReadTimeout:       time.Second,
				ShutdownTimeout:   time.Second,
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			copied := testCase.settings.Copy()

			assert.Equal(t, testCase.expected, copied)
		})
	}
}

func Test_Settings_MergeWith(t *testing.T) {
	t.Parallel()

	someHandler := http.NewServeMux()
	someLogger := &testLogger{}

	testCases := map[string]struct {
		settings Settings
		other    Settings
		expected Settings
	}{
		"merge empty with empty": {},
		"merge empty with filled": {
			other: Settings{
				Address:           ":8001",
				Handler:           someHandler,
				Logger:            someLogger,
				ReadHeaderTimeout: time.Second,
				ReadTimeout:       time.Second,
				ShutdownTimeout:   time.Second,
			},
			expected: Settings{
				Address:           ":8001",
				Handler:           someHandler,
				Logger:            someLogger,
				ReadHeaderTimeout: time.Second,
				ReadTimeout:       time.Second,
				ShutdownTimeout:   time.Second,
			},
		},
		"merge filled with empty": {
			settings: Settings{
				Address:           ":8001",
				Handler:           someHandler,
				Logger:            someLogger,
				ReadHeaderTimeout: time.Second,
				ReadTimeout:       time.Second,
				ShutdownTimeout:   time.Second,
			},
			expected: Settings{
				Address:           ":8001",
				Handler:           someHandler,
				Logger:            someLogger,
				ReadHeaderTimeout: time.Second,
				ReadTimeout:       time.Second,
				ShutdownTimeout:   time.Second,
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			testCase.settings.MergeWith(testCase.other)

			assert.Equal(t, testCase.expected, testCase.settings)
		})
	}
}

func Test_Settings_OverrideWith(t *testing.T) {
	t.Parallel()

	someHandler := http.NewServeMux()
	someLogger := &testLogger{}

	testCases := map[string]struct {
		settings Settings
		other    Settings
		expected Settings
	}{
		"override empty with empty": {},
		"override empty with filled": {
			other: Settings{
				Address:           ":8001",
				Handler:           someHandler,
				Logger:            someLogger,
				ReadHeaderTimeout: time.Second,
				ReadTimeout:       time.Second,
				ShutdownTimeout:   time.Second,
			},
			expected: Settings{
				Address:           ":8001",
				Handler:           someHandler,
				Logger:            someLogger,
				ReadHeaderTimeout: time.Second,
				ReadTimeout:       time.Second,
				ShutdownTimeout:   time.Second,
			},
		},
		"override filled with empty": {
			settings: Settings{
				Address:           ":8001",
				Handler:           someHandler,
				Logger:            someLogger,
				ReadHeaderTimeout: time.Second,
				ReadTimeout:       time.Second,
				ShutdownTimeout:   time.Second,
			},
			expected: Settings{
				Address:           ":8001",
				Handler:           someHandler,
				Logger:            someLogger,
				ReadHeaderTimeout: time.Second,
				ReadTimeout:       time.Second,
				ShutdownTimeout:   time.Second,
			},
		},
		"override filled with filled": {
			settings: Settings{
				Address:           ":8001",
				Handler:           someHandler,
				Logger:            someLogger,
				ReadHeaderTimeout: time.Second,
				ReadTimeout:       time.Second,
				ShutdownTimeout:   time.Second,
			},
			other: Settings{
				Address:           ":8002",
				ReadHeaderTimeout: time.Hour,
				ReadTimeout:       time.Hour,
				ShutdownTimeout:   time.Hour,
			},
			expected: Settings{
				Address:           ":8002",
				Handler:           someHandler,
				Logger:            someLogger,
				ReadHeaderTimeout: time.Hour,
				ReadTimeout:       time.Hour,
				ShutdownTimeout:   time.Hour,
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			testCase.settings.OverrideWith(testCase.other)

			assert.Equal(t, testCase.expected, testCase.settings)
		})
	}
}

func Test_Settings_Validate(t *testing.T) {
	t.Parallel()

	someHandler := http.NewServeMux()
	someLogger := &testLogger{}

	testCases := map[string]struct {
		settings   Settings
		errWrapped error
		errMessage string
	}{
		"bad address": {
			settings: Settings{
				Address: "noport",
			},
			errWrapped: address.ErrValueNotValid,
			errMessage: "value is not valid: address noport: missing port in address",
		},
		"nil handler": {
			settings: Settings{
				Address: ":8000",
			},
			errWrapped: ErrHandlerIsNotSet,
			errMessage: ErrHandlerIsNotSet.Error(),
		},
		"nil logger": {
			settings: Settings{
				Address: ":8000",
				Handler: someHandler,
			},
			errWrapped: ErrLoggerIsNotSet,
			errMessage: ErrLoggerIsNotSet.Error(),
		},
		"read header timeout too small": {
			settings: Settings{
				Address:           ":8000",
				Handler:           someHandler,
				Logger:            someLogger,
				ReadHeaderTimeout: time.Nanosecond,
			},
			errWrapped: ErrReadHeaderTimeoutTooSmall,
			errMessage: "read header timeout is too small: 1ns must be at least 1ms",
		},
		"read timeout too small": {
			settings: Settings{
				Address:           ":8000",
				Handler:           someHandler,
				Logger:            someLogger,
				ReadHeaderTimeout: time.Millisecond,
				ReadTimeout:       time.Nanosecond,
			},
			errWrapped: ErrReadTimeoutTooSmall,
			errMessage: "read timeout is too small: 1ns must be at least 1ms",
		},
		"shutdown timeout too small": {
			settings: Settings{
				Address:           ":8000",
				Handler:           someHandler,
				Logger:            someLogger,
				ReadHeaderTimeout: time.Millisecond,
				ReadTimeout:       time.Millisecond,
				ShutdownTimeout:   time.Millisecond,
			},
			errWrapped: ErrShutdownTimeoutTooSmall,
			errMessage: "shutdown timeout is too small: 1ms must be at least 5ms",
		},
		"valid settings": {
			settings: Settings{
				Address:           ":8000",
				Handler:           someHandler,
				Logger:            someLogger,
				ReadHeaderTimeout: time.Millisecond,
				ReadTimeout:       time.Millisecond,
				ShutdownTimeout:   time.Second,
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := testCase.settings.Validate()

			assert.ErrorIs(t, err, testCase.errWrapped)
			if err != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
		})
	}
}

func Test_Settings_String(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		settings Settings
		s        string
	}{
		"all values": {
			settings: Settings{
				Address:           ":8000",
				ReadHeaderTimeout: time.Millisecond,
				ReadTimeout:       time.Millisecond,
				ShutdownTimeout:   time.Second,
			},
			s: `HTTP server settings:
├── Listening address: :8000
├── Read header timeout: 1ms
├── Read timeout: 1ms
└── Shutdown timeout: 1s`,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := testCase.settings.String()

			assert.Equal(t, testCase.s, s)
		})
	}
}
