package privateinternetaccess

import (
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
)

func getPort(openvpnSelection settings.OpenVPNSelection) (
	port uint16, err error) {
	customPort := *openvpnSelection.CustomPort
	tcp := *openvpnSelection.TCP
	if customPort == 0 {
		return getDefaultPort(tcp, *openvpnSelection.PIAEncPreset), nil
	}

	if err := checkPort(customPort, tcp); err != nil {
		return 0, err
	}

	return customPort, nil
}

func getDefaultPort(tcp bool, encryptionPreset string) (port uint16) {
	if tcp {
		switch encryptionPreset {
		case constants.PIAEncryptionPresetNone, constants.PIAEncryptionPresetNormal:
			port = 502
		case constants.PIAEncryptionPresetStrong:
			port = 501
		}
	} else {
		switch encryptionPreset {
		case constants.PIAEncryptionPresetNone, constants.PIAEncryptionPresetNormal:
			port = 1198
		case constants.PIAEncryptionPresetStrong:
			port = 1197
		}
	}
	return port
}

var ErrInvalidPort = errors.New("invalid port number")

func checkPort(port uint16, tcp bool) (err error) {
	if tcp {
		switch port {
		case 80, 110, 443: //nolint:gomnd
			return nil
		default:
			return fmt.Errorf("%w: %d for protocol TCP", ErrInvalidPort, port)
		}
	}
	switch port {
	case 53, 1194, 1197, 1198, 8080, 9201: //nolint:gomnd
		return nil
	default:
		return fmt.Errorf("%w: %d for protocol UDP", ErrInvalidPort, port)
	}
}
