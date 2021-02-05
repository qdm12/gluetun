package configuration

import (
	"github.com/qdm12/gluetun/internal/constants"
)

func (settings *Provider) purevpnLines() (lines []string) {
	if len(settings.ServerSelection.Regions) > 0 {
		lines = append(lines, lastIndent+"Regions: "+commaJoin(settings.ServerSelection.Regions))
	}

	if len(settings.ServerSelection.Countries) > 0 {
		lines = append(lines, lastIndent+"Countries: "+commaJoin(settings.ServerSelection.Countries))
	}

	if len(settings.ServerSelection.Cities) > 0 {
		lines = append(lines, lastIndent+"Cities: "+commaJoin(settings.ServerSelection.Cities))
	}

	return lines
}

func (settings *Provider) readPurevpn(r reader) (err error) {
	settings.Name = constants.Purevpn

	settings.ServerSelection.Protocol, err = readProtocol(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Regions, err = r.env.CSVInside("REGION", constants.PurevpnRegionChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.Countries, err = r.env.CSVInside("COUNTRY", constants.PurevpnCountryChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.PurevpnCityChoices())
	if err != nil {
		return err
	}

	return nil
}
