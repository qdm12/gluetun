package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (r *Reader) readVPN() (vpn settings.VPN, err error) {
	vpn.Type = strings.ToLower(os.Getenv("VPN_TYPE"))

	vpn.Provider, err = r.readProvider(vpn.Type)
	if err != nil {
		return vpn, fmt.Errorf("cannot read provider settings: %w", err)
	}

	vpn.OpenVPN, err = r.readOpenVPN()
	if err != nil {
		return vpn, fmt.Errorf("cannot read OpenVPN settings: %w", err)
	}

	vpn.Wireguard, err = r.readWireguard()
	if err != nil {
		return vpn, fmt.Errorf("cannot read Wireguard settings: %w", err)
	}

	return vpn, nil
}
