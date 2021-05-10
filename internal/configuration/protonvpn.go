package configuration

import (
	"github.com/qdm12/gluetun/internal/constants"
)

func (settings *Provider) protonvpnLines() (lines []string) {
	if len(settings.ServerSelection.Countries) > 0 {
		lines = append(lines, lastIndent+"Countries: "+commaJoin(settings.ServerSelection.Countries))
	}

	if len(settings.ServerSelection.Regions) > 0 {
		lines = append(lines, lastIndent+"Regions: "+commaJoin(settings.ServerSelection.Regions))
	}

	if len(settings.ServerSelection.Cities) > 0 {
		lines = append(lines, lastIndent+"Cities: "+commaJoin(settings.ServerSelection.Cities))
	}

	if len(settings.ServerSelection.Names) > 0 {
		lines = append(lines, lastIndent+"Names: "+commaJoin(settings.ServerSelection.Names))
	}

	if len(settings.ServerSelection.Hostnames) > 0 {
		lines = append(lines, lastIndent+"Hostnames: "+commaJoin(settings.ServerSelection.Hostnames))
	}

	return lines
}

func (settings *Provider) readProtonvpn(r reader) (err error) {
	settings.Name = constants.Protonvpn

	settings.ServerSelection.TCP, err = readProtocol(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.CustomPort, err = readPortOrZero(r.env, "PORT")
	if err != nil {
		return err
	}

	settings.ServerSelection.Countries, err = r.env.CSVInside("COUNTRY", constants.ProtonvpnCountryChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.Regions, err = r.env.CSVInside("REGION", constants.ProtonvpnRegionChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.ProtonvpnCityChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.Names, err = r.env.CSVInside("SERVER_NAME", constants.ProtonvpnNameChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.ProtonvpnHostnameChoices())
	if err != nil {
		return err
	}

	return nil
}
