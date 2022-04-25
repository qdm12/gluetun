package privateinternetaccess

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *PIA) BuildConf(connection models.Connection,
	settings settings.OpenVPN) (lines []string) {
	providerSettings := utils.OpenVPNProviderSettings{
		RemoteCertTLS: true,
		RenegDisabled: true,
		AuthUserPass:  true,
	}

	switch *settings.PIAEncPreset {
	case constants.PIAEncryptionPresetNormal:
		providerSettings.Ciphers = []string{constants.AES128cbc}
		providerSettings.Auth = constants.SHA1
		providerSettings.CRLVerify = constants.PiaX509CRLNormal
		providerSettings.CA = constants.PiaCANormal
	case constants.PIAEncryptionPresetNone:
		providerSettings.Ciphers = []string{"none"}
		providerSettings.Auth = "none"
		providerSettings.CRLVerify = constants.PiaX509CRLNormal
		providerSettings.CA = constants.PiaCANormal
	default: // strong
		providerSettings.Ciphers = []string{constants.AES256cbc}
		providerSettings.Auth = constants.SHA256
		providerSettings.CRLVerify = constants.PiaX509CRLStrong
		providerSettings.CA = constants.PiaCAStrong
	}

	return utils.OpenVPNConfig(providerSettings, connection, settings)
}
