package hidemyass

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (h *HideMyAss) GetConnection(selection settings.ServerSelection) (
	connection models.Connection, err error) {
	defaults := utils.NewConnectionDefaults(8080, 553, 0) //nolint:gomnd
	return utils.GetConnection(h.servers, selection, defaults, h.randSource)
}
