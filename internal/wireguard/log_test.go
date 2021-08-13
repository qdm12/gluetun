package wireguard

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_makeDeviceLogger(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	logger := NewMockLogger(ctrl)

	deviceLogger := makeDeviceLogger(logger)

	logger.EXPECT().Debug("test 1")
	deviceLogger.Verbosef("test %d", 1)

	logger.EXPECT().Error("test 2")
	deviceLogger.Errorf("test %d", 2)
}
