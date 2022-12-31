package netlink

import (
	"fmt"

	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
)

//nolint:revive
const (
	FAMILY_ALL = netlink.FAMILY_ALL
	FAMILY_V4  = netlink.FAMILY_V4
	FAMILY_V6  = netlink.FAMILY_V6
)

func (n *NetLink) IsWireguardSupported() (ok bool, err error) {
	const wireguardLinkName = "test"
	device, err := tun.CreateTUN(wireguardLinkName, device.DefaultMTU)
	if err != nil {
		return false, fmt.Errorf("creating tun device: %w", err)
	}

	err = device.Close()
	if err != nil {
		return false, fmt.Errorf("closing tun device: %w", err)
	}

	families, err := netlink.GenlFamilyList()
	if err != nil {
		return false, fmt.Errorf("listing genl families: %w", err)
	}
	for _, family := range families {
		if family.Name == "wireguard" {
			return true, nil
		}
	}

	return false, nil
}
