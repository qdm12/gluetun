package vpn

import (
	"context"
	"time"

	"github.com/qdm12/gluetun/internal/models"
)

func (l *Loop) cleanup(ctx context.Context, pfEnabled bool) {
	for _, vpnPort := range l.vpnInputPorts {
		err := l.fw.RemoveAllowedPort(ctx, vpnPort)
		if err != nil {
			l.logger.Error("cannot remove allowed input port from firewall: " + err.Error())
		}
	}

	l.publicip.SetData(models.PublicIP{}) // clear public IP address data

	if pfEnabled {
		const pfTimeout = 100 * time.Millisecond
		err := l.stopPortForwarding(ctx, pfTimeout)
		if err != nil {
			l.logger.Error("cannot stop port forwarding: " + err.Error())
		}
	}
}
