package vpn

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/portforward/service"
)

func (l *Loop) startPortForwarding(data tunnelUpData) (err error) {
	gateway, err := l.routing.VPNLocalGatewayIP(data.vpnIntf)
	if err != nil {
		return fmt.Errorf("obtaining VPN local gateway IP for interface %s: %w", data.vpnIntf, err)
	}
	l.logger.Info("VPN gateway IP address: " + gateway.String())

	partialUpdate := service.Settings{
		PortForwarder: data.portForwarder,
		Gateway:       gateway,
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
