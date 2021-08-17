package configuration

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

func (settings *Provider) vpnUnlimitedLines() (lines []string) {
	if len(settings.ServerSelection.Countries) > 0 {
		lines = append(lines, lastIndent+"Countries: "+commaJoin(settings.ServerSelection.Countries))
	}

	if len(settings.ServerSelection.Cities) > 0 {
		lines = append(lines, lastIndent+"Cities: "+commaJoin(settings.ServerSelection.Cities))
	}

	if len(settings.ServerSelection.Hostnames) > 0 {
		lines = append(lines, lastIndent+"Hostnames: "+commaJoin(settings.ServerSelection.Hostnames))
	}

	if settings.ServerSelection.FreeOnly {
		lines = append(lines, lastIndent+"Free servers only")
	}

	if settings.ServerSelection.StreamOnly {
		lines = append(lines, lastIndent+"Stream servers only")
	}

	return lines
}

func (settings *Provider) readVPNUnlimited(r reader) (err error) {
	settings.Name = constants.VPNUnlimited

	settings.ServerSelection.TCP, err = readProtocol(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Countries, err = r.env.CSVInside("COUNTRY", constants.VPNUnlimitedCountryChoices())
	if err != nil {
		return fmt.Errorf("environment variable COUNTRY: %w", err)
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.VPNUnlimitedCityChoices())
	if err != nil {
		return fmt.Errorf("environment variable CITY: %w", err)
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.VPNUnlimitedHostnameChoices())
	if err != nil {
		return fmt.Errorf("environment variable SERVER_HOSTNAME: %w", err)
	}

	settings.ServerSelection.FreeOnly, err = r.env.YesNo("FREE_ONLY", params.Default("no"))
	if err != nil {
		return fmt.Errorf("environment variable FREE_ONLY: %w", err)
	}

	settings.ServerSelection.StreamOnly, err = r.env.YesNo("STREAM_ONLY", params.Default("no"))
	if err != nil {
		return fmt.Errorf("environment variable STREAM_ONLY: %w", err)
	}

	return nil
}

func (settings *OpenVPN) readVPNUnlimited(r reader) (err error) {
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
