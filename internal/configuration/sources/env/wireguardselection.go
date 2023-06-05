package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readWireguardSelection() (
	selection settings.WireguardSelection, err error) {
	selection.EndpointIP, err = s.env.NetipAddr("VPN_ENDPOINT_IP", env.RetroKeys("WIREGUARD_ENDPOINT_IP"))
	if err != nil {
		return selection, err
	}

	selection.EndpointPort, err = s.env.Uint16Ptr("VPN_ENDPOINT_PORT", env.RetroKeys("WIREGUARD_ENDPOINT_PORT"))
	if err != nil {
		return selection, err
	}

	selection.PublicKey = s.env.String("WIREGUARD_PUBLIC_KEY", env.ForceLowercase(false))

	return selection, nil
}
