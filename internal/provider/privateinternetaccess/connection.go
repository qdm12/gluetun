package privateinternetaccess

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *PIA) GetConnection(selection settings.ServerSelection) (
	connection models.Connection, err error) {
	// Set port defaults depending on encryption preset.
	var defaults utils.ConnectionDefaults
	switch *selection.OpenVPN.PIAEncPreset {
	case constants.PIAEncryptionPresetNone, constants.PIAEncryptionPresetNormal:
		defaults.OpenVPNTCPPort = 502
		defaults.OpenVPNUDPPort = 1198
	case constants.PIAEncryptionPresetStrong:
		defaults.OpenVPNTCPPort = 501
		defaults.OpenVPNUDPPort = 1197
	}

	return utils.GetConnection(p.servers, selection, defaults, p.randSource)
}
