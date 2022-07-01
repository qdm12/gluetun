package routing

import (
	"net"
	"sync"

	"github.com/qdm12/gluetun/internal/netlink"
)

type NetLinker interface {
	AddrList(link netlink.Link, family int) (
		addresses []netlink.Addr, err error)
	AddrAdd(link netlink.Link, addr *netlink.Addr) error
	IsWireguardSupported() (ok bool, err error)
	RouteList(link netlink.Link, family int) (
		routes []netlink.Route, err error)
	RouteAdd(route *netlink.Route) error
	RouteDel(route *netlink.Route) error
	RouteReplace(route *netlink.Route) error
	RuleList(family int) (rules []netlink.Rule, err error)
	RuleAdd(rule *netlink.Rule) error
	RuleDel(rule *netlink.Rule) error
	LinkList() (links []netlink.Link, err error)
	LinkByName(name string) (link netlink.Link, err error)
	LinkByIndex(index int) (link netlink.Link, err error)
	LinkAdd(link netlink.Link) (err error)
	LinkDel(link netlink.Link) (err error)
	LinkSetUp(link netlink.Link) (err error)
	LinkSetDown(link netlink.Link) (err error)
}

type Routing struct {
	netLinker       NetLinker
	logger          Logger
	outboundSubnets []net.IPNet
	stateMutex      sync.RWMutex
}

// New creates a new routing instance.
func New(netLinker NetLinker, logger Logger) *Routing {
	return &Routing{
		netLinker: netLinker,
		logger:    logger,
	}
}
