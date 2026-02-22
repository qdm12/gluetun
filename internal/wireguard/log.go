package wireguard

import (
	amneziaDevice "github.com/amnezia-vpn/amneziawg-go/device"
	wgDevice "golang.zx2c4.com/wireguard/device"
)

//go:generate mockgen -destination=log_mock_test.go -package wireguard . Logger

type Logger interface {
	Debug(s string)
	Debugf(format string, args ...interface{})
	Info(s string)
	Error(s string)
	Errorf(format string, args ...interface{})
}

func makeWgDeviceLogger(logger Logger) (deviceLogger *wgDevice.Logger) {
	return &wgDevice.Logger{
		Verbosef: logger.Debugf,
		Errorf:   logger.Errorf,
	}
}

func makeAmneziaDeviceLogger(logger Logger) (deviceLogger *amneziaDevice.Logger) {
	return &amneziaDevice.Logger{
		Verbosef: logger.Debugf,
		Errorf:   logger.Errorf,
	}
}
