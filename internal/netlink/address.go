package netlink

import "github.com/vishvananda/netlink"

var _ Addresser = (*NetLink)(nil)

type Addresser interface {
	AddrAdd(link netlink.Link, addr *netlink.Addr) error
}

func (n *NetLink) AddrAdd(link netlink.Link, addr *netlink.Addr) error {
	return netlink.AddrAdd(link, addr)
}
