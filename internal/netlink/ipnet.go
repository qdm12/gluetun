package netlink

import (
	"net"

	"github.com/vishvananda/netlink"
)

func NewIPNet(ip net.IP) *net.IPNet {
	return netlink.NewIPNet(ip)
}
