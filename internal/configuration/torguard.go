package configuration

import (
	"github.com/qdm12/gluetun/internal/constants"
)

func (settings *Provider) torguardLines() (lines []string) {
	if len(settings.ServerSelection.Countries) > 0 {
		lines = append(lines, lastIndent+"Countries: "+commaJoin(settings.ServerSelection.Countries))
	}

	if len(settings.ServerSelection.Cities) > 0 {
		lines = append(lines, lastIndent+"Cities: "+commaJoin(settings.ServerSelection.Cities))
	}

	if len(settings.ServerSelection.Hostnames) > 0 {
		lines = append(lines, lastIndent+"Hostnames: "+commaJoin(settings.ServerSelection.Hostnames))
	}

	return lines
}

func (settings *Provider) readTorguard(r reader) (err error) {
	settings.Name = constants.Torguard

	settings.ServerSelection.Protocol, err = readProtocol(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Countries, err = r.env.CSVInside("COUNTRY", constants.TorguardCountryChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.TorguardCityChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.TorguardHostnamesChoices())
	if err != nil {
		return err
	}

	return nil
}
