package routing

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/golibs/command/mock_command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_removeRoute(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		subnet    net.IPNet
		runOutput string
		runErr    error
		err       error
	}{
		"no output no error": {
			subnet: net.IPNet{
				IP:   net.IP{192, 168, 1, 0},
				Mask: net.IPMask{255, 255, 255, 0},
			},
		},
		"error only": {
			subnet: net.IPNet{
				IP:   net.IP{192, 168, 1, 0},
				Mask: net.IPMask{255, 255, 255, 0},
			},
			runErr: fmt.Errorf("error"),
			err:    fmt.Errorf("cannot delete route for 192.168.1.0/24: : error"),
		},
		"error and output": {
			subnet: net.IPNet{
				IP:   net.IP{192, 168, 1, 0},
				Mask: net.IPMask{255, 255, 255, 0},
			},
			runErr:    fmt.Errorf("error"),
			runOutput: "output",
			err:       fmt.Errorf("cannot delete route for 192.168.1.0/24: output: error"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			commander := mock_command.NewMockCommander(mockCtrl)

			commander.EXPECT().Run(context.Background(), "ip", "route", "del", tc.subnet.String()).
				Return(tc.runOutput, tc.runErr).Times(1)
			r := &routing{commander: commander}
			err := r.removeRoute(context.Background(), tc.subnet)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
