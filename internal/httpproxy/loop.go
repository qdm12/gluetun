// Package httpproxy defines an interface to run an HTTP(s) proxy server.
package httpproxy

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/logging"
)

type Looper interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	SetStatus(status models.LoopStatus) (outcome string, err error)
	GetStatus() (status models.LoopStatus)
	GetSettings() (settings configuration.HTTPProxy)
	SetSettings(settings configuration.HTTPProxy) (outcome string)
}

type looper struct {
	state state
	// Other objects
	logger logging.Logger
	// Internal channels and locks
	loopLock      sync.Mutex
	running       chan models.LoopStatus
	stop, stopped chan struct{}
	start         chan struct{}
	backoffTime   time.Duration
}

const defaultBackoffTime = 10 * time.Second

func NewLooper(logger logging.Logger, settings configuration.HTTPProxy) Looper {
	return &looper{
		state: state{
			status:   constants.Stopped,
			settings: settings,
		},
		logger:      logger.NewChild(logging.SetPrefix("http proxy: ")),
		start:       make(chan struct{}),
		running:     make(chan models.LoopStatus),
		stop:        make(chan struct{}),
		stopped:     make(chan struct{}),
		backoffTime: defaultBackoffTime,
	}
}

func (l *looper) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	crashed := false

	if l.GetSettings().Enabled {
		go func() {
			_, _ = l.SetStatus(constants.Running)
		}()
	}

	select {
	case <-l.start:
	case <-ctx.Done():
		return
	}

	defer l.logger.Warn("loop exited")

	for ctx.Err() == nil {
		runCtx, runCancel := context.WithCancel(ctx)

		settings := l.GetSettings()
		address := fmt.Sprintf(":%d", settings.Port)
		server := New(runCtx, address, l.logger, settings.Stealth, settings.Log, settings.User, settings.Password)

		runWg := &sync.WaitGroup{}
		runWg.Add(1)
		errorCh := make(chan error)
		go server.Run(runCtx, runWg, errorCh)

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
				l.logger.Warn("context canceled: exiting loop")
				runCancel()
				runWg.Wait()
				return
			case <-l.start:
				l.logger.Info("starting")
				runCancel()
				runWg.Wait()
				stayHere = false
			case <-l.stop:
				l.logger.Info("stopping")
				runCancel()
				runWg.Wait()
				l.stopped <- struct{}{}
			case err := <-errorCh:
				runWg.Wait()
				l.state.setStatusWithLock(constants.Crashed)
				l.logAndWait(ctx, err)
				crashed = true
				stayHere = false
			}
		}
		runCancel() // repetition for linter only
	}
}

func (l *looper) logAndWait(ctx context.Context, err error) {
	l.logger.Error(err)
	l.logger.Info("retrying in %s", l.backoffTime)
	timer := time.NewTimer(l.backoffTime)
	l.backoffTime *= 2
	select {
	case <-timer.C:
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
	}
}
