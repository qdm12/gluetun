package env

import (
	"fmt"
	"testing"

	"github.com/qdm12/golibs/logging/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_FatalOnError(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		err error
	}{
		"nil": {},
		"err": {fmt.Errorf("error")},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var logged string
			var exitCode int
			logger := &mocks.Logger{}
			if tc.err != nil {
				logger.On("Error", tc.err).
					Run(func(args mock.Arguments) {
						err := args.Get(0).(error)
						logged = err.Error()
					}).Once()
			}
			osExit := func(n int) { exitCode = n }
			e := &env{logger, osExit}
			e.FatalOnError(tc.err)
			if tc.err != nil {
				assert.Equal(t, logged, tc.err.Error())
				assert.Equal(t, exitCode, 1)
			} else {
				assert.Empty(t, logged)
				assert.Zero(t, exitCode)
			}
		})
	}
}

func Test_PrintVersion(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		program        string
		commandVersion string
		commandErr     error
	}{
		"no data": {},
		"data":    {"binu", "2.3-5", nil},
		"error":   {"binu", "", fmt.Errorf("error")},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var logged string
			logger := &mocks.Logger{}
			if tc.commandErr != nil {
				logger.On("Error", tc.commandErr).
					Run(func(args mock.Arguments) {
						err := args.Get(0).(error)
						logged = err.Error()
					}).Once()
			} else {
				logger.On("Info", "%s version: %s", tc.program, tc.commandVersion).
					Run(func(args mock.Arguments) {
						format := args.Get(0).(string)
						program := args.Get(1).(string)
						version := args.Get(2).(string)
						logged = fmt.Sprintf(format, program, version)
					}).Once()
			}
			e := &env{logger: logger}
			commandFn := func() (string, error) { return tc.commandVersion, tc.commandErr }
			e.PrintVersion(tc.program, commandFn)
			if tc.commandErr != nil {
				assert.Equal(t, logged, tc.commandErr.Error())
			} else {
				assert.Equal(t, logged, fmt.Sprintf("%s version: %s", tc.program, tc.commandVersion))
			}
		})
	}
}
