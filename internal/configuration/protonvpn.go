package configuration

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

func (settings *Provider) readProtonvpn(r reader) (err error) {
	settings.Name = constants.Protonvpn

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Countries, err = r.env.CSVInside("COUNTRY", constants.ProtonvpnCountryChoices())
	if err != nil {
		return fmt.Errorf("environment variable COUNTRY: %w", err)
	}

	settings.ServerSelection.Regions, err = r.env.CSVInside("REGION", constants.ProtonvpnRegionChoices())
	if err != nil {
		return fmt.Errorf("environment variable REGION: %w", err)
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.ProtonvpnCityChoices())
	if err != nil {
		return fmt.Errorf("environment variable CITY: %w", err)
	}

	settings.ServerSelection.Names, err = r.env.CSVInside("SERVER_NAME", constants.ProtonvpnNameChoices())
	if err != nil {
		return fmt.Errorf("environment variable SERVER_NAME: %w", err)
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.ProtonvpnHostnameChoices())
	if err != nil {
		return fmt.Errorf("environment variable SERVER_HOSTNAME: %w", err)
	}

	settings.ServerSelection.FreeOnly, err = r.env.YesNo("FREE_ONLY", params.Default("no"))
	if err != nil {
		return fmt.Errorf("environment variable FREE_ONLY: %w", err)
	}

	return settings.ServerSelection.OpenVPN.readProtocolAndPort(r.env)
}
