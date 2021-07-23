package httpproxy

import (
	"context"
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
)

func (l *looper) Run(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	crashed := false

	if l.GetSettings().Enabled {
		go func() {
			_, _ = l.SetStatus(ctx, constants.Running)
		}()
	}

	select {
	case <-l.start:
	case <-ctx.Done():
		return
	}

	for ctx.Err() == nil {
		runCtx, runCancel := context.WithCancel(ctx)

		settings := l.GetSettings()
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
			l.state.setStatusWithLock(constants.Running)
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
				l.state.setStatusWithLock(constants.Crashed)
				l.logAndWait(ctx, err)
				crashed = true
				stayHere = false
			}
		}
		runCancel() // repetition for linter only
	}
}
