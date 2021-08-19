package configuration

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
)

func (settings *Provider) readFastestvpn(r reader) (err error) {
	settings.Name = constants.Fastestvpn

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.FastestvpnHostnameChoices())
	if err != nil {
		return fmt.Errorf("environment variable SERVER_HOSTNAME: %w", err)
	}

	settings.ServerSelection.Countries, err = r.env.CSVInside("COUNTRY", constants.FastestvpnCountriesChoices())
	if err != nil {
		return fmt.Errorf("environment variable COUNTRY: %w", err)
	}

	return settings.ServerSelection.OpenVPN.readProtocolOnly(r.env)
}
