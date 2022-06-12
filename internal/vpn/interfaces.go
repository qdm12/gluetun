package vpn

import (
	"context"
	"net"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/qdm12/gluetun/internal/portforward"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/gluetun/internal/publicip/ipinfo"
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
}

type StateManager interface {
	GetSettings() (vpn settings.VPN)
	SetSettings(ctx context.Context, vpn settings.VPN) (outcome string)
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
	IsWireguardSupported() (ok bool, err error)
	RouteList(link netlink.Link, family int) (
		routes []netlink.Route, err error)
	RouteAdd(route *netlink.Route) error
	RuleAdd(rule *netlink.Rule) error
	RuleDel(rule *netlink.Rule) error
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
	SetData(data ipinfo.Response)
}

type statusManager interface {
	GetStatus() (status models.LoopStatus)
	SetStatus(status models.LoopStatus)
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
	Lock()
	Unlock()
}

type runner interface {
	Run(ctx context.Context, waitError chan<- error, tunnelReady chan<- struct{})
}
