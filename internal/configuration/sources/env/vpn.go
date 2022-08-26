package env

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *Source) readVPN() (vpn settings.VPN, err error) {
	vpn.Type = strings.ToLower(getCleanedEnv("VPN_TYPE"))

	vpn.Provider, err = s.readProvider(vpn.Type)
	if err != nil {
		return vpn, fmt.Errorf("VPN provider: %w", err)
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
