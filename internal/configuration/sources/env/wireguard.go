package env

import (
	"fmt"
	"net/netip"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readWireguard() (wireguard settings.Wireguard, err error) {
	defer func() {
		err = unsetEnvKeys([]string{"WIREGUARD_PRIVATE_KEY", "WIREGUARD_PRESHARED_KEY"}, err)
	}()
	wireguard.PrivateKey = s.env.Get("WIREGUARD_PRIVATE_KEY", env.ForceLowercase(false))
	wireguard.PreSharedKey = s.env.Get("WIREGUARD_PRESHARED_KEY", env.ForceLowercase(false))
	envKey, _ := s.getEnvWithRetro("VPN_INTERFACE",
		[]string{"WIREGUARD_INTERFACE"}, env.ForceLowercase(false))
	wireguard.Interface = s.env.String(envKey)
	wireguard.Implementation = s.env.String("WIREGUARD_IMPLEMENTATION")
	wireguard.Addresses, err = s.readWireguardAddresses()
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

func (s *Source) readWireguardAddresses() (addresses []netip.Prefix, err error) {
	key, value := s.getEnvWithRetro("WIREGUARD_ADDRESSES",
		[]string{"WIREGUARD_ADDRESS"})
	if value == nil {
		return nil, nil
	}

	addressStrings := strings.Split(*value, ",")
	addresses = make([]netip.Prefix, len(addressStrings))
	for i, addressString := range addressStrings {
		addressString = strings.TrimSpace(addressString)
		addresses[i], err = netip.ParsePrefix(addressString)
		if err != nil {
			return nil, fmt.Errorf("environment variable %s: %w", key, err)
		}
	}

	return addresses, nil
}
