package secrets

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *Source) readVPN() (vpn settings.VPN, err error) {
	vpn.OpenVPN, err = s.readOpenVPN()
	if err != nil {
		return vpn, fmt.Errorf("reading OpenVPN settings: %w", err)
	}

	vpn.Wireguard, err = s.readWireguard()
	if err != nil {
		return vpn, fmt.Errorf("reading Wireguard settings: %w", err)
	}

	return vpn, nil
}
