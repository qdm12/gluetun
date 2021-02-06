package configuration

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

func (settings *Provider) windscribeLines() (lines []string) {
	if len(settings.ServerSelection.Regions) > 0 {
		lines = append(lines, lastIndent+"Regions: "+commaJoin(settings.ServerSelection.Regions))
	}

	if len(settings.ServerSelection.Cities) > 0 {
		lines = append(lines, lastIndent+"Cities: "+commaJoin(settings.ServerSelection.Cities))
	}

	if len(settings.ServerSelection.Hostnames) > 0 {
		lines = append(lines, lastIndent+"Hostnames: "+commaJoin(settings.ServerSelection.Hostnames))
	}

	lines = append(lines, lastIndent+"Custom port: "+strconv.Itoa(int(settings.ServerSelection.CustomPort)))

	return lines
}

func (settings *Provider) readWindscribe(r reader) (err error) {
	settings.Name = constants.Windscribe

	settings.ServerSelection.Protocol, err = readProtocol(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Regions, err = r.env.CSVInside("REGION", constants.WindscribeRegionChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.WindscribeCityChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME",
		constants.WindscribeHostnameChoices(), params.RetroKeys([]string{"HOSTNAME"}, r.onRetroActive))
	if err != nil {
		return err
	}

	settings.ServerSelection.CustomPort, err = readCustomPort(r.env, settings.ServerSelection.Protocol,
		[]uint16{21, 22, 80, 123, 143, 443, 587, 1194, 3306, 8080, 54783},
		[]uint16{53, 80, 123, 443, 1194, 54783})
	if err != nil {
		return err
	}

	return nil
}
