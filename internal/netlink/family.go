package netlink

import (
	"fmt"
	"math/rand"

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
	ok, err = isWireguardFamilyPresent()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	// Some host systems need the tun device to be created
	// after a system boot in order to detect the wireguard family.
	// See https://github.com/qdm12/gluetun/issues/984
	wgInterfaceName := randomInterfaceName()
	device, err := tun.CreateTUN(wgInterfaceName, device.DefaultMTU)
	if err != nil {
		return false, fmt.Errorf("creating tun device: %w", err)
	}

	err = device.Close()
	if err != nil {
		return false, fmt.Errorf("closing tun device: %w", err)
	}

	return isWireguardFamilyPresent()
}

func isWireguardFamilyPresent() (found bool, err error) {
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

func randomInterfaceName() (interfaceName string) {
	const size = 15
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, size)
	for i := range b {
		letterIndex := rand.Intn(len(letterRunes)) //nolint:gosec
		b[i] = letterRunes[letterIndex]
	}
	return string(b)
}
