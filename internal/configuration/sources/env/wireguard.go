package env

import (
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (r *Reader) readWireguard() (wireguard settings.Wireguard, err error) {
	defer func() {
		err = unsetEnvKeys([]string{"WIREGUARD_PRIVATE_KEY", "WIREGUARD_PRESHARED_KEY"}, err)
	}()
	wireguard.PrivateKey = envToStringPtr("WIREGUARD_PRIVATE_KEY")
	wireguard.PreSharedKey = envToStringPtr("WIREGUARD_PRESHARED_KEY")
	_, wireguard.Interface = r.getEnvWithRetro("VPN_INTERFACE", "WIREGUARD_INTERFACE")
	wireguard.Addresses, err = r.readWireguardAddresses()
	if err != nil {
		return wireguard, err // already wrapped
	}
	return wireguard, nil
}

func (r *Reader) readWireguardAddresses() (addresses []net.IPNet, err error) {
	key, addressesCSV := r.getEnvWithRetro("WIREGUARD_ADDRESSES", "WIREGUARD_ADDRESS")
	if addressesCSV == "" {
		return nil, nil
	}

	addressStrings := strings.Split(addressesCSV, ",")
	addresses = make([]net.IPNet, len(addressStrings))
	for i, addressString := range addressStrings {
		ip, ipNet, err := net.ParseCIDR(addressString)
		if err != nil {
			return nil, fmt.Errorf("environment variable %s: %w", key, err)
		}
		ipNet.IP = ip
		addresses[i] = *ipNet
	}

	return addresses, nil
}
