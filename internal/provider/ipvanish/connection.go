package ipvanish

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (i *Ipvanish) GetConnection(selection settings.ServerSelection) (
	connection models.Connection, err error) {
	defaults := utils.NewConnectionDefaults(0, 443, 0) //nolint:gomnd
	return utils.GetConnection(i.servers, selection, defaults, i.randSource)
}
