package env

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func readWireguard() (wireguard settings.Wireguard, err error) {
	defer func() {
		err = unsetEnvKeys([]string{"WIREGUARD_PRIVATE_KEY", "WIREGUARD_PRESHARED_KEY"}, err)
	}()
	wireguard.PrivateKey = envToStringPtr("WIREGUARD_PRIVATE_KEY")
	wireguard.PreSharedKey = envToStringPtr("WIREGUARD_PRESHARED_KEY")
	wireguard.Interface = os.Getenv("WIREGUARD_INTERFACE")
	wireguard.Addresses, err = readWireguardAddresses()
	if err != nil {
		return wireguard, err // already wrapped
	}
	return wireguard, nil
}

func readWireguardAddresses() (addresses []*net.IPNet, err error) {
	addressesCSV := os.Getenv("WIREGUARD_ADDRESS")
	addressStrings := strings.Split(addressesCSV, ",")
	addresses = make([]*net.IPNet, len(addressStrings))
	for i, addressString := range addressStrings {
		var ip net.IP
		ip, addresses[i], err = net.ParseCIDR(addressString)
		if err != nil {
			return nil, fmt.Errorf("environment variable WIREGUARD_ADDRESS: address %s: %w",
				addressString, err)
		}
		addresses[i].IP = ip
	}

	return addresses, nil
}
