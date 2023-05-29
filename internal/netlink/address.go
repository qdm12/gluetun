package netlink

import (
	"net/netip"

	"github.com/vishvananda/netlink"
)

type Addr struct {
	Network netip.Prefix
}

func (a Addr) String() string {
	return a.Network.String()
}

func (n *NetLink) AddrList(link Link, family int) (
	addresses []Addr, err error) {
	netlinkLink := linkToNetlinkLink(&link)
	netlinkAddresses, err := netlink.AddrList(netlinkLink, family)
	if err != nil {
		return nil, err
	}

	addresses = make([]Addr, len(netlinkAddresses))
	for i := range netlinkAddresses {
		addresses[i].Network = netIPNetToNetipPrefix(netlinkAddresses[i].IPNet)
	}

	return addresses, nil
}

func (n *NetLink) AddrReplace(link Link, addr Addr) error {
	netlinkLink := linkToNetlinkLink(&link)
	netlinkAddress := netlink.Addr{
		IPNet: netipPrefixToIPNet(addr.Network),
	}

	return netlink.AddrReplace(netlinkLink, &netlinkAddress)
}
