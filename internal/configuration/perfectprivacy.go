package configuration

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
)

func (settings *Provider) readPerfectPrivacy(r reader) (err error) {
	settings.Name = constants.Perfectprivacy
	servers := r.servers.GetPerfectprivacy()

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.PerfectprivacyCityChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable CITY: %w", err)
	}

	return settings.ServerSelection.OpenVPN.readPerfectPrivacy(r)
}

func (settings *OpenVPNSelection) readPerfectPrivacy(r reader) (err error) {
	settings.TCP, err = readOpenVPNProtocol(r)
	if err != nil {
		return err
	}

	portValidation := openvpnPortValidation{
		tcp:        settings.TCP,
		allowedTCP: []uint16{44, 443, 4433},
		allowedUDP: []uint16{44, 443, 4433},
	}
	settings.CustomPort, err = readOpenVPNCustomPort(r, portValidation)
	if err != nil {
		return err
	}

	return nil
}
