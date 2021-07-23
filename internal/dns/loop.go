// Package dns defines interfaces to interact with DNS and DNS over TLS.
package dns

import (
	"context"
	"net/http"
	"time"

	"github.com/qdm12/dns/pkg/blacklist"
	"github.com/qdm12/dns/pkg/unbound"
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/logging"
)

type Looper interface {
	Run(ctx context.Context, done chan<- struct{})
	RunRestartTicker(ctx context.Context, done chan<- struct{})
	GetStatus() (status models.LoopStatus)
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
	GetSettings() (settings configuration.DNS)
	SetSettings(ctx context.Context, settings configuration.DNS) (
		outcome string)
}

type looper struct {
	state        *state
	conf         unbound.Configurator
	resolvConf   string
	blockBuilder blacklist.Builder
	client       *http.Client
	logger       logging.Logger
	userTrigger  bool
	start        <-chan struct{}
	running      chan<- models.LoopStatus
	stop         <-chan struct{}
	stopped      chan<- struct{}
	updateTicker <-chan struct{}
	backoffTime  time.Duration
	timeNow      func() time.Time
	timeSince    func(time.Time) time.Duration
}

const defaultBackoffTime = 10 * time.Second

func NewLoop(conf unbound.Configurator, settings configuration.DNS, client *http.Client,
	logger logging.Logger) Looper {
	start := make(chan struct{})
	running := make(chan models.LoopStatus)
	stop := make(chan struct{})
	stopped := make(chan struct{})
	updateTicker := make(chan struct{})

	state := newState(constants.Stopped, settings, start, running, stop, stopped, updateTicker)

	return &looper{
		state:        state,
		conf:         conf,
		resolvConf:   "/etc/resolv.conf",
		blockBuilder: blacklist.NewBuilder(client),
		client:       client,
		logger:       logger,
		userTrigger:  true,
		start:        start,
		running:      running,
		stop:         stop,
		stopped:      stopped,
		updateTicker: updateTicker,
		backoffTime:  defaultBackoffTime,
		timeNow:      time.Now,
		timeSince:    time.Since,
	}
}

func (l *looper) logAndWait(ctx context.Context, err error) {
	if err != nil {
		l.logger.Warn(err.Error())
	}
	l.logger.Info("attempting restart in " + l.backoffTime.String())
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

func (l *looper) signalOrSetStatus(status models.LoopStatus) {
	if l.userTrigger {
		l.userTrigger = false
		select {
		case l.running <- status:
		default: // receiver dropped out - avoid deadlock on events routing when shutting down
		}
	} else {
		l.state.SetStatus(status)
	}
}
