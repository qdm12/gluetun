package configuration

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
)

func (settings *Provider) readPrivatevpn(r reader) (err error) {
	settings.Name = constants.Privatevpn
	servers := r.servers.GetPrivatevpn()

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Countries, err = r.env.CSVInside("COUNTRY", constants.PrivatevpnCountryChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable COUNTRY: %w", err)
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.PrivatevpnCityChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable CITY: %w", err)
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME",
		constants.PrivatevpnHostnameChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable SERVER_HOSTNAME: %w", err)
	}

	return settings.ServerSelection.OpenVPN.readProtocolAndPort(r.env)
}
