//go:build linux

package netlink

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

func (n *NetLink) IsWireguardSupported() (ok bool, err error) {
	families, err := netlink.GenlFamilyList()
	if err != nil {
		return false, fmt.Errorf("listing gen 1 families: %w", err)
	}
	for _, family := range families {
		if family.Name == "wireguard" {
			return true, nil
		}
	}
	return false, nil
}
