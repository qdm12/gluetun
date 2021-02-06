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

	if numbersUint16 := settings.ServerSelection.Numbers; len(numbersUint16) > 0 {
		numbersString := make([]string, len(numbersUint16))
		for i, numberUint16 := range numbersUint16 {
			numbersString[i] = strconv.Itoa(int(numberUint16))
		}
		lines = append(lines, lastIndent+"Numbers: "+commaJoin(numbersString))
	}

	return lines
}

func (settings *Provider) readNordvpn(r reader) (err error) {
	settings.Name = constants.Nordvpn

	settings.ServerSelection.Protocol, err = readProtocol(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Regions, err = r.env.CSVInside("REGION", constants.NordvpnRegionChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.Numbers, err = readNordVPNServerNumbers(r.env)
	if err != nil {
		return err
	}

	return nil
}

func readNordVPNServerNumbers(env params.Env) (numbers []uint16, err error) {
	possibilities := make([]string, 65537)
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
