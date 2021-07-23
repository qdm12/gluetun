// Package httpproxy defines an interface to run an HTTP(s) proxy server.
package httpproxy

import (
	"context"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/logging"
)

type Looper interface {
	Run(ctx context.Context, done chan<- struct{})
	SetStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
	GetStatus() (status models.LoopStatus)
	GetSettings() (settings configuration.HTTPProxy)
	SetSettings(ctx context.Context, settings configuration.HTTPProxy) (
		outcome string)
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
		logger:      logger,
		start:       make(chan struct{}),
		running:     make(chan models.LoopStatus),
		stop:        make(chan struct{}),
		stopped:     make(chan struct{}),
		backoffTime: defaultBackoffTime,
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
