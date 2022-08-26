package env

import (
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/govalid/port"
)

func (s *Source) readOpenVPNSelection() (
	selection settings.OpenVPNSelection, err error) {
	confFile := getCleanedEnv("OPENVPN_CUSTOM_CONFIG")
	if confFile != "" {
		selection.ConfFile = &confFile
	}

	selection.TCP, err = s.readOpenVPNProtocol()
	if err != nil {
		return selection, err
	}

	selection.CustomPort, err = s.readOpenVPNCustomPort()
	if err != nil {
		return selection, err
	}

	selection.PIAEncPreset = s.readPIAEncryptionPreset()

	return selection, nil
}

var ErrOpenVPNProtocolNotValid = errors.New("OpenVPN protocol is not valid")

func (s *Source) readOpenVPNProtocol() (tcp *bool, err error) {
	envKey, protocol := s.getEnvWithRetro("OPENVPN_PROTOCOL", "PROTOCOL")

	switch strings.ToLower(protocol) {
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

func (s *Source) readOpenVPNCustomPort() (customPort *uint16, err error) {
	key, value := s.getEnvWithRetro("VPN_ENDPOINT_PORT", "PORT", "OPENVPN_PORT")
	if value == "" {
		return nil, nil //nolint:nilnil
	}

	customPort = new(uint16)
	*customPort, err = port.Validate(value)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s: %w", key, err)
	}

	return customPort, nil
}
