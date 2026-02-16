package mullvad

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
)

func (p *Provider) OpenVPNConfig(_ models.Connection, _ settings.OpenVPN, _ bool) (lines []string) {
	panic("OpenVPN is no longer supported by Mullvad as of January 15th, 2026: " +
		"https://mullvad.net/en/blog/removing-openvpn-15th-january-2026")
}
