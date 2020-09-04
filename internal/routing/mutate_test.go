package routing

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/command/mock_command"
	"github.com/qdm12/golibs/files/mock_files"
	"github.com/qdm12/golibs/logging/mock_logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_DeleteRouteVia(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tests := map[string]struct {
		subnet    net.IPNet
		runOutput string
		runErr    error
		err       error
	}{
		"no output no error": {
			subnet: net.IPNet{
				IP:   net.IP{192, 168, 2, 0},
				Mask: net.IPMask{255, 255, 255, 0},
			},
		},
		"error only": {
			subnet: net.IPNet{
				IP:   net.IP{192, 168, 2, 0},
				Mask: net.IPMask{255, 255, 255, 0},
			},
			runErr: fmt.Errorf("error"),
			err:    fmt.Errorf("cannot delete route for 192.168.2.0/24: : error"),
		},
		"error and output": {
			subnet: net.IPNet{
				IP:   net.IP{192, 168, 2, 0},
				Mask: net.IPMask{255, 255, 255, 0},
			},
			runErr:    fmt.Errorf("error"),
			runOutput: "output",
			err:       fmt.Errorf("cannot delete route for 192.168.2.0/24: output: error"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			subnetStr := tc.subnet.String()

			logger := mock_logging.NewMockLogger(mockCtrl)
			logger.EXPECT().Info("deleting route for %s")
			commander := mock_command.NewMockCommander(mockCtrl)
			commander.EXPECT().Run(ctx, "ip", "route", "del", subnetStr).
				Return(tc.runOutput, tc.runErr).Times(1)
			fileManager := mock_files.NewMockFileManager(mockCtrl)
			routesData := []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0   0002A8C0  0100000A  0003   0 0 0  00FFFFFF   0 0  0
`)
			fileManager.EXPECT().ReadFile(string(constants.NetRoute)).Return(routesData, nil)
			r := &routing{
				logger:      logger,
				commander:   commander,
				fileManager: fileManager,
			}

			err := r.DeleteRouteVia(ctx, tc.subnet)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
