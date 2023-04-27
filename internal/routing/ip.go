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

func ensureNoIPv6WrappedIPv4(candidateIP net.IP) (resultIP net.IP) {
	const ipv4Size = 4
	if candidateIP.To4() == nil || len(candidateIP) == ipv4Size { // ipv6 or ipv4
		return candidateIP
	}

	// ipv6-wrapped ipv4
	resultIP = make(net.IP, ipv4Size)
	copy(resultIP, candidateIP[12:16])
	return resultIP
}

func (r *Routing) assignedIP(interfaceName string, family int) (ip net.IP, err error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("network interface %s not found: %w", interfaceName, err)
	}
	addresses, err := iface.Addrs()
	if err != nil {
		return nil, fmt.Errorf("listing interface %s addresses: %w", interfaceName, err)
	}
	for _, address := range addresses {
		switch value := address.(type) {
		case *net.IPAddr:
			ip = value.IP
		case *net.IPNet:
			ip = value.IP
		default:
			continue
		}

		if !ipMatchesFamily(ip, family) {
			continue
		}

		// Ensure we don't return an IPv6-wrapped IPv4 address
		// since netip.Address String method works differently than
		// net.IP String method for this kind of addresses.
		ip = ensureNoIPv6WrappedIPv4(ip)
		return ip, nil
	}
	return nil, fmt.Errorf("%w: interface %s in %d addresses",
		errInterfaceIPNotFound, interfaceName, len(addresses))
}
