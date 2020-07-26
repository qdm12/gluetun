package dns

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/golibs/command/mock_command"
	"github.com/qdm12/golibs/logging/mock_logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/qdm12/gluetun/internal/constants"
)

func Test_Start(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	logger := mock_logging.NewMockLogger(mockCtrl)
	logger.EXPECT().Info("starting unbound").Times(1)
	commander := mock_command.NewMockCommander(mockCtrl)
	commander.EXPECT().Start(context.Background(), "unbound", "-d", "-c", string(constants.UnboundConf), "-vv").
		Return(nil, nil, nil, nil).Times(1)
	c := &configurator{commander: commander, logger: logger}
	stdout, waitFn, err := c.Start(context.Background(), 2)
	assert.Nil(t, stdout)
	assert.Nil(t, waitFn)
	assert.NoError(t, err)
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
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			commander := mock_command.NewMockCommander(mockCtrl)
			commander.EXPECT().Run(context.Background(), "unbound", "-V").
				Return(tc.runOutput, tc.runErr).Times(1)
			c := &configurator{commander: commander}
			version, err := c.Version(context.Background())
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.version, version)
		})
	}
}
