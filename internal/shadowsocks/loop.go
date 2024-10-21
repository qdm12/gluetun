package shadowsocks

import (
	"context"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
)

type Loop struct {
	state state
	// Other objects
	logger Logger
	// Internal channels and locks
	refreshing  bool
	refresh     chan struct{}
	changed     chan models.LoopStatus
	backoffTime time.Duration

	runCancel context.CancelFunc
	runDone   <-chan struct{}
}

const defaultBackoffTime = 10 * time.Second

func NewLoop(settings settings.Shadowsocks, logger Logger) *Loop {
	return &Loop{
		state: state{
			status:   constants.Stopped,
			settings: settings,
		},
		logger:      logger,
		refresh:     make(chan struct{}, 1), // capacity of 1 to handle crash auto-restart
		changed:     make(chan models.LoopStatus),
		backoffTime: defaultBackoffTime,
	}
}

func (l *Loop) Start(ctx context.Context) (runError <-chan error, err error) {
	runCtx, runCancel := context.WithCancel(context.Background())
	l.runCancel = runCancel
	ready := make(chan struct{})
	done := make(chan struct{})
	l.runDone = done

	go l.run(runCtx, ready, done)

	<-ready

	return nil, nil //nolint:nilnil
}

func (l *Loop) run(ctx context.Context, ready, done chan<- struct{}) {
	defer close(done)
	close(ready)

	for ctx.Err() == nil {
		// What if update and crash at the same time ish?
		settings := l.GetSettings()

		var service *service
		var runError <-chan error
		var err error
		if *settings.Enabled {
			service = newService(settings.Settings, l.logger)
			runError, err = service.Start(ctx)
			if err != nil {
				runErrorCh := make(chan error, 1)
				runError = runErrorCh
				runErrorCh <- err
			} else if l.refreshing {
				l.changed <- constants.Running
			} else { // auto-restart due to crash
				l.state.setStatusWithLock(constants.Running)
				l.backoffTime = defaultBackoffTime
			}
		} else {
			if l.refreshing {
				l.changed <- constants.Stopped
			} else { // auto-restart due to crash
				l.state.setStatusWithLock(constants.Stopped)
				l.backoffTime = defaultBackoffTime
			}
		}
		l.refreshing = false

		select {
		case <-l.refresh:
			l.refreshing = true
			if service != nil {
				err = service.Stop()
				if err != nil {
					l.logger.Error("stopping service: " + err.Error())
				}
			}
		case err = <-runError:
			if l.refreshing {
				l.changed <- constants.Crashed
			} else {
				l.state.setStatusWithLock(constants.Crashed)
			}
			l.logAndWait(ctx, err)
		case <-ctx.Done():
			if service != nil {
				err = service.Stop()
				if err != nil {
					l.logger.Error("stopping service: " + err.Error())
				}
			}
			return
		}
	}
}

func (l *Loop) Stop() (err error) {
	l.runCancel()
	<-l.runDone
	return nil
}

func (l *Loop) logAndWait(ctx context.Context, err error) {
	if err != nil {
		l.logger.Error(err.Error())
	}
	l.logger.Info("retrying in " + l.backoffTime.String())
	timer := time.NewTimer(l.backoffTime)
	l.backoffTime *= 2
	select {
	case <-timer.C:
	case <-ctx.Done():
		_ = timer.Stop()
	case <-l.refresh: // user-triggered refresh
	}
}
