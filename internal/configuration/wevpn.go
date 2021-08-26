package configuration

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

func (settings *Provider) readWevpn(r reader) (err error) {
	settings.Name = constants.Wevpn

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.WevpnCityChoices())
	if err != nil {
		return fmt.Errorf("environment variable CITY: %w", err)
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.WevpnHostnameChoices())
	if err != nil {
		return fmt.Errorf("environment variable SERVER_HOSTNAME: %w", err)
	}

	return settings.ServerSelection.OpenVPN.readWevpn(r.env)
}

func (settings *OpenVPNSelection) readWevpn(env params.Interface) (err error) {
	settings.TCP, err = readProtocol(env)
	if err != nil {
		return err
	}

	settings.CustomPort, err = readOpenVPNCustomPort(env, settings.TCP,
		[]uint16{53, 1195, 1199, 2018},
		[]uint16{80, 1194, 1198})
	if err != nil {
		return err
	}

	return nil
}

func (settings *OpenVPN) readWevpn(r reader) (err error) {
	settings.ClientKey, err = readClientKey(r)
	if err != nil {
		return err
	}

	return nil
}
