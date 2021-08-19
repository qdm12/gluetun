package configuration

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

func (settings *Provider) readPrivateInternetAccess(r reader) (err error) {
	settings.Name = constants.PrivateInternetAccess

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

	return settings.ServerSelection.OpenVPN.readPrivateInternetAccess(r)
}

func (settings *OpenVPNSelection) readPrivateInternetAccess(r reader) (err error) {
	settings.EncPreset, err = getPIAEncryptionPreset(r)
	if err != nil {
		return err
	}

	settings.CustomPort, err = readPortOrZero(r.env, "PORT")
	if err != nil {
		return fmt.Errorf("environment variable PORT: %w", err)
	}

	return nil
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
