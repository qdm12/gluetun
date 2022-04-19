package mullvad

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (m *Mullvad) GetConnection(selection settings.ServerSelection) (
	connection models.Connection, err error) {
	defaults := utils.NewConnectionDefaults(443, 1194, 51820) //nolint:gomnd
	return utils.GetConnection(m.servers, selection, defaults, m.randSource)
}
