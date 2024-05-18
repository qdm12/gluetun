package utils

import (
	"net/http"
	"net/netip"
)

// PortForwardObjects contains fields that may or may not need to be set
// depending on the port forwarding provider code.
type PortForwardObjects struct {
	// Logger is a logger, used by both Private Internet Access and ProtonVPN.
	Logger Logger
	// Gateway is the VPN gateway IP address, used by ProtonVPN.
	Gateway netip.Addr
	// Client is used to query the VPN gateway for Private Internet Access.
	Client *http.Client
	// ServerHostname is used by Private Internet Access for port forwarding.
	ServerHostname string
	// CanPortForward is used by Private Internet Access for port forwarding.
	CanPortForward bool
}

type Routing interface {
	VPNLocalGatewayIP(vpnInterface string) (gateway netip.Addr, err error)
}
