package privateinternetaccess

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Provider) GetConnection(selection settings.ServerSelection, ipv6Supported bool) (
	connection models.Connection, err error,
) {
	defaults := getConnectionDefaults()
	return utils.GetConnection(p.Name(),
		p.storage, selection, defaults, ipv6Supported, p.randSource)
}

func getConnectionDefaults() utils.ConnectionDefaults {
	return utils.NewConnectionDefaults(8443, 8080, 0) //nolint:mnd
}
