package routing

import (
	"fmt"
	"net"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:generate mockgen -destination=mockCommander_test.go -package=routing github.com/qdm12/golibs/command Commander

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
			mockCommander := NewMockCommander(mockCtrl)

			mockCommander.EXPECT().Run("ip", "route", "del", tc.subnet.String()).
				Return(tc.runOutput, tc.runErr).Times(1)
			r := &routing{commander: mockCommander}
			err := r.removeRoute(tc.subnet)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
