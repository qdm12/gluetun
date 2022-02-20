package wireguard

import (
	"fmt"
	"net"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func configureDevice(client *wgctrl.Client, settings Settings) (err error) {
	deviceConfig, err := makeDeviceConfig(settings)
	if err != nil {
		return fmt.Errorf("cannot make device configuration: %w", err)
	}

	err = client.ConfigureDevice(settings.InterfaceName, deviceConfig)
	if err != nil {
		return fmt.Errorf("cannot configure device: %w", err)
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

	var preSharedKey *wgtypes.Key
	if settings.PreSharedKey != "" {
		preSharedKeyValue, err := wgtypes.ParseKey(settings.PreSharedKey)
		if err != nil {
			return config, ErrPreSharedKeyInvalid
		}
		preSharedKey = &preSharedKeyValue
	}

	firewallMark := settings.FirewallMark

	config = wgtypes.Config{
		PrivateKey:   &privateKey,
		ReplacePeers: true,
		FirewallMark: &firewallMark,
		Peers: []wgtypes.PeerConfig{
			{
				PublicKey:    publicKey,
				PresharedKey: preSharedKey,
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
