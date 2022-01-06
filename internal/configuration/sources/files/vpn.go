package files

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (r *Reader) readVPN() (vpn settings.VPN, err error) {
	vpn.OpenVPN, err = r.readOpenVPN()
	if err != nil {
		return vpn, fmt.Errorf("cannot read OpenVPN settings: %w", err)
	}

	return vpn, nil
}
