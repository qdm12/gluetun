package routing

import (
	"net/netip"
	"sync"

	"github.com/qdm12/gluetun/internal/netlink"
)

type NetLinker interface {
	Addresser
	Router
	Ruler
	Linker
}

type Addresser interface {
	AddrList(linkIndex uint32, family uint8) (
		addresses []netip.Prefix, err error)
	AddrReplace(linkIndex uint32, prefix netip.Prefix) error
}

type Router interface {
	RouteList(family uint8) (routes []netlink.Route, err error)
	RouteAdd(route netlink.Route) error
	RouteDel(route netlink.Route) error
	RouteReplace(route netlink.Route) error
}

type Ruler interface {
	RuleList(family uint8) (rules []netlink.Rule, err error)
	RuleAdd(rule netlink.Rule) error
	RuleDel(rule netlink.Rule) error
}

type Linker interface {
	LinkList() (links []netlink.Link, err error)
	LinkByName(name string) (link netlink.Link, err error)
	LinkByIndex(index uint32) (link netlink.Link, err error)
	LinkAdd(link netlink.Link) (linkIndex uint32, err error)
	LinkDel(index uint32) (err error)
	LinkSetUp(index uint32) (err error)
	LinkSetDown(index uint32) (err error)
}

type Routing struct {
	netLinker       NetLinker
	logger          Logger
	outboundSubnets []netip.Prefix
	stateMutex      sync.RWMutex
}

// New creates a new routing instance.
func New(netLinker NetLinker, logger Logger) *Routing {
	return &Routing{
		netLinker: netLinker,
		logger:    logger,
	}
}
