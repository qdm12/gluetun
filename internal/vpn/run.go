package vpn

import (
	"context"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/provider"
)

type Runner interface {
	Run(ctx context.Context, done chan<- struct{})
}

type vpnRunner interface {
	Run(ctx context.Context, errCh chan<- error, ready chan<- struct{})
}

func (l *Loop) Run(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	select {
	case <-l.start:
	case <-ctx.Done():
		return
	}

	for ctx.Err() == nil {
		settings, allServers := l.state.GetSettingsAndServers()

		providerConf := provider.New(settings.Provider.Name, allServers, time.Now)

		var vpnRunner vpnRunner
		var serverName, vpnInterface string
		var err error
		if settings.Type == constants.OpenVPN {
			vpnInterface = settings.OpenVPN.Interface
			vpnRunner, serverName, err = setupOpenVPN(ctx, l.fw,
				l.openvpnConf, providerConf, settings, l.starter, l.logger)
		} else { // Wireguard
			vpnInterface = settings.Wireguard.Interface
			vpnRunner, serverName, err = setupWireguard(ctx, l.netLinker, l.fw, providerConf, settings, l.logger)
		}
		if err != nil {
			l.crashed(ctx, err)
			continue
		}
		tunnelUpData := tunnelUpData{
			portForwarding: settings.Provider.PortForwarding.Enabled,
			serverName:     serverName,
			portForwarder:  providerConf,
			vpnIntf:        vpnInterface,
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
				const pfTimeout = 100 * time.Millisecond
				l.stopPortForwarding(context.Background(),
					settings.Provider.PortForwarding.Enabled, pfTimeout)
				openvpnCancel()
				<-waitError
				close(waitError)
				return
			case <-l.stop:
				l.userTrigger = true
				l.logger.Info("stopping")
				l.stopPortForwarding(ctx, settings.Provider.PortForwarding.Enabled, 0)
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
				close(waitError)

				l.statusManager.Lock() // prevent SetStatus from running in parallel

				l.stopPortForwarding(ctx, settings.Provider.PortForwarding.Enabled, 0)
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
