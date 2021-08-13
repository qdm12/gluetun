package wireguard

import (
	"errors"
	"fmt"
	"net"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var (
	errMakeConfig      = errors.New("cannot make device configuration")
	errConfigureDevice = errors.New("cannot configure device")
)

func configureDevice(client *wgctrl.Client, settings Settings) (err error) {
	deviceConfig, err := makeDeviceConfig(settings)
	if err != nil {
		return fmt.Errorf("%w: %s", errMakeConfig, err)
	}

	err = client.ConfigureDevice(settings.InterfaceName, deviceConfig)
	if err != nil {
		return fmt.Errorf("%w: %s", errConfigureDevice, err)
	}

	return nil
}

func makeDeviceConfig(settings Settings) (config wgtypes.Config, err error) {
	privateKey, err := wgtypes.ParseKey(settings.PrivateKey)
	if err != nil {
		return config, ErrPrivateKeyInvalid
	}

	publicKey, err := wgtypes.ParseKey(settings.PublicKey)
	if err != nil {
		return config, fmt.Errorf("%w: %s", ErrPublicKeyInvalid, settings.PublicKey)
	}

	firewallMark := settings.FirewallMark

	config = wgtypes.Config{
		PrivateKey:   &privateKey,
		ReplacePeers: true,
		FirewallMark: &firewallMark,
		Peers: []wgtypes.PeerConfig{
			{
				PublicKey: publicKey,
				AllowedIPs: []net.IPNet{
					*allIPv4(),
					*allIPv6(),
				},
				ReplaceAllowedIPs: true,
				Endpoint:          settings.Endpoint,
			},
		},
	}

	return config, nil
}

func allIPv4() (ipNet *net.IPNet) {
	return &net.IPNet{
		IP:   net.IPv4(0, 0, 0, 0),
		Mask: []byte{0, 0, 0, 0},
	}
}

func allIPv6() (ipNet *net.IPNet) {
	return &net.IPNet{
		IP:   net.IPv6zero,
		Mask: []byte(net.IPv6zero),
	}
}
