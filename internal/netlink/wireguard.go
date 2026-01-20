//go:build linux

package netlink

import (
	"github.com/qdm12/gluetun/internal/mod"
	"github.com/vishvananda/netlink"
)

func (n *NetLink) IsWireguardSupported() bool {
	// Check for Wireguard family without loading the wireguard module.
	// Some kernels have the wireguard module built-in, and don't have a
	// modules directory, such as WSL2 kernels.
	ok := hasWireguardFamily()
	if ok {
		return true
	}

	// Try loading the wireguard module, since some systems do not load
	// it after a boot. If this fails, wireguard is assumed to not be supported.
	n.debugLogger.Debugf("wireguard family not found, trying to load wireguard kernel module")
	err := mod.Probe("wireguard")
	if err != nil {
		n.debugLogger.Debugf("failed loading wireguard kernel module: %s", err)
		return false
	}
	n.debugLogger.Debugf("wireguard kernel module loaded successfully")

	// Re-check if the Wireguard family is now available, after loading
	// the wireguard kernel module.
	return hasWireguardFamily()
}

func hasWireguardFamily() bool {
	_, err := netlink.GenlFamilyGet("wireguard")
	return err == nil
}
