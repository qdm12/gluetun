package configuration

import (
	"fmt"
	"strconv"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

func (settings *Provider) nordvpnLines() (lines []string) {
	if len(settings.ServerSelection.Regions) > 0 {
		lines = append(lines, lastIndent+"Regions: "+commaJoin(settings.ServerSelection.Regions))
	}

	if len(settings.ServerSelection.Hostnames) > 0 {
		lines = append(lines, lastIndent+"Hostnames: "+commaJoin(settings.ServerSelection.Hostnames))
	}

	if len(settings.ServerSelection.Names) > 0 {
		lines = append(lines, lastIndent+"Names: "+commaJoin(settings.ServerSelection.Hostnames))
	}

	if numbersUint16 := settings.ServerSelection.Numbers; len(numbersUint16) > 0 {
		numbersString := make([]string, len(numbersUint16))
		for i, numberUint16 := range numbersUint16 {
			numbersString[i] = strconv.Itoa(int(numberUint16))
		}
		lines = append(lines, lastIndent+"Numbers: "+commaJoin(numbersString))
	}

	lines = append(lines, settings.ServerSelection.OpenVPN.lines()...)

	return lines
}

func (settings *Provider) readNordvpn(r reader) (err error) {
	settings.Name = constants.Nordvpn

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Regions, err = r.env.CSVInside("REGION", constants.NordvpnRegionChoices())
	if err != nil {
		return fmt.Errorf("environment variable REGION: %w", err)
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.NordvpnHostnameChoices())
	if err != nil {
		return fmt.Errorf("environment variable SERVER_HOSTNAME: %w", err)
	}

	settings.ServerSelection.Names, err = r.env.CSVInside("SERVER_NAME", constants.NordvpnHostnameChoices())
	if err != nil {
		return fmt.Errorf("environment variable SERVER_NAME: %w", err)
	}

	settings.ServerSelection.Numbers, err = readNordVPNServerNumbers(r.env)
	if err != nil {
		return err
	}

	return settings.ServerSelection.OpenVPN.readProtocolOnly(r.env)
}

func readNordVPNServerNumbers(env params.Env) (numbers []uint16, err error) {
	const possiblePortsCount = 65537
	possibilities := make([]string, possiblePortsCount)
	for i := range possibilities {
		possibilities[i] = fmt.Sprintf("%d", i)
	}
	possibilities[65536] = ""
	values, err := env.CSVInside("SERVER_NUMBER", possibilities)
	if err != nil {
		return nil, err
	}
	numbers = make([]uint16, len(values))
	for i := range values {
		n, err := strconv.Atoi(values[i])
		if err != nil {
			return nil, err
		}
		numbers[i] = uint16(n)
	}
	return numbers, nil
}
