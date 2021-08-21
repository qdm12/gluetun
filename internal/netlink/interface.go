package netlink

import "github.com/vishvananda/netlink"

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . NetLinker

var _ NetLinker = (*NetLink)(nil)

type NetLinker interface {
	AddrAdd(link netlink.Link, addr *netlink.Addr) error
	RouteAdd(route *netlink.Route) error
	RuleAdd(rule *netlink.Rule) error
	RuleDel(rule *netlink.Rule) error
}
