package httpproxy

import (
	"context"
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
)

type Runner interface {
	Run(ctx context.Context, done chan<- struct{})
}

func (l *Loop) Run(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	crashed := false

	if l.state.GetSettings().Enabled {
		go func() {
			_, _ = l.statusManager.ApplyStatus(ctx, constants.Running)
		}()
	}

	select {
	case <-l.start:
	case <-ctx.Done():
		return
	}

	for ctx.Err() == nil {
		runCtx, runCancel := context.WithCancel(ctx)

		settings := l.state.GetSettings()
		address := fmt.Sprintf(":%d", settings.Port)
		server := New(runCtx, address, l.logger, settings.Stealth, settings.Log, settings.User, settings.Password)

		errorCh := make(chan error)
		go server.Run(runCtx, errorCh)

		// TODO stable timer, check Shadowsocks
		if !crashed {
			l.running <- constants.Running
			crashed = false
		} else {
			l.backoffTime = defaultBackoffTime
			l.statusManager.SetStatus(constants.Running)
		}

		stayHere := true
		for stayHere {
			select {
			case <-ctx.Done():
				runCancel()
				<-errorCh
				return
			case <-l.start:
				l.logger.Info("starting")
				runCancel()
				<-errorCh
				stayHere = false
			case <-l.stop:
				l.logger.Info("stopping")
				runCancel()
				<-errorCh
				l.stopped <- struct{}{}
			case err := <-errorCh:
				l.statusManager.SetStatus(constants.Crashed)
				l.logAndWait(ctx, err)
				crashed = true
				stayHere = false
			}
		}
		runCancel() // repetition for linter only
	}
}
