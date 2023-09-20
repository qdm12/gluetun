package vpn

import (
	"fmt"
)

func (l *Loop) startPortForwarding(data tunnelUpData) (err error) {
	if !data.portForwarding {
		return nil
	}

	// only used for PIA for now
	gateway, err := l.routing.VPNLocalGatewayIP(data.vpnIntf)
	if err != nil {
		return fmt.Errorf("obtaining VPN local gateway IP for interface %s: %w", data.vpnIntf, err)
	}
	l.logger.Info("VPN gateway IP address: " + gateway.String())

	settings := l.portForward.GetSettings()
	settings.PortForwarder = data.portForwarder
	settings.Gateway = gateway
	settings.ServerName = data.serverName
	settings.Interface = data.vpnIntf
	l.portForward.Update(settings)

	return nil
}

func (l *Loop) stopPortForwarding() {
	settings := l.portForward.GetSettings()
	settings.Settings.Enabled = ptrTo(false)
	l.portForward.Update(settings)
}
