package configuration

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
)

func (settings *Provider) readCyberghost(r reader) (err error) {
	settings.Name = constants.Cyberghost
	servers := r.servers.GetCyberghost()

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Groups, err = r.env.CSVInside("CYBERGHOST_GROUP",
		constants.CyberghostGroupChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable CYBERGHOST_GROUP: %w", err)
	}

	settings.ServerSelection.Regions, err = r.env.CSVInside("REGION", constants.CyberghostRegionChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable REGION: %w", err)
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME",
		constants.CyberghostHostnameChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable SERVER_HOSTNAME: %w", err)
	}

	return settings.ServerSelection.OpenVPN.readProtocolAndPort(r.env)
}

func (settings *OpenVPN) readCyberghost(r reader) (err error) {
	settings.ClientKey, err = readClientKey(r)
	if err != nil {
		return fmt.Errorf("%w: %s", errClientKey, err)
	}

	settings.ClientCrt, err = readClientCertificate(r)
	if err != nil {
		return fmt.Errorf("%w: %s", errClientCert, err)
	}

	return nil
}
