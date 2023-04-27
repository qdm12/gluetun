package env

import (
	"fmt"
	"net/netip"
	"os"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *Source) readWireguard() (wireguard settings.Wireguard, err error) {
	defer func() {
		err = unsetEnvKeys([]string{"WIREGUARD_PRIVATE_KEY", "WIREGUARD_PRESHARED_KEY"}, err)
	}()
	wireguard.PrivateKey = envToStringPtr("WIREGUARD_PRIVATE_KEY")
	wireguard.PreSharedKey = envToStringPtr("WIREGUARD_PRESHARED_KEY")
	_, wireguard.Interface = s.getEnvWithRetro("VPN_INTERFACE", "WIREGUARD_INTERFACE")
	wireguard.Implementation = os.Getenv("WIREGUARD_IMPLEMENTATION")
	wireguard.Addresses, err = s.readWireguardAddresses()
	if err != nil {
		return wireguard, err // already wrapped
	}
	return wireguard, nil
}

func (s *Source) readWireguardAddresses() (addresses []netip.Prefix, err error) {
	key, addressesCSV := s.getEnvWithRetro("WIREGUARD_ADDRESSES", "WIREGUARD_ADDRESS")
	if addressesCSV == "" {
		return nil, nil
	}

	addressStrings := strings.Split(addressesCSV, ",")
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
