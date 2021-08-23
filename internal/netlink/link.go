package netlink

import "github.com/vishvananda/netlink"

var _ Linker = (*NetLink)(nil)

type Linker interface {
	LinkList() (links []netlink.Link, err error)
	LinkByName(name string) (link netlink.Link, err error)
	LinkByIndex(index int) (link netlink.Link, err error)
}

func (n *NetLink) LinkList() (links []netlink.Link, err error) {
	return netlink.LinkList()
}

func (n *NetLink) LinkByName(name string) (link netlink.Link, err error) {
	return netlink.LinkByName(name)
}

func (n *NetLink) LinkByIndex(index int) (link netlink.Link, err error) {
	return netlink.LinkByIndex(index)
}
