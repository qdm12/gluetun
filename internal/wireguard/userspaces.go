package wireguard

import (
	amneziaconn "github.com/amnezia-vpn/amneziawg-go/conn"
	amneziadevice "github.com/amnezia-vpn/amneziawg-go/device"
	amneziatun "github.com/amnezia-vpn/amneziawg-go/tun"
	wgconn "golang.zx2c4.com/wireguard/conn"
	wgdevice "golang.zx2c4.com/wireguard/device"
	wgtun "golang.zx2c4.com/wireguard/tun"
)

func defaultUserSpaceBackend() userSpaceBackend {
	return userSpaceBackend{
		createTun: func(name string, mtu int) (tunDevice, error) {
			return wgtun.CreateTUN(name, mtu)
		},
		createBind: func() bind {
			return wgconn.NewDefaultBind()
		},
		createDevice: func(td tunDevice, b bind, logger Logger) userspaceDevice {
			wgtun, _ := td.(wgtun.Device)
			wgBind, _ := b.(wgconn.Bind)
			wgLogger := wgdevice.Logger{
				Verbosef: logger.Debugf,
				Errorf:   logger.Errorf,
			}
			device := wgdevice.NewDevice(wgtun, wgBind, &wgLogger)
			return device
		},
		preStart: nil,
	}
}

func amneziaUserSpaceBackend() userSpaceBackend {
	return userSpaceBackend{
		createTun: func(name string, mtu int) (tunDevice, error) {
			return amneziatun.CreateTUN(name, mtu)
		},
		createBind: func() bind {
			return amneziaconn.NewDefaultBind()
		},
		createDevice: func(td tunDevice, b bind, logger Logger) userspaceDevice {
			wgamneziaTun, _ := td.(amneziatun.Device)
			wgamneziaBind, _ := b.(amneziaconn.Bind)
			wgamneziaLogger := amneziadevice.Logger{
				Verbosef: logger.Debugf,
				Errorf:   logger.Errorf,
			}
			device := amneziadevice.NewDevice(wgamneziaTun, wgamneziaBind, &wgamneziaLogger)
			return device
		},
		preStart: func(ud userspaceDevice, s Settings) error {
			uapiConfig := s.AmneziaWG.uapiConfig()
			err := ud.IpcSet(uapiConfig)
			return err
		},
	}
}
