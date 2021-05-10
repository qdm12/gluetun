package configuration

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

func (settings *Provider) mullvadLines() (lines []string) {
	if len(settings.ServerSelection.Countries) > 0 {
		lines = append(lines, lastIndent+"Countries: "+commaJoin(settings.ServerSelection.Countries))
	}

	if len(settings.ServerSelection.Cities) > 0 {
		lines = append(lines, lastIndent+"Cities: "+commaJoin(settings.ServerSelection.Cities))
	}

	if len(settings.ServerSelection.Hostnames) > 0 {
		lines = append(lines, lastIndent+"Hostnames: "+commaJoin(settings.ServerSelection.Hostnames))
	}

	if len(settings.ServerSelection.ISPs) > 0 {
		lines = append(lines, lastIndent+"ISPs: "+commaJoin(settings.ServerSelection.ISPs))
	}

	if settings.ServerSelection.CustomPort > 0 {
		lines = append(lines, lastIndent+"Custom port: "+strconv.Itoa(int(settings.ServerSelection.CustomPort)))
	}

	if settings.ExtraConfigOptions.OpenVPNIPv6 {
		lines = append(lines, lastIndent+"IPv6: enabled")
	}

	return lines
}

func (settings *Provider) readMullvad(r reader) (err error) {
	settings.Name = constants.Mullvad

	settings.ServerSelection.TCP, err = readProtocol(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Countries, err = r.env.CSVInside("COUNTRY", constants.MullvadCountryChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.MullvadCityChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME", constants.MullvadHostnameChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.ISPs, err = r.env.CSVInside("ISP", constants.MullvadISPChoices())
	if err != nil {
		return err
	}

	settings.ServerSelection.CustomPort, err = readCustomPort(r.env, settings.ServerSelection.TCP,
		[]uint16{80, 443, 1401}, []uint16{53, 1194, 1195, 1196, 1197, 1300, 1301, 1302, 1303, 1400})
	if err != nil {
		return err
	}

	settings.ServerSelection.Owned, err = r.env.YesNo("OWNED", params.Default("no"))
	if err != nil {
		return err
	}

	settings.ExtraConfigOptions.OpenVPNIPv6, err = r.env.OnOff("OPENVPN_IPV6", params.Default("off"))
	if err != nil {
		return err
	}

	return nil
}
