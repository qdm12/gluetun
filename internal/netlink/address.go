package netlink

import "github.com/vishvananda/netlink"

func (n *NetLink) AddrAdd(link netlink.Link, addr *netlink.Addr) error {
	return netlink.AddrAdd(link, addr)
}
