package secrets

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func readVPN() (vpn settings.VPN, err error) {
	vpn.OpenVPN, err = readOpenVPN()
	if err != nil {
		return vpn, fmt.Errorf("reading OpenVPN settings: %w", err)
	}

	return vpn, nil
}
