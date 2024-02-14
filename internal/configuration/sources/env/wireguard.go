package env

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readWireguard() (wireguard settings.Wireguard, err error) {

	wireguard.PrivateKey, err = s.readSecretFileAsStringPtr(
		"WIREGUARD_PRIVATE_KEY",
		"/run/secrets/wireguard_private_key",
		[]string{},
	)
	if err != nil {
		return wireguard, fmt.Errorf("reading user file: %w", err)
	}

	wireguard.PrivateKey, err = s.readSecretFileAsStringPtr(
		"WIREGUARD_PRESHARED_KEY",
		"/run/secrets/wireguard_preshared_key",
		[]string{},
	)
	if err != nil {
		return wireguard, fmt.Errorf("reading user file: %w", err)
	}

	wireguard.Interface = s.env.String("VPN_INTERFACE",
		env.RetroKeys("WIREGUARD_INTERFACE"), env.ForceLowercase(false))
	wireguard.Implementation = s.env.String("WIREGUARD_IMPLEMENTATION")
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
