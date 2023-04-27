package netlink

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

//nolint:revive
const (
	FAMILY_ALL = netlink.FAMILY_ALL
	FAMILY_V4  = netlink.FAMILY_V4
	FAMILY_V6  = netlink.FAMILY_V6
)

func FamilyToString(family int) string {
	switch family {
	case FAMILY_ALL:
		return "all"
	case FAMILY_V4:
		return "v4"
	case FAMILY_V6:
		return "v6"
	default:
		return fmt.Sprint(family)
	}
}

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
