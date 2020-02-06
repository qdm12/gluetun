package dns

import (
	"fmt"
	"testing"

	commandMocks "github.com/qdm12/golibs/command/mocks"
	loggingMocks "github.com/qdm12/golibs/logging/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func Test_Start(t *testing.T) {
	t.Parallel()
	logger := &loggingMocks.Logger{}
	logger.On("Info", "%s: starting unbound", logPrefix).Once()
	commander := &commandMocks.Commander{}
	commander.On("Start", "unbound", "-d", "-c", string(constants.UnboundConf), "-vv").
		Return(nil, nil, nil).Once()
	c := &configurator{commander: commander, logger: logger}
	stdout, err := c.Start(2)
	assert.Nil(t, stdout)
	assert.NoError(t, err)
	logger.AssertExpectations(t)
	commander.AssertExpectations(t)
}

func Test_Version(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		runOutput string
		runErr    error
		version   string
		err       error
	}{
		"no data": {
			err: fmt.Errorf(`unbound version was not found in ""`),
		},
		"2 lines with version": {
			runOutput: "Version  \nVersion 1.0-a hello\n",
			version:   "1.0-a",
		},
		"run error": {
			runErr: fmt.Errorf("error"),
			err:    fmt.Errorf("unbound version: error"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			commander := &commandMocks.Commander{}
			commander.On("Run", "unbound", "-V").
				Return(tc.runOutput, tc.runErr).Once()
			c := &configurator{commander: commander}
			version, err := c.Version()
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.version, version)
			commander.AssertExpectations(t)
		})
	}
}
