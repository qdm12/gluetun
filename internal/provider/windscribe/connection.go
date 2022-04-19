package windscribe

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Windscribe) GetConnection(selection settings.ServerSelection) (
	connection models.Connection, err error) {
	defaults := utils.NewConnectionDefaults(443, 1194, 1194) //nolint:gomnd
	return utils.GetConnection(p.servers, selection, defaults, p.randSource)
}
