package configuration

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
)

func (settings *Provider) cyberghostLines() (lines []string) {
	lines = append(lines, lastIndent+"Server groups: "+commaJoin(settings.ServerSelection.Groups))

	if len(settings.ServerSelection.Regions) > 0 {
		lines = append(lines, lastIndent+"Regions: "+commaJoin(settings.ServerSelection.Regions))
	}

	if len(settings.ServerSelection.Hostnames) > 0 {
		lines = append(lines, lastIndent+"Hostnames: "+commaJoin(settings.ServerSelection.Hostnames))
	}

	lines = append(lines, settings.ServerSelection.OpenVPN.lines()...)

	return lines
}

func (settings *Provider) readCyberghost(r reader) (err error) {
	settings.Name = constants.Cyberghost

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Groups, err = r.env.CSVInside("CYBERGHOST_GROUP",
		constants.CyberghostGroupChoices())
	if err != nil {
		return fmt.Errorf("environment variable CYBERGHOST_GROUP: %w", err)
	}

	settings.ServerSelection.Regions, err = r.env.CSVInside("REGION", constants.CyberghostRegionChoices())
	if err != nil {
		return fmt.Errorf("environment variable REGION: %w", err)
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.CyberghostHostnameChoices())
	if err != nil {
		return fmt.Errorf("environment variable SERVER_HOSTNAME: %w", err)
	}

	return settings.ServerSelection.OpenVPN.readProtocolAndPort(r.env)
}

func (settings *OpenVPN) readCyberghost(r reader) (err error) {
	settings.ClientKey, err = readClientKey(r)
	if err != nil {
		return err
	}

	settings.ClientCrt, err = readClientCertificate(r)
	if err != nil {
		return err
	}

	return nil
}
