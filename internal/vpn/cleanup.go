package vpn

import (
	"context"

	"github.com/qdm12/gluetun/internal/models"
)

func (l *Loop) cleanup() {
	for _, vpnPort := range l.vpnInputPorts {
		err := l.fw.RemoveAllowedPort(context.Background(), vpnPort)
		if err != nil {
			l.logger.Error("cannot remove allowed input port from firewall: " + err.Error())
		}
	}

	l.publicip.SetData(models.PublicIP{}) // clear public IP address data

	l.stopPortForwarding()
}
