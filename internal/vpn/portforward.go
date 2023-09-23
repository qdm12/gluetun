package vpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/portforward/service"
)

func (l *Loop) startPortForwarding(data tunnelUpData) (err error) {
	partialUpdate := service.Settings{
		PortForwarder: data.portForwarder,
		Interface:     data.vpnIntf,
		ServerName:    data.serverName,
		VPNProvider:   data.portForwarder.Name(),
	}
	return l.portForward.UpdateWith(partialUpdate)
}

func (l *Loop) stopPortForwarding(vpnProvider string) (err error) {
	partialUpdate := service.Settings{
		VPNProvider: vpnProvider,
		UserSettings: settings.PortForwarding{
			Enabled: ptrTo(false),
		},
	}
	return l.portForward.UpdateWith(partialUpdate)
}
