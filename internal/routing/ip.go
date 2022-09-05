package routing

import (
	"errors"
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/netlink"
)

func IPIsPrivate(ip net.IP) bool {
	return ip.IsPrivate() || ip.IsLoopback() ||
		ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()
}

var (
	errInterfaceIPNotFound = errors.New("IP address not found for interface")
)

func ipMatchesFamily(ip net.IP, family int) bool {
	return (family == netlink.FAMILY_V6 && ip.To4() == nil) ||
		(family == netlink.FAMILY_V4 && ip.To4() != nil)
}

func (r *Routing) assignedIP(interfaceName string, family int) (ip net.IP, err error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("network interface %s not found: %w", interfaceName, err)
	}
	addresses, err := iface.Addrs()
	if err != nil {
		return nil, fmt.Errorf("cannot list interface %s addresses: %w", interfaceName, err)
	}
	for _, address := range addresses {
		switch value := address.(type) {
		case *net.IPAddr:
			if ipMatchesFamily(value.IP, family) {
				return value.IP, nil
			}
		case *net.IPNet:
			if ipMatchesFamily(value.IP, family) {
				return value.IP, nil
			}
		}
	}
	return nil, fmt.Errorf("%w: interface %s in %d addresses",
		errInterfaceIPNotFound, interfaceName, len(addresses))
}
