package httpproxy

import (
	"context"

	"github.com/qdm12/gluetun/internal/constants"
)

func (l *Loop) Run(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	if !*l.state.GetSettings().Enabled {
		select {
		case <-l.start:
		case <-ctx.Done():
			return
		}
	}

	for ctx.Err() == nil {
		settings := l.state.GetSettings()
		server, err := New(settings.ListeningAddress, l.logger,
			*settings.Stealth, *settings.Log, *settings.User,
			*settings.Password, settings.ReadHeaderTimeout, settings.ReadTimeout)
		if err != nil {
			l.statusManager.SetStatus(constants.Crashed)
			l.logAndWait(ctx, err)
			continue
		}

		errorCh, err := server.Start(ctx)
		if err != nil {
			l.statusManager.SetStatus(constants.Crashed)
			l.logAndWait(ctx, err)
			continue
		}

		// TODO stable timer, check Shadowsocks
		if l.userTrigger {
			l.running <- constants.Running
			l.userTrigger = false
		} else {
			l.backoffTime = defaultBackoffTime
			l.statusManager.SetStatus(constants.Running)
		}

		stayHere := true
		for stayHere {
			select {
			case <-ctx.Done():
				_ = server.Stop()
				return
			case <-l.start:
				l.userTrigger = true
				l.logger.Info("starting")
				_ = server.Stop()
				stayHere = false
			case <-l.stop:
				l.userTrigger = true
				l.logger.Info("stopping")
				_ = server.Stop()
				l.stopped <- struct{}{}
			case err := <-errorCh:
				l.statusManager.SetStatus(constants.Crashed)
				l.logAndWait(ctx, err)
				stayHere = false
			}
		}
	}
}
