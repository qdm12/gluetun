package configuration

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/params"
)

func (settings *Provider) privateinternetaccessLines() (lines []string) {
	if len(settings.ServerSelection.Regions) > 0 {
		lines = append(lines, lastIndent+"Regions: "+commaJoin(settings.ServerSelection.Regions))
	}

	lines = append(lines, lastIndent+"Encryption preset: "+settings.ServerSelection.EncryptionPreset)

	lines = append(lines, lastIndent+"Custom port: "+strconv.Itoa(int(settings.ServerSelection.CustomPort)))

	if settings.PortForwarding.Enabled {
		lines = append(lines, lastIndent+"Port forwarding:")
		for _, line := range settings.PortForwarding.lines() {
			lines = append(lines, indent+line)
		}
	}

	return lines
}

func (settings *Provider) readPrivateInternetAccess(r reader) (err error) {
	settings.Name = constants.PrivateInternetAccess

	settings.ServerSelection.Protocol, err = readProtocol(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	encryptionPreset, err := r.env.Inside("PIA_ENCRYPTION",
		[]string{constants.PIAEncryptionPresetNormal, constants.PIAEncryptionPresetStrong},
		params.RetroKeys([]string{"ENCRYPTION"}, r.onRetroActive),
		params.Default(constants.PIACertificateStrong),
	)
	if err != nil {
		return err
	}
	settings.ServerSelection.EncryptionPreset = encryptionPreset
	settings.ExtraConfigOptions.EncryptionPreset = encryptionPreset

	settings.ServerSelection.Regions, err = r.env.CSVInside("REGION", constants.PIAGeoChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.CustomPort, err = r.env.Port("PORT", params.Default("0"))
	if err != nil {
		return err
	}

	settings.PortForwarding.Enabled, err = r.env.OnOff("PORT_FORWARDING", params.Default("off"))
	if err != nil {
		return err
	}

	if settings.PortForwarding.Enabled {
		filepathStr, err := r.env.Path("PORT_FORWARDING_STATUS_FILE",
			params.Default("/tmp/gluetun/forwarded_port"), params.CaseSensitiveValue())
		if err != nil {
			return err
		}
		settings.PortForwarding.Filepath = models.Filepath(filepathStr)
	}

	return nil
}
