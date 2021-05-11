package protonvpn

import (
	"errors"
	"fmt"
)

func getPort(tcp bool, customPort uint16) (port uint16, err error) {
	if customPort == 0 {
		const defaultTCPPort, defaultUDPPort = 443, 1194
		if tcp {
			return defaultTCPPort, nil
		}
		return defaultUDPPort, nil
	}

	if err := checkPort(customPort, tcp); err != nil {
		return 0, err
	}

	return customPort, nil
}

var ErrInvalidPort = errors.New("invalid port number")

func checkPort(port uint16, tcp bool) (err error) {
	if tcp {
		switch port {
		case 443, 5995, 8443: //nolint:gomnd
			return nil
		default:
			return fmt.Errorf("%w: %d for protocol TCP", ErrInvalidPort, port)
		}
	}
	switch port {
	case 80, 443, 1194, 4569, 5060: //nolint:gomnd
		return nil
	default:
		return fmt.Errorf("%w: %d for protocol UDP", ErrInvalidPort, port)
	}
}
