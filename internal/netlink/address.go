package netlink

import "github.com/vishvananda/netlink"

type Addr = netlink.Addr

var _ Addresser = (*NetLink)(nil)

type Addresser interface {
	AddrList(link netlink.Link, family int) (
		addresses []netlink.Addr, err error)
	AddrAdd(link netlink.Link, addr *netlink.Addr) error
}

func (n *NetLink) AddrList(link Link, family int) (
	addresses []Addr, err error) {
	return netlink.AddrList(link, family)
}

func (n *NetLink) AddrAdd(link Link, addr *Addr) error {
	return netlink.AddrAdd(link, addr)
}
