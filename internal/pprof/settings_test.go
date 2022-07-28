package pprof

import (
	"net/http"
	"testing"
	"time"

	"github.com/qdm12/gluetun/internal/httpserver"
	"github.com/qdm12/govalid/address"
	"github.com/stretchr/testify/assert"
)

func Test_Settings_SetDefaults(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initial  Settings
		expected Settings
	}{
		"empty settings": {
			expected: Settings{
				Enabled: boolPtr(false),
				HTTPServer: httpserver.Settings{
					Address:           "localhost:6060",
					ReadHeaderTimeout: 3 * time.Second,
					ReadTimeout:       5 * time.Minute,
					ShutdownTimeout:   3 * time.Second,
				},
			},
		},
		"non empty settings": {
			initial: Settings{
				Enabled:          boolPtr(true),
				BlockProfileRate: 1,
				MutexProfileRate: 1,
				HTTPServer: httpserver.Settings{
					Address:           ":6061",
					ReadHeaderTimeout: time.Second,
					ReadTimeout:       time.Second,
					ShutdownTimeout:   time.Second,
				},
			},
			expected: Settings{
				Enabled:          boolPtr(true),
				BlockProfileRate: 1,
				MutexProfileRate: 1,
				HTTPServer: httpserver.Settings{
					Address:           ":6061",
					ReadHeaderTimeout: time.Second,
					ReadTimeout:       time.Second,
					ShutdownTimeout:   time.Second,
				},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			testCase.initial.SetDefaults()

			assert.Equal(t, testCase.expected, testCase.initial)
		})
	}
}

func Test_Settings_Copy(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initial  Settings
		expected Settings
	}{
		"empty settings": {},
		"non empty settings": {
			initial: Settings{
				Enabled:          boolPtr(true),
				BlockProfileRate: 1,
				MutexProfileRate: 1,
				HTTPServer: httpserver.Settings{
					Address:         ":6061",
					ShutdownTimeout: time.Second,
				},
			},
			expected: Settings{
				Enabled:          boolPtr(true),
				BlockProfileRate: 1,
				MutexProfileRate: 1,
				HTTPServer: httpserver.Settings{
					Address:         ":6061",
					ShutdownTimeout: time.Second,
				},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			copied := testCase.initial.Copy()

			assert.Equal(t, testCase.expected, copied)
		})
	}
}

func Test_Settings_MergeWith(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		settings Settings
		other    Settings
		expected Settings
	}{
		"merge empty with empty": {},
		"merge empty with filled": {
			other: Settings{
				Enabled:          boolPtr(true),
				BlockProfileRate: 1,
				MutexProfileRate: 1,
				HTTPServer: httpserver.Settings{
					Address: ":8001",
				},
			},
			expected: Settings{
				Enabled:          boolPtr(true),
				BlockProfileRate: 1,
				MutexProfileRate: 1,
				HTTPServer: httpserver.Settings{
					Address: ":8001",
				},
			},
		},
		"merge filled with empty": {
			settings: Settings{
				Enabled:          boolPtr(true),
				BlockProfileRate: 1,
				MutexProfileRate: 1,
				HTTPServer: httpserver.Settings{
					Address: ":8001",
				},
			},
			expected: Settings{
				Enabled:          boolPtr(true),
				BlockProfileRate: 1,
				MutexProfileRate: 1,
				HTTPServer: httpserver.Settings{
					Address: ":8001",
				},
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

	testCases := map[string]struct {
		settings Settings
		other    Settings
		expected Settings
	}{
		"override empty with empty": {},
		"override empty with filled": {
			other: Settings{
				Enabled:          boolPtr(true),
				BlockProfileRate: 1,
				MutexProfileRate: 1,
				HTTPServer: httpserver.Settings{
					Address: ":8001",
				},
			},
			expected: Settings{
				Enabled:          boolPtr(true),
				BlockProfileRate: 1,
				MutexProfileRate: 1,
				HTTPServer: httpserver.Settings{
					Address: ":8001",
				},
			},
		},
		"override filled with empty": {
			settings: Settings{
				Enabled:          boolPtr(true),
				BlockProfileRate: 1,
				MutexProfileRate: 1,
				HTTPServer: httpserver.Settings{
					Address: ":8001",
				},
			},
			expected: Settings{
				Enabled:          boolPtr(true),
				BlockProfileRate: 1,
				MutexProfileRate: 1,
				HTTPServer: httpserver.Settings{
					Address: ":8001",
				},
			},
		},
		"override filled with filled": {
			settings: Settings{
				Enabled:          boolPtr(false),
				BlockProfileRate: 1,
				MutexProfileRate: 1,
				HTTPServer: httpserver.Settings{
					Address: ":8001",
				},
			},
			other: Settings{
				Enabled:          boolPtr(true),
				BlockProfileRate: 2,
				MutexProfileRate: 3,
				HTTPServer: httpserver.Settings{
					Address: ":8002",
				},
			},
			expected: Settings{
				Enabled:          boolPtr(true),
				BlockProfileRate: 2,
				MutexProfileRate: 3,
				HTTPServer: httpserver.Settings{
					Address: ":8002",
				},
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

	testCases := map[string]struct {
		settings   Settings
		errWrapped error
		errMessage string
	}{
		"negative block profile rate": {
			settings: Settings{
				BlockProfileRate: -1,
			},
			errWrapped: ErrBlockProfileRateNegative,
			errMessage: ErrBlockProfileRateNegative.Error(),
		},
		"negative mutex profile rate": {
			settings: Settings{
				MutexProfileRate: -1,
			},
			errWrapped: ErrMutexProfileRateNegative,
			errMessage: ErrMutexProfileRateNegative.Error(),
		},
		"http server validation error": {
			settings: Settings{
				HTTPServer: httpserver.Settings{},
			},
			errWrapped: address.ErrValueNotValid,
			errMessage: "value is not valid: missing port in address",
		},
		"valid settings": {
			settings: Settings{
				HTTPServer: httpserver.Settings{
					Address:           ":8000",
					Handler:           http.NewServeMux(),
					Logger:            &MockLogger{},
					ReadHeaderTimeout: time.Second,
					ReadTimeout:       time.Second,
					ShutdownTimeout:   time.Second,
				},
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
		"disabled pprof": {
			settings: Settings{
				Enabled: boolPtr(false),
			},
		},
		"all values": {
			settings: Settings{
				Enabled:          boolPtr(true),
				BlockProfileRate: 2,
				MutexProfileRate: 1,
				HTTPServer: httpserver.Settings{
					Address:         ":8000",
					ShutdownTimeout: time.Second,
				},
			},
			s: `Pprof settings:
├── Block profile rate: 2
├── Mutex profile rate: 1
└── HTTP server settings:
    ├── Listening address: :8000
    ├── Read header timeout: 0s
    ├── Read timeout: 0s
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
