package configuration

import (
	"github.com/qdm12/gluetun/internal/constants"
)

func (settings *Provider) surfsharkLines() (lines []string) {
	if len(settings.ServerSelection.Regions) > 0 {
		lines = append(lines, lastIndent+"Regions: "+commaJoin(settings.ServerSelection.Regions))
	}

	if len(settings.ServerSelection.Hostnames) > 0 {
		lines = append(lines, lastIndent+"Hostnames: "+commaJoin(settings.ServerSelection.Hostnames))
	}

	return lines
}

func (settings *Provider) readSurfshark(r reader) (err error) {
	settings.Name = constants.Surfshark

	settings.ServerSelection.Protocol, err = readProtocol(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Regions, err = r.env.CSVInside("REGION", constants.SurfsharkRegionChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.SurfsharkHostnameChoices())
	if err != nil {
		return err
	}

	return nil
}
