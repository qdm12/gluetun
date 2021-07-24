// Package httpproxy defines an interface to run an HTTP(s) proxy server.
package httpproxy

import (
	"context"
	"time"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/httpproxy/state"
	"github.com/qdm12/gluetun/internal/loopstate"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/logging"
)

type Looper interface {
	Runner
	loopstate.Getter
	loopstate.Applier
	SettingsGetterSetter
}

type looper struct {
	statusManager loopstate.Manager
	state         state.Manager
	// Other objects
	logger logging.Logger
	// Internal channels and locks
	running       chan models.LoopStatus
	stop, stopped chan struct{}
	start         chan struct{}
	backoffTime   time.Duration
}

const defaultBackoffTime = 10 * time.Second

func NewLooper(logger logging.Logger, settings configuration.HTTPProxy) Looper {
	start := make(chan struct{})
	running := make(chan models.LoopStatus)
	stop := make(chan struct{})
	stopped := make(chan struct{})

	statusManager := loopstate.New(constants.Stopped,
		start, running, stop, stopped)
	state := state.New(statusManager, settings)

	return &looper{
		statusManager: statusManager,
		state:         state,
		logger:        logger,
		start:         start,
		running:       running,
		stop:          stop,
		stopped:       stopped,
		backoffTime:   defaultBackoffTime,
	}
}

func (l *looper) logAndWait(ctx context.Context, err error) {
	l.logger.Error(err.Error())
	l.logger.Info("retrying in " + l.backoffTime.String())
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
