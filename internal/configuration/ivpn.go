package configuration

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

func (settings *Provider) readIvpn(r reader) (err error) {
	settings.Name = constants.Ivpn
	servers := r.servers.GetIvpn()

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Countries, err = r.env.CSVInside("COUNTRY", constants.IvpnCountryChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable COUNTRY: %w", err)
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.IvpnCityChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable CITY: %w", err)
	}

	settings.ServerSelection.ISPs, err = r.env.CSVInside("ISP", constants.IvpnISPChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable ISP: %w", err)
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.IvpnHostnameChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable SERVER_HOSTNAME: %w", err)
	}

	err = settings.ServerSelection.OpenVPN.readIVPN(r)
	if err != nil {
		return err
	}

	return settings.ServerSelection.Wireguard.readIVPN(r.env)
}

func (settings *OpenVPNSelection) readIVPN(r reader) (err error) {
	settings.TCP, err = readOpenVPNProtocol(r)
	if err != nil {
		return err
	}

	settings.CustomPort, err = readOpenVPNCustomPort(r.env, settings.TCP,
		[]uint16{80, 443, 1443}, []uint16{53, 1194, 2049, 2050})
	if err != nil {
		return err
	}

	return nil
}

func (settings *WireguardSelection) readIVPN(env params.Interface) (err error) {
	settings.EndpointPort, err = readWireguardCustomPort(env,
		[]uint16{2049, 2050, 53, 30587, 41893, 48574, 58237})
	if err != nil {
		return err
	}

	return nil
}
