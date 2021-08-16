package openvpn

import (
	"context"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
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
		settings, allServers := l.state.GetSettingsAndServers()

		providerConf := provider.New(settings.Provider.Name, allServers, time.Now)

		var connection models.OpenVPNConnection
		var lines []string
		var err error
		if settings.Config == "" {
			connection, err = providerConf.GetOpenVPNConnection(settings.Provider.ServerSelection)
			if err != nil {
				l.signalOrSetStatus(constants.Crashed)
				l.logAndWait(ctx, err)
				continue
			}
			lines = providerConf.BuildConf(connection, l.username, settings)
		} else {
			lines, connection, err = l.processCustomConfig(settings)
			if err != nil {
				l.signalOrSetStatus(constants.Crashed)
				l.logAndWait(ctx, err)
				continue
			}
		}

		if err := l.writeOpenvpnConf(lines); err != nil {
			l.signalOrSetStatus(constants.Crashed)
			l.logAndWait(ctx, err)
			continue
		}

		if settings.User != "" {
			err := l.conf.WriteAuthFile(
				settings.User, settings.Password, l.puid, l.pgid)
			if err != nil {
				l.signalOrSetStatus(constants.Crashed)
				l.logAndWait(ctx, err)
				continue
			}
		}

		if err := l.fw.SetVPNConnection(ctx, connection); err != nil {
			l.signalOrSetStatus(constants.Crashed)
			l.logAndWait(ctx, err)
			continue
		}

		openvpnCtx, openvpnCancel := context.WithCancel(context.Background())

		stdoutLines, stderrLines, waitError, err := l.conf.Start(
			openvpnCtx, settings.Version, settings.Flags)
		if err != nil {
			openvpnCancel()
			l.signalOrSetStatus(constants.Crashed)
			l.logAndWait(ctx, err)
			continue
		}

		linesCollectionCtx, linesCollectionCancel := context.WithCancel(context.Background())
		lineCollectionDone := make(chan struct{})
		go l.collectLines(linesCollectionCtx, lineCollectionDone,
			stdoutLines, stderrLines)
		closeStreams := func() {
			linesCollectionCancel()
			<-lineCollectionDone
		}

		l.backoffTime = defaultBackoffTime
		l.signalOrSetStatus(constants.Running)

		stayHere := true
		for stayHere {
			select {
			case <-l.startPFCh:
				l.startPortForwarding(ctx, settings.Provider.PortForwarding.Enabled,
					providerConf, connection.Hostname)
			case <-ctx.Done():
				const pfTimeout = 100 * time.Millisecond
				l.stopPortForwarding(context.Background(),
					settings.Provider.PortForwarding.Enabled, pfTimeout)
				openvpnCancel()
				<-waitError
				close(waitError)
				closeStreams()
				return
			case <-l.stop:
				l.userTrigger = true
				l.logger.Info("stopping")
				l.stopPortForwarding(ctx, settings.Provider.PortForwarding.Enabled, 0)
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
