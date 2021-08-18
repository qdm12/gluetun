package openvpn

import (
	"context"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/provider"
)

type Runner interface {
	Run(ctx context.Context, done chan<- struct{})
}

func (l *Loop) Run(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	select {
	case <-l.start:
	case <-ctx.Done():
		return
	}

	for ctx.Err() == nil {
		VPNSettings, providerSettings, allServers := l.state.GetSettingsAndServers()

		providerConf := provider.New(providerSettings.Name, allServers, time.Now)

		serverName, err := setup(ctx, l.fw, l.openvpnConf, providerConf, VPNSettings.OpenVPN, providerSettings)
		if err != nil {
			l.crashed(ctx, err)
			continue
		}
		tunnelUpData := tunnelUpData{
			portForwarding: providerSettings.PortForwarding.Enabled,
			serverName:     serverName,
			portForwarder:  providerConf,
		}

		openvpnCtx, openvpnCancel := context.WithCancel(context.Background())
		waitError := make(chan error)
		tunnelReady := make(chan struct{})

		go l.openvpnConf.Run(openvpnCtx, waitError, tunnelReady,
			l.logger, VPNSettings.OpenVPN)

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
					providerSettings.PortForwarding.Enabled, pfTimeout)
				openvpnCancel()
				<-waitError
				close(waitError)
				return
			case <-l.stop:
				l.userTrigger = true
				l.logger.Info("stopping")
				l.stopPortForwarding(ctx, providerSettings.PortForwarding.Enabled, 0)
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

				l.stopPortForwarding(ctx, providerSettings.PortForwarding.Enabled, 0)
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
