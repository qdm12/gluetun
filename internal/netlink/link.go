package netlink

import "github.com/vishvananda/netlink"

type (
	Link      = netlink.Link
	Bridge    = netlink.Bridge
	Wireguard = netlink.Wireguard
)

var _ Linker = (*NetLink)(nil)

type Linker interface {
	LinkList() (links []netlink.Link, err error)
	LinkByName(name string) (link netlink.Link, err error)
	LinkByIndex(index int) (link netlink.Link, err error)
	LinkAdd(link netlink.Link) (err error)
	LinkDel(link netlink.Link) (err error)
	LinkSetUp(link netlink.Link) (err error)
	LinkSetDown(link netlink.Link) (err error)
}

func (n *NetLink) LinkList() (links []Link, err error) {
	return netlink.LinkList()
}

func (n *NetLink) LinkByName(name string) (link Link, err error) {
	return netlink.LinkByName(name)
}

func (n *NetLink) LinkByIndex(index int) (link Link, err error) {
	return netlink.LinkByIndex(index)
}

func (n *NetLink) LinkAdd(link Link) (err error) {
	return netlink.LinkAdd(link)
}

func (n *NetLink) LinkDel(link Link) (err error) {
	return netlink.LinkDel(link)
}

func (n *NetLink) LinkSetUp(link Link) (err error) {
	return netlink.LinkSetUp(link)
}

func (n *NetLink) LinkSetDown(link Link) (err error) {
	return netlink.LinkSetDown(link)
}
