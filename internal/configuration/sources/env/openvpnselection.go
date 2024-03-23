package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readOpenVPNSelection() (
	selection settings.OpenVPNSelection, err error) {
	selection.ConfFile = s.env.Get("OPENVPN_CUSTOM_CONFIG", env.ForceLowercase(false))

	selection.Protocol = s.env.String("OPENVPN_PROTOCOL", env.RetroKeys("PROTOCOL"))
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
