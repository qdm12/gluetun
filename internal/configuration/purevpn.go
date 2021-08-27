package configuration

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
)

func (settings *Provider) readPurevpn(r reader) (err error) {
	settings.Name = constants.Purevpn
	servers := r.servers.GetPurevpn()

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Regions, err = r.env.CSVInside("REGION", constants.PurevpnRegionChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable REGION: %w", err)
	}

	settings.ServerSelection.Countries, err = r.env.CSVInside("COUNTRY", constants.PurevpnCountryChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable COUNTRY: %w", err)
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.PurevpnCityChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable CITY: %w", err)
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.PurevpnHostnameChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable SERVER_HOSTNAME: %w", err)
	}

	return settings.ServerSelection.OpenVPN.readProtocolOnly(r.env)
}
