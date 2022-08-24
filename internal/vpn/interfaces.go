package vpn

import (
	"context"
	"net"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/qdm12/gluetun/internal/portforward"
	"github.com/qdm12/gluetun/internal/provider"
)

type Firewall interface {
	SetVPNConnection(ctx context.Context, connection models.Connection, interfaceName string) error
	SetAllowedPort(ctx context.Context, port uint16, interfaceName string) error
	RemoveAllowedPort(ctx context.Context, port uint16) error
}

type Routing interface {
	VPNLocalGatewayIP(vpnInterface string) (gateway net.IP, err error)
}

type PortForward interface {
	Start(ctx context.Context, data portforward.StartData) (outcome string, err error)
	Stop(ctx context.Context) (outcome string, err error)
}

type OpenVPN interface {
	WriteConfig(lines []string) error
	WriteAuthFile(user, password string) error
	WriteAskPassFile(passphrase string) error
}

type Providers interface {
	Get(providerName string) provider.Provider
}

type Storage interface {
	FilterServers(provider string, selection settings.ServerSelection) (servers []models.Server, err error)
	GetServerByName(provider, name string) (server models.Server, ok bool)
}

type NetLinker interface {
	AddrAdd(link netlink.Link, addr *netlink.Addr) error
	Router
	Ruler
	Linker
	IsWireguardSupported() (ok bool, err error)
}

type Router interface {
	RouteList(link netlink.Link, family int) (
		routes []netlink.Route, err error)
	RouteAdd(route *netlink.Route) error
}

type Ruler interface {
	RuleAdd(rule *netlink.Rule) error
	RuleDel(rule *netlink.Rule) error
}

type Linker interface {
	LinkList() (links []netlink.Link, err error)
	LinkByName(name string) (link netlink.Link, err error)
	LinkAdd(link netlink.Link) (err error)
	LinkDel(link netlink.Link) (err error)
	LinkSetUp(link netlink.Link) (err error)
	LinkSetDown(link netlink.Link) (err error)
}

type DNSLoop interface {
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
	GetSettings() (settings settings.DNS)
}

type PublicIPLoop interface {
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
	SetData(data models.PublicIP)
}
