package service

import (
	"net/netip"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/provider"
)

type Settings struct {
	Settings      settings.PortForwarding
	PortForwarder provider.PortForwarder
	Gateway       netip.Addr // needed for PIA
	ServerName    string     // needed for PIA
	Interface     string     // tun0 for example
}
