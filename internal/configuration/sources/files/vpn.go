package files

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *Source) readVPN() (vpn settings.VPN, err error) {
	vpn.Provider, err = s.readProvider()
	if err != nil {
		return vpn, fmt.Errorf("provider: %w", err)
	}

	vpn.OpenVPN, err = s.readOpenVPN()
	if err != nil {
		return vpn, fmt.Errorf("OpenVPN: %w", err)
	}

	vpn.Wireguard, err = s.readWireguard()
	if err != nil {
		return vpn, fmt.Errorf("wireguard: %w", err)
	}

	return vpn, nil
}
