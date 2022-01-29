package env

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/govalid/port"
)

func (r *Reader) readOpenVPNSelection() (
	selection settings.OpenVPNSelection, err error) {
	confFile := os.Getenv("OPENVPN_CUSTOM_CONFIG")
	if confFile != "" {
		selection.ConfFile = &confFile
	}

	selection.TCP, err = r.readOpenVPNProtocol()
	if err != nil {
		return selection, err
	}

	selection.CustomPort, err = r.readOpenVPNCustomPort()
	if err != nil {
		return selection, err
	}

	selection.PIAEncPreset = r.readPIAEncryptionPreset()

	return selection, nil
}

var ErrOpenVPNProtocolNotValid = errors.New("OpenVPN protocol is not valid")

func (r *Reader) readOpenVPNProtocol() (tcp *bool, err error) {
		// Retro-compatibility
	envKey := "PROTOCOL"
	protocol := strings.ToLower(os.Getenv("PROTOCOL"))
	if protocol == "" {
		protocol = strings.ToLower(os.Getenv("OPENVPN_PROTOCOL"))
		if protocol != "" {
			envKey = "OPENVPN_PROTOCOL"
		}
	} else {
			r.onRetroActive("PROTOCOL", "OPENVPN_PROTOCOL")
		}

	switch protocol {
	case "":
		return nil, nil //nolint:nilnil
	case constants.UDP:
		return boolPtr(false), nil
	case constants.TCP:
		return boolPtr(true), nil
	default:
		return nil, fmt.Errorf("environment variable %s: %w: %s",
			envKey, ErrOpenVPNProtocolNotValid, protocol)
	}
}

func (r *Reader) readOpenVPNCustomPort() (customPort *uint16, err error) {
	const currentKey = "VPN_ENDPOINT_PORT"
	key := "PORT"
	s := os.Getenv(key) // Retro-compatibility
	if s == "" {
		key = "OPENVPN_PORT" // Retro-compatibility
		s = os.Getenv(key)
		if s == "" {
			key = currentKey
			s = os.Getenv(key)
			if s == "" {
				return nil, nil //nolint:nilnil
			}
		}
	}

	if key != currentKey {
		r.onRetroActive(key, currentKey)
	}

	customPort = new(uint16)
	*customPort, err = port.Validate(s)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s: %w", key, err)
	}

	return customPort, nil
}
