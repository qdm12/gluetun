package configuration

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

func (settings *Provider) readMullvad(r reader) (err error) {
	settings.Name = constants.Mullvad

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Countries, err = r.env.CSVInside("COUNTRY", constants.MullvadCountryChoices())
	if err != nil {
		return fmt.Errorf("environment variable COUNTRY: %w", err)
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.MullvadCityChoices())
	if err != nil {
		return fmt.Errorf("environment variable CITY: %w", err)
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.MullvadHostnameChoices())
	if err != nil {
		return fmt.Errorf("environment variable SERVER_HOSTNAME: %w", err)
	}

	settings.ServerSelection.ISPs, err = r.env.CSVInside("ISP", constants.MullvadISPChoices())
	if err != nil {
		return fmt.Errorf("environment variable ISP: %w", err)
	}

	settings.ServerSelection.Owned, err = r.env.YesNo("OWNED", params.Default("no"))
	if err != nil {
		return fmt.Errorf("environment variable OWNED: %w", err)
	}

	err = settings.ServerSelection.OpenVPN.readMullvad(r.env)
	if err != nil {
		return err
	}

	return settings.ServerSelection.Wireguard.readMullvad(r.env)
}

func (settings *OpenVPNSelection) readMullvad(env params.Interface) (err error) {
	settings.TCP, err = readProtocol(env)
	if err != nil {
		return err
	}

	settings.CustomPort, err = readOpenVPNCustomPort(env, settings.TCP,
		[]uint16{80, 443, 1401}, []uint16{53, 1194, 1195, 1196, 1197, 1300, 1301, 1302, 1303, 1400})
	if err != nil {
		return err
	}

	return nil
}

func (settings *WireguardSelection) readMullvad(env params.Interface) (err error) {
	settings.CustomPort, err = readWireguardCustomPort(env, nil)
	if err != nil {
		return err
	}

	return nil
}
