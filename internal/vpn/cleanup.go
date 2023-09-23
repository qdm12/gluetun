package vpn

import (
	"context"
	"errors"

	"github.com/qdm12/gluetun/internal/models"
)

func (l *Loop) cleanup(vpnProvider string) {
	for _, vpnPort := range l.vpnInputPorts {
		err := l.fw.RemoveAllowedPort(context.Background(), vpnPort)
		if err != nil {
			l.logger.Error("cannot remove allowed input port from firewall: " + err.Error())
		}
	}

	l.publicip.SetData(models.PublicIP{}) // clear public IP address data

	err := l.stopPortForwarding(vpnProvider)
	if err != nil {
		portForwardingAlreadyStopped := errors.Is(err, context.Canceled)
		if !portForwardingAlreadyStopped {
			l.logger.Error("stopping port forwarding: " + err.Error())
		}
	}
}
