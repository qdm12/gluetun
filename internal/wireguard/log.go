package wireguard

import (
	"fmt"

	"golang.zx2c4.com/wireguard/device"
)

//go:generate mockgen -destination=log_mock_test.go -package wireguard . Logger

type Logger interface {
	Debug(s string)
	Info(s string)
	Error(s string)
}

func makeDeviceLogger(logger Logger) (deviceLogger *device.Logger) {
	return &device.Logger{
		Verbosef: func(format string, args ...interface{}) {
			logger.Debug(fmt.Sprintf(format, args...))
		},
		Errorf: func(format string, args ...interface{}) {
			logger.Error(fmt.Sprintf(format, args...))
		},
	}
}
