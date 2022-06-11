package netlink

import "github.com/vishvananda/netlink"

type (
	Link      = netlink.Link
	Bridge    = netlink.Bridge
	Wireguard = netlink.Wireguard
)

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
