package configuration

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

func (settings *Provider) readWindscribe(r reader) (err error) {
	settings.Name = constants.Windscribe
	servers := r.servers.GetWindscribe()

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Regions, err = r.env.CSVInside("REGION", constants.WindscribeRegionChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable REGION: %w", err)
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.WindscribeCityChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable CITY: %w", err)
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME",
		constants.WindscribeHostnameChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable SERVER_HOSTNAME: %w", err)
	}

	err = settings.ServerSelection.OpenVPN.readWindscribe(r)
	if err != nil {
		return err
	}

	return settings.ServerSelection.Wireguard.readWindscribe(r.env)
}

func (settings *OpenVPNSelection) readWindscribe(r reader) (err error) {
	settings.TCP, err = readOpenVPNProtocol(r)
	if err != nil {
		return err
	}

	settings.CustomPort, err = readOpenVPNCustomPort(r, openvpnPortValidation{
		tcp:        settings.TCP,
		allowedTCP: []uint16{21, 22, 80, 123, 143, 443, 587, 1194, 3306, 8080, 54783},
		allowedUDP: []uint16{53, 80, 123, 443, 1194, 54783},
	})
	if err != nil {
		return err
	}

	return nil
}

func (settings *WireguardSelection) readWindscribe(env params.Interface) (err error) {
	settings.EndpointPort, err = readWireguardCustomPort(env,
		[]uint16{53, 80, 123, 443, 1194, 65142})
	if err != nil {
		return err
	}

	return nil
}
