package configuration

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
)

func (settings *Provider) readWevpn(r reader) (err error) {
	settings.Name = constants.Wevpn
	servers := r.servers.GetWevpn()

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.WevpnCityChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable CITY: %w", err)
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.WevpnHostnameChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable SERVER_HOSTNAME: %w", err)
	}

	return settings.ServerSelection.OpenVPN.readWevpn(r)
}

func (settings *OpenVPNSelection) readWevpn(r reader) (err error) {
	settings.TCP, err = readOpenVPNProtocol(r)
	if err != nil {
		return err
	}

	validation := openvpnPortValidation{
		tcp:        settings.TCP,
		allowedTCP: []uint16{53, 1195, 1199, 2018},
		allowedUDP: []uint16{80, 1194, 1198},
	}
	settings.CustomPort, err = readOpenVPNCustomPort(r, validation)
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
