package netlink

import (
	"fmt"
	"net"
	"net/netip"

	"github.com/jsimonetti/rtnetlink/rtnl"
)

func (n *NetLink) AddrList(linkIndex uint32, family uint8) (
	ipPrefixes []netip.Prefix, err error,
) {
	conn, err := rtnl.Dial(nil)
	if err != nil {
		return nil, fmt.Errorf("dialing netlink: %w", err)
	}
	defer conn.Close()

	ifc := &net.Interface{
		Index: int(linkIndex),
	}
	ipNets, err := conn.Addrs(ifc, int(family))
	if err != nil {
		return nil, fmt.Errorf("failed to list addresses: %w", err)
	}

	ipPrefixes = make([]netip.Prefix, len(ipNets))
	for i := range ipNets {
		ipPrefixes[i] = netIPNetToNetipPrefix(ipNets[i])
	}

	return ipPrefixes, nil
}

func (n *NetLink) AddrReplace(linkIndex uint32, prefix netip.Prefix) error {
	conn, err := rtnl.Dial(nil)
	if err != nil {
		return fmt.Errorf("dialing netlink: %w", err)
	}
	defer conn.Close()

	ipNet := netipPrefixToIPNet(prefix)

	// Remove any address identical to the one we want to add
	family := FamilyV4
	if prefix.Addr().Is6() {
		family = FamilyV6
	}
	ifc := &net.Interface{
		Index: int(linkIndex),
	}
	addresses, err := conn.Addrs(ifc, int(family))
	if err != nil {
		return fmt.Errorf("listing addresses: %w", err)
	}
	for _, address := range addresses {
		if address.IP.Equal(ipNet.IP) &&
			net.IP(address.Mask).String() == net.IP(ipNet.Mask).String() {
			err = conn.AddrDel(ifc, address)
			if err != nil {
				return fmt.Errorf("deleting address from interface: %w", err)
			}
			break
		}
	}

	// Add the new address to the interface
	err = conn.AddrAdd(ifc, ipNet)
	if err != nil {
		return fmt.Errorf("adding address to interface: %w", err)
	}

	return nil
}
