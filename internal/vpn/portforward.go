package vpn

import (
	"context"
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/portforward/service"
	pfutils "github.com/qdm12/gluetun/internal/provider/utils"
)

func getPortForwarder(provider Provider, providers Providers, //nolint:ireturn
	customPortForwarderName string) (portForwarder PortForwarder) {
	if customPortForwarderName != "" {
		provider = providers.Get(customPortForwarderName)
	}
	portForwarder, ok := provider.(PortForwarder)
	if ok {
		return portForwarder
	}
	return newNoPortForwarder(provider.Name())
}

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

type noPortForwarder struct {
	providerName string
}

func newNoPortForwarder(providerName string) *noPortForwarder {
	return &noPortForwarder{
		providerName: providerName,
	}
}

var ErrPortForwardingNotSupported = errors.New("custom port forwarding obtention is not supported")

func (n *noPortForwarder) Name() string {
	return n.providerName
}

func (n *noPortForwarder) PortForward(context.Context, pfutils.PortForwardObjects) (
	port uint16, err error) {
	return 0, fmt.Errorf("%w: for %s", ErrPortForwardingNotSupported, n.providerName)
}

func (n *noPortForwarder) KeepPortForward(context.Context, pfutils.PortForwardObjects) (err error) {
	return fmt.Errorf("%w: for %s", ErrPortForwardingNotSupported, n.providerName)
}
