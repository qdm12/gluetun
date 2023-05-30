package env

import (
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readOpenVPNSelection() (
	selection settings.OpenVPNSelection, err error) {
	selection.ConfFile = env.Get("OPENVPN_CUSTOM_CONFIG", env.ForceLowercase(false))

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
	envKey, protocolPtr := s.getEnvWithRetro("OPENVPN_PROTOCOL", []string{"PROTOCOL"})
	if protocolPtr == nil {
		return nil, nil //nolint:nilnil
	}
	protocol := *protocolPtr

	switch strings.ToLower(protocol) {
	case "":
		return nil, nil //nolint:nilnil
	case constants.UDP:
		return ptrTo(false), nil
	case constants.TCP:
		return ptrTo(true), nil
	default:
		return nil, fmt.Errorf("environment variable %s: %w: %s",
			envKey, ErrOpenVPNProtocolNotValid, protocol)
	}
}

func (s *Source) readOpenVPNCustomPort() (customPort *uint16, err error) {
	key, _ := s.getEnvWithRetro("VPN_ENDPOINT_PORT", []string{"PORT", "OPENVPN_PORT"})
	return env.Uint16Ptr(key)
}
