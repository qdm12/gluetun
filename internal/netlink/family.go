package netlink

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

const (
	FamilyAll = 0
	FamilyV4  = 2
	FamilyV6  = 10
)

func FamilyToString(family int) string {
	switch family {
	case FamilyAll:
		return "all" //nolint:goconst
	case FamilyV4:
		return "v4"
	case FamilyV6:
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
