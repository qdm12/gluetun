//go:build linux

package netlink

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/mod"
	"github.com/vishvananda/netlink"
)

func (n *NetLink) IsWireguardSupported() (ok bool, err error) {
	// Check for Wireguard family without loading the wireguard module.
	// Some kernels have the wireguard module built-in, and don't have a
	// modules directory, such as WSL2 kernels.
	ok, err = hasWireguardFamily()
	if err != nil {
		return false, fmt.Errorf("checking for wireguard family: %w", err)
	}
	if ok {
		return true, nil
	}

	// Try loading the wireguard module, since some systems do not load
	// it after a boot. If this fails, wireguard is assumed to not be supported.
	err = mod.Probe("wireguard")
	if err != nil {
		n.debugLogger.Debugf("failed trying to load Wireguard kernel module: %s", err)
		return false, nil
	}

	// Re-check if the Wireguard family is now available, after loading
	// the wireguard kernel module.
	ok, err = hasWireguardFamily()
	if err != nil {
		return false, fmt.Errorf("checking for wireguard family: %w", err)
	}
	return ok, nil
}

func hasWireguardFamily() (ok bool, err error) {
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
