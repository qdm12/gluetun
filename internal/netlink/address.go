//go:build linux || darwin

package netlink

import (
	"github.com/vishvananda/netlink"
)

func (n *NetLink) AddrList(link Link, family int) (
	addresses []Addr, err error,
) {
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
