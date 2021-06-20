package configuration

import (
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

func (settings *Provider) cyberghostLines() (lines []string) {
	lines = append(lines, lastIndent+"Server group: "+settings.ServerSelection.Group)

	if len(settings.ServerSelection.Regions) > 0 {
		lines = append(lines, lastIndent+"Regions: "+commaJoin(settings.ServerSelection.Regions))
	}

	if len(settings.ServerSelection.Hostnames) > 0 {
		lines = append(lines, lastIndent+"Hostnames: "+commaJoin(settings.ServerSelection.Hostnames))
	}

	if settings.ExtraConfigOptions.ClientKey != "" {
		lines = append(lines, lastIndent+"Client key is set")
	}

	if settings.ExtraConfigOptions.ClientCertificate != "" {
		lines = append(lines, lastIndent+"Client certificate is set")
	}

	return lines
}

func (settings *Provider) readCyberghost(r reader) (err error) {
	settings.Name = constants.Cyberghost

	settings.ServerSelection.TCP, err = readProtocol(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ExtraConfigOptions.ClientKey, err = readClientKey(r)
	if err != nil {
		return err
	}

	settings.ExtraConfigOptions.ClientCertificate, err = readClientCertificate(r)
	if err != nil {
		return err
	}

	settings.ServerSelection.Group, err = r.env.Inside("CYBERGHOST_GROUP",
		constants.CyberghostGroupChoices(), params.Default("Premium UDP Europe"))
	if err != nil {
		return err
	}

	settings.ServerSelection.Regions, err = r.env.CSVInside("REGION", constants.CyberghostRegionChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.CyberghostHostnameChoices())
	if err != nil {
		return err
	}

	return nil
}
