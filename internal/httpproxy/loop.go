package httpproxy

import (
	"context"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/httpproxy/state"
	"github.com/qdm12/gluetun/internal/loopstate"
	"github.com/qdm12/gluetun/internal/models"
)

type Loop struct {
	statusManager *loopstate.State
	state         *state.State
	// Other objects
	logger Logger
	// Internal channels and locks
	running       chan models.LoopStatus
	stop, stopped chan struct{}
	start         chan struct{}
	userTrigger   bool
	backoffTime   time.Duration
}

const defaultBackoffTime = 10 * time.Second

func NewLoop(logger Logger, settings settings.HTTPProxy) *Loop {
	start := make(chan struct{})
	running := make(chan models.LoopStatus)
	stop := make(chan struct{})
	stopped := make(chan struct{})

	statusManager := loopstate.New(constants.Stopped,
		start, running, stop, stopped)
	state := state.New(statusManager, settings)

	return &Loop{
		statusManager: statusManager,
		state:         state,
		logger:        logger,
		start:         start,
		running:       running,
		stop:          stop,
		stopped:       stopped,
		userTrigger:   true,
		backoffTime:   defaultBackoffTime,
	}
}

func (l *Loop) logAndWait(ctx context.Context, err error) {
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
