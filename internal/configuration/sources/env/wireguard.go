package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readWireguard() (wireguard settings.Wireguard, err error) {
	defer func() {
		err = unsetEnvKeys([]string{"WIREGUARD_PRIVATE_KEY", "WIREGUARD_PRESHARED_KEY"}, err)
	}()
	wireguard.PrivateKey = s.env.Get("WIREGUARD_PRIVATE_KEY", env.ForceLowercase(false))
	wireguard.PreSharedKey = s.env.Get("WIREGUARD_PRESHARED_KEY", env.ForceLowercase(false))
	wireguard.Interface = s.env.String("VPN_INTERFACE",
		env.RetroKeys("WIREGUARD_INTERFACE"), env.ForceLowercase(false))
	wireguard.Implementation = s.env.String("WIREGUARD_IMPLEMENTATION")
	wireguard.PersistenKeepAlive, err = s.env.Int("WIREGUARD_IMPLEMENTATION", env.ForceLowercase(false))
	if err != nil {
		return wireguard, err // already wrapped
	}
	wireguard.Addresses, err = s.env.CSVNetipPrefixes("WIREGUARD_ADDRESSES",
		env.RetroKeys("WIREGUARD_ADDRESS"))
	if err != nil {
		return wireguard, err // already wrapped
	}
	wireguard.AllowedIPs, err = s.env.CSVNetipPrefixes("WIREGUARD_ALLOWED_IPS")
	if err != nil {
		return wireguard, err // already wrapped
	}
	mtuPtr, err := s.env.Uint16Ptr("WIREGUARD_MTU")
	if err != nil {
		return wireguard, err
	} else if mtuPtr != nil {
		wireguard.MTU = *mtuPtr
	}
	return wireguard, nil
}
