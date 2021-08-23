package netlink

import "github.com/vishvananda/netlink"

type LinkAttrs = netlink.LinkAttrs

func NewLinkAttrs() LinkAttrs {
	return netlink.NewLinkAttrs()
}
