package vpn

import (
	"context"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/golibs/logging"
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

		portForwarding := settings.Provider.PortForwarding.Enabled
		var vpnRunner vpnRunner
		var serverName, vpnInterface string
		var err error
		subLogger := l.logger.NewChild(logging.Settings{Prefix: settings.Type + ": "})
		if settings.Type == constants.OpenVPN {
			vpnInterface = settings.OpenVPN.Interface
			vpnRunner, serverName, err = setupOpenVPN(ctx, l.fw,
				l.openvpnConf, providerConf, settings, l.starter, subLogger)
		} else { // Wireguard
			vpnInterface = settings.Wireguard.Interface
			vpnRunner, serverName, err = setupWireguard(ctx, l.netLinker, l.fw, providerConf, settings, subLogger)
		}
		if err != nil {
			l.crashed(ctx, err)
			continue
		}
		tunnelUpData := tunnelUpData{
			portForwarding: portForwarding,
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
				l.cleanup(context.Background(), portForwarding)
				openvpnCancel()
				<-waitError
				close(waitError)
				return
			case <-l.stop:
				l.userTrigger = true
				l.logger.Info("stopping")
				l.cleanup(context.Background(), portForwarding)
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

				l.cleanup(context.Background(), portForwarding)
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
