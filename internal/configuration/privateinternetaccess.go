package configuration

import (
	"fmt"
	"strconv"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

func (settings *Provider) privateinternetaccessLines() (lines []string) {
	if len(settings.ServerSelection.Regions) > 0 {
		lines = append(lines, lastIndent+"Regions: "+commaJoin(settings.ServerSelection.Regions))
	}

	if len(settings.ServerSelection.Hostnames) > 0 {
		lines = append(lines, lastIndent+"Hostnames: "+commaJoin(settings.ServerSelection.Hostnames))
	}

	if len(settings.ServerSelection.Names) > 0 {
		lines = append(lines, lastIndent+"Names: "+commaJoin(settings.ServerSelection.Names))
	}

	if settings.ServerSelection.CustomPort > 0 {
		lines = append(lines, lastIndent+"Custom port: "+strconv.Itoa(int(settings.ServerSelection.CustomPort)))
	}

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

	settings.ServerSelection.TCP, err = readProtocol(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Regions, err = r.env.CSVInside("REGION", constants.PIAGeoChoices())
	if err != nil {
		return fmt.Errorf("environment variable REGION: %w", err)
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.PIAHostnameChoices())
	if err != nil {
		return fmt.Errorf("environment variable SERVER_HOSTNAME: %w", err)
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_NAME", constants.PIANameChoices())
	if err != nil {
		return fmt.Errorf("environment variable SERVER_NAME: %w", err)
	}

	settings.ServerSelection.CustomPort, err = readPortOrZero(r.env, "PORT")
	if err != nil {
		return fmt.Errorf("environment variable PORT: %w", err)
	}

	settings.ServerSelection.EncryptionPreset, err = getPIAEncryptionPreset(r)
	if err != nil {
		return err
	}

	settings.PortForwarding.Enabled, err = r.env.OnOff("PORT_FORWARDING", params.Default("off"))
	if err != nil {
		return fmt.Errorf("environment variable PORT_FORWARDING: %w", err)
	}

	if settings.PortForwarding.Enabled {
		settings.PortForwarding.Filepath, err = r.env.Path("PORT_FORWARDING_STATUS_FILE",
			params.Default("/tmp/gluetun/forwarded_port"), params.CaseSensitiveValue())
		if err != nil {
			return fmt.Errorf("environment variable PORT_FORWARDING_STATUS_FILE: %w", err)
		}
	}

	return nil
}

func (settings *OpenVPN) readPrivateInternetAccess(r reader) (err error) {
	settings.EncPreset, err = getPIAEncryptionPreset(r)
	return err
}

func getPIAEncryptionPreset(r reader) (encryptionPreset string, err error) {
	encryptionPreset, err = r.env.Inside("PIA_ENCRYPTION",
		[]string{constants.PIAEncryptionPresetNone, constants.PIAEncryptionPresetNormal, constants.PIAEncryptionPresetStrong},
		params.RetroKeys([]string{"ENCRYPTION"}, r.onRetroActive),
		params.Default(constants.PIACertificateStrong),
	)
	if err != nil {
		return "", fmt.Errorf("environment variable PIA_ENCRYPTION: %w", err)
	}
	return encryptionPreset, nil
}
