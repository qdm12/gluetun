package wireguard

import "github.com/vishvananda/netlink"

//go:generate mockgen -destination=netlinker_mock_test.go -package wireguard . NetLinker

type NetLinker interface {
	AddrAdd(link netlink.Link, addr *netlink.Addr) error
	RouteAdd(route *netlink.Route) error
	RuleAdd(rule *netlink.Rule) error
	RuleDel(rule *netlink.Rule) error
}
