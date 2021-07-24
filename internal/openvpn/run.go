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

		lineCollectionDone := make(chan struct{})
		go l.collectLines(stdoutLines, stderrLines, lineCollectionDone)
		closeStreams := func() {
			close(stdoutLines)
			close(stderrLines)
			<-lineCollectionDone
		}

		// Needs the stream line from main.go to know when the tunnel is up
		portForwardDone := make(chan struct{})
		go func(ctx context.Context) {
			defer close(portForwardDone)
			select {
			// TODO have a way to disable pf with a context
			case <-ctx.Done():
				return
			case gateway := <-l.portForwardSignals:
				l.portForward(ctx, providerConf, l.client, gateway)
			}
		}(openvpnCtx)

		l.backoffTime = defaultBackoffTime
		l.signalOrSetStatus(constants.Running)

		stayHere := true
		for stayHere {
			select {
			case <-ctx.Done():
				openvpnCancel()
				<-waitError
				close(waitError)
				closeStreams()
				<-portForwardDone
				return
			case <-l.stop:
				l.userTrigger = true
				l.logger.Info("stopping")
				openvpnCancel()
				<-waitError
				// do not close waitError or the waitError
				// select case will trigger
				closeStreams()
				<-portForwardDone
				l.stopped <- struct{}{}
			case <-l.start:
				l.userTrigger = true
				l.logger.Info("starting")
				stayHere = false
			case err := <-waitError: // unexpected error
				close(waitError)
				closeStreams()

				l.statusManager.Lock() // prevent SetStatus from running in parallel

				openvpnCancel()
				l.statusManager.SetStatus(constants.Crashed)
				<-portForwardDone
				l.logAndWait(ctx, err)
				stayHere = false

				l.statusManager.Unlock()
			}
		}
		openvpnCancel()
	}
}
