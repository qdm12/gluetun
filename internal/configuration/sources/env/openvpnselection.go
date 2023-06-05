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
	selection.ConfFile = s.env.Get("OPENVPN_CUSTOM_CONFIG", env.ForceLowercase(false))

	selection.TCP, err = s.readOpenVPNProtocol()
	if err != nil {
		return selection, err
	}

	selection.CustomPort, err = s.env.Uint16Ptr("VPN_ENDPOINT_PORT",
		env.RetroKeys("PORT", "OPENVPN_PORT"))
	if err != nil {
		return selection, err
	}

	selection.PIAEncPreset = s.readPIAEncryptionPreset()

	return selection, nil
}

var ErrOpenVPNProtocolNotValid = errors.New("OpenVPN protocol is not valid")

func (s *Source) readOpenVPNProtocol() (tcp *bool, err error) {
	const currentKey = "OPENVPN_PROTOCOL"
	envKey := firstKeySet(s.env, "PROTOCOL", currentKey)
	switch envKey {
	case "":
		return nil, nil //nolint:nilnil
	case currentKey:
	default: // Retro compatibility
		s.handleDeprecatedKey(envKey, currentKey)
	}

	protocol := s.env.String(envKey)
	switch strings.ToLower(protocol) {
	case constants.UDP:
		return ptrTo(false), nil
	case constants.TCP:
		return ptrTo(true), nil
	default:
		return nil, fmt.Errorf("environment variable %s: %w: %s",
			envKey, ErrOpenVPNProtocolNotValid, protocol)
	}
}
