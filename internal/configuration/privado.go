package configuration

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
)

func (settings *Provider) readPrivado(r reader) (err error) {
	settings.Name = constants.Privado
	servers := r.servers.GetPrivado()

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Countries, err = r.env.CSVInside("COUNTRY", constants.PrivadoCountryChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable COUNTRY: %w", err)
	}

	settings.ServerSelection.Regions, err = r.env.CSVInside("REGION", constants.PrivadoRegionChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable REGION: %w", err)
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.PrivadoCityChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable CITY: %w", err)
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.PrivadoHostnameChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable SERVER_HOSTNAME: %w", err)
	}

	return nil
}
