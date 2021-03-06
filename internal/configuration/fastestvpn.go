package configuration

import (
	"github.com/qdm12/gluetun/internal/constants"
)

func (settings *Provider) fastestvpnLines() (lines []string) {
	if len(settings.ServerSelection.Hostnames) > 0 {
		lines = append(lines, lastIndent+"Hostnames: "+commaJoin(settings.ServerSelection.Hostnames))
	}

	if len(settings.ServerSelection.Countries) > 0 {
		lines = append(lines, lastIndent+"Countries: "+commaJoin(settings.ServerSelection.Countries))
	}

	return lines
}

func (settings *Provider) readFastestvpn(r reader) (err error) {
	settings.Name = constants.Fastestvpn

	settings.ServerSelection.Protocol, err = readProtocol(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.FastestvpnHostnameChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.Countries, err = r.env.CSVInside("COUNTRY", constants.FastestvpnCountriesChoices())
	if err != nil {
		return err
	}

	return nil
}
