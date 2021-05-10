package configuration

import (
	"github.com/qdm12/gluetun/internal/constants"
)

func (settings *Provider) vyprvpnLines() (lines []string) {
	if len(settings.ServerSelection.Regions) > 0 {
		lines = append(lines, lastIndent+"Regions: "+commaJoin(settings.ServerSelection.Regions))
	}

	return lines
}

func (settings *Provider) readVyprvpn(r reader) (err error) {
	settings.Name = constants.Vyprvpn

	settings.ServerSelection.TCP, err = readProtocol(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Regions, err = r.env.CSVInside("REGION", constants.VyprvpnRegionChoices())
	if err != nil {
		return err
	}

	return nil
}
