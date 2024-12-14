package vpn

import (
	"context"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/log"
)

func (l *Loop) Run(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	select {
	case <-l.start:
	case <-ctx.Done():
		return
	}

	for ctx.Err() == nil {
		settings := l.state.GetSettings()

		providerConf := l.providers.Get(settings.Provider.Name)

		portForwarder := getPortForwarder(providerConf, l.providers,
			*settings.Provider.PortForwarding.Provider)

		var vpnRunner interface {
			Run(ctx context.Context, waitError chan<- error, tunnelReady chan<- struct{})
		}
		var serverName, vpnInterface string
		var canPortForward bool
		var err error
		subLogger := l.logger.New(log.SetComponent(settings.Type))
		if settings.Type == vpn.OpenVPN {
			vpnInterface = settings.OpenVPN.Interface
			vpnRunner, serverName, canPortForward, err = setupOpenVPN(ctx, l.fw,
				l.openvpnConf, providerConf, settings, l.ipv6SupportLevel, l.starter, subLogger)
		} else { // Wireguard
			vpnInterface = settings.Wireguard.Interface
			vpnRunner, serverName, canPortForward, err = setupWireguard(ctx, l.netLinker, l.fw,
				providerConf, settings, l.ipv6SupportLevel, subLogger)
		}
		if err != nil {
			l.crashed(ctx, err)
			continue
		}
		tunnelUpData := tunnelUpData{
			serverName:     serverName,
			canPortForward: canPortForward,
			portForwarder:  portForwarder,
			vpnIntf:        vpnInterface,
			username:       settings.Provider.PortForwarding.Username,
			password:       settings.Provider.PortForwarding.Password,
		}

		openvpnCtx, openvpnCancel := context.WithCancel(context.Background())
		waitError := make(chan error)
		tunnelReady := make(chan struct{})

		go vpnRunner.Run(openvpnCtx, waitError, tunnelReady)

		if err := l.waitForError(ctx, waitError); err != nil {
			openvpnCancel()
			l.crashed(ctx, err)
			continue
		}

		l.backoffTime = defaultBackoffTime
		l.signalOrSetStatus(constants.Running)

		stayHere := true
		for stayHere {
			select {
			case <-tunnelReady:
				go l.onTunnelUp(openvpnCtx, tunnelUpData)
			case <-ctx.Done():
				l.cleanup()
				openvpnCancel()
				<-waitError
				close(waitError)
				return
			case <-l.stop:
				l.userTrigger = true
				l.logger.Info("stopping")
				l.cleanup()
				openvpnCancel()
				<-waitError
				// do not close waitError or the waitError
				// select case will trigger
				l.stopped <- struct{}{}
			case <-l.start:
				l.userTrigger = true
				l.logger.Info("starting")
				stayHere = false
			case err := <-waitError: // unexpected error
				l.statusManager.Lock() // prevent SetStatus from running in parallel

				l.cleanup()
				openvpnCancel()
				l.statusManager.SetStatus(constants.Crashed)
				l.logAndWait(ctx, err)
				stayHere = false

				l.statusManager.Unlock()
			}
		}
		openvpnCancel()
	}
}
