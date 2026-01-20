package wireguard

import "github.com/qdm12/gluetun/internal/netlink"

//go:generate mockgen -destination=netlinker_mock_test.go -package wireguard . NetLinker

type NetLinker interface {
	AddrReplace(link netlink.Link, addr netlink.Addr) error
	Router
	Ruler
	Linker
	IsWireguardSupported() bool
}

type Router interface {
	RouteList(family int) (routes []netlink.Route, err error)
	RouteAdd(route netlink.Route) error
}

type Ruler interface {
	RuleAdd(rule netlink.Rule) error
	RuleDel(rule netlink.Rule) error
}

type Linker interface {
	LinkAdd(link netlink.Link) (linkIndex int, err error)
	LinkList() (links []netlink.Link, err error)
	LinkByName(name string) (link netlink.Link, err error)
	LinkSetUp(link netlink.Link) (linkIndex int, err error)
	LinkSetDown(link netlink.Link) error
	LinkDel(link netlink.Link) error
}
