package configuration

import (
	"github.com/qdm12/gluetun/internal/constants"
)

func (settings *Provider) privadoLines() (lines []string) {
	if len(settings.ServerSelection.Hostnames) > 0 {
		lines = append(lines, lastIndent+"Hostnames: "+commaJoin(settings.ServerSelection.Hostnames))
	}

	return lines
}

func (settings *Provider) readPrivado(r reader) (err error) {
	settings.Name = constants.Privado

	settings.ServerSelection.Protocol, err = readProtocol(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.PrivadoHostnameChoices())
	if err != nil {
		return err
	}

	return nil
}
