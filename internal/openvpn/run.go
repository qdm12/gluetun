package openvpn

import (
	"context"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/openvpn/custom"
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
		openVPNSettings, providerSettings, allServers := l.state.GetSettingsAndServers()

		providerConf := provider.New(providerSettings.Name, allServers, time.Now)

		var connection models.OpenVPNConnection
		var lines []string
		var err error
		if openVPNSettings.Config == "" {
			connection, err = providerConf.GetOpenVPNConnection(providerSettings.ServerSelection)
			if err == nil {
				lines = providerConf.BuildConf(connection, openVPNSettings)
			}
		} else {
			lines, connection, err = custom.BuildConfig(openVPNSettings)
		}
		if err != nil {
			l.crashed(ctx, err)
			continue
		}

		if err := l.openvpnConf.WriteConfig(lines); err != nil {
			l.crashed(ctx, err)
			continue
		}

		if openVPNSettings.User != "" {
			err := l.openvpnConf.WriteAuthFile(openVPNSettings.User, openVPNSettings.Password)
			if err != nil {
				l.crashed(ctx, err)
				continue
			}
		}

		if err := l.fw.SetVPNConnection(ctx, connection); err != nil {
			l.crashed(ctx, err)
			continue
		}

		openvpnCtx, openvpnCancel := context.WithCancel(context.Background())

		stdoutLines, stderrLines, waitError, err := l.openvpnConf.Start(
			openvpnCtx, openVPNSettings.Version, openVPNSettings.Flags)
		if err != nil {
			openvpnCancel()
			l.crashed(ctx, err)
			continue
		}

		linesCollectionCtx, linesCollectionCancel := context.WithCancel(context.Background())
		lineCollectionDone := make(chan struct{})
		tunnelUpData := tunnelUpData{
			portForwarding: providerSettings.PortForwarding.Enabled,
			serverName:     connection.Hostname,
			portForwarder:  providerConf,
		}
		go l.collectLines(linesCollectionCtx, lineCollectionDone,
			stdoutLines, stderrLines, tunnelUpData)
		closeStreams := func() {
			linesCollectionCancel()
			<-lineCollectionDone
		}

		l.backoffTime = defaultBackoffTime
		l.signalOrSetStatus(constants.Running)

		stayHere := true
		for stayHere {
			select {
			case <-ctx.Done():
				const pfTimeout = 100 * time.Millisecond
				l.stopPortForwarding(context.Background(),
					providerSettings.PortForwarding.Enabled, pfTimeout)
				openvpnCancel()
				<-waitError
				close(waitError)
				closeStreams()
				return
			case <-l.stop:
				l.userTrigger = true
				l.logger.Info("stopping")
				l.stopPortForwarding(ctx, providerSettings.PortForwarding.Enabled, 0)
				openvpnCancel()
				<-waitError
				// do not close waitError or the waitError
				// select case will trigger
				closeStreams()
				l.stopped <- struct{}{}
			case <-l.start:
				l.userTrigger = true
				l.logger.Info("starting")
				stayHere = false
			case err := <-waitError: // unexpected error
				close(waitError)
				closeStreams()

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
