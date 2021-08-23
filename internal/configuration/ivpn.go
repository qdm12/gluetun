package configuration

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
)

func (settings *Provider) readIvpn(r reader) (err error) {
	settings.Name = constants.Ivpn

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Countries, err = r.env.CSVInside("COUNTRY", constants.IvpnCountryChoices())
	if err != nil {
		return fmt.Errorf("environment variable COUNTRY: %w", err)
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.IvpnCityChoices())
	if err != nil {
		return fmt.Errorf("environment variable CITY: %w", err)
	}

	settings.ServerSelection.ISPs, err = r.env.CSVInside("ISP", constants.IvpnISPChoices())
	if err != nil {
		return fmt.Errorf("environment variable ISP: %w", err)
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.IvpnHostnameChoices())
	if err != nil {
		return fmt.Errorf("environment variable SERVER_HOSTNAME: %w", err)
	}

	return settings.ServerSelection.OpenVPN.readProtocolOnly(r.env)
}
