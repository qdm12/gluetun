package netlink

import (
	"errors"
	"fmt"
	"os"

	"github.com/mdlayher/genetlink"
	"github.com/qdm12/gluetun/internal/mod"
)

func (n *NetLink) IsWireguardSupported() (ok bool, err error) {
	// Check for Wireguard family without loading the wireguard module.
	// Some kernels have the wireguard module built-in, and don't have a
	// modules directory, such as WSL2 kernels.
	ok, err = hasWireguardFamily()
	if err != nil {
		return false, fmt.Errorf("checking wireguard family: %w", err)
	} else if ok {
		return true, nil
	}

	// Try loading the wireguard module, since some systems do not load
	// it after a boot. If this fails, wireguard is assumed to not be supported.
	n.debugLogger.Debugf("wireguard family not found, trying to load wireguard kernel module")
	err = mod.Probe("wireguard")
	if err != nil {
		n.debugLogger.Debugf("failed loading wireguard kernel module: %s", err)
		return false, nil
	}
	n.debugLogger.Debugf("wireguard kernel module loaded successfully")

	// Re-check if the Wireguard family is now available, after loading
	// the wireguard kernel module.
	ok, err = hasWireguardFamily()
	if err != nil {
		return false, fmt.Errorf("checking wireguard family: %w", err)
	}
	return ok, nil
}

func hasWireguardFamily() (ok bool, err error) {
	conn, err := genetlink.Dial(nil)
	if err != nil {
		return false, fmt.Errorf("dialing netlink: %w", err)
	}
	defer conn.Close()

	_, err = conn.GetFamily("wireguard")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("getting wireguard family: %w", err)
	}

	return true, nil
}
