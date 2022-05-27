package wireguard

import (
	"golang.zx2c4.com/wireguard/device"
)

//go:generate mockgen -destination=log_mock_test.go -package wireguard . Logger

type Logger interface {
	Debug(s string)
	Debugf(format string, args ...interface{})
	Info(s string)
	Error(s string)
	Errorf(format string, args ...interface{})
}

func makeDeviceLogger(logger Logger) (deviceLogger *device.Logger) {
	return &device.Logger{
		Verbosef: logger.Debugf,
		Errorf:   logger.Errorf,
	}
}
