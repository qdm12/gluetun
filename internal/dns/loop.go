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
	"github.com/qdm12/gluetun/internal/dns/state"
	"github.com/qdm12/gluetun/internal/loopstate"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/logging"
)

var _ Looper = (*Loop)(nil)

type Looper interface {
	Runner
	RestartTickerRunner
	loopstate.Applier
	loopstate.Getter
	SettingsGetterSetter
}

type Loop struct {
	statusManager loopstate.Manager
	state         state.Manager
	conf          unbound.Configurator
	resolvConf    string
	blockBuilder  blacklist.Builder
	client        *http.Client
	logger        logging.Logger
	userTrigger   bool
	start         <-chan struct{}
	running       chan<- models.LoopStatus
	stop          <-chan struct{}
	stopped       chan<- struct{}
	updateTicker  <-chan struct{}
	backoffTime   time.Duration
	timeNow       func() time.Time
	timeSince     func(time.Time) time.Duration
}

const defaultBackoffTime = 10 * time.Second

func NewLoop(conf unbound.Configurator, settings configuration.DNS, client *http.Client,
	logger logging.Logger) *Loop {
	start := make(chan struct{})
	running := make(chan models.LoopStatus)
	stop := make(chan struct{})
	stopped := make(chan struct{})
	updateTicker := make(chan struct{})

	statusManager := loopstate.New(constants.Stopped, start, running, stop, stopped)
	state := state.New(statusManager, settings, updateTicker)

	return &Loop{
		statusManager: statusManager,
		state:         state,
		conf:          conf,
		resolvConf:    "/etc/resolv.conf",
		blockBuilder:  blacklist.NewBuilder(client),
		client:        client,
		logger:        logger,
		userTrigger:   true,
		start:         start,
		running:       running,
		stop:          stop,
		stopped:       stopped,
		updateTicker:  updateTicker,
		backoffTime:   defaultBackoffTime,
		timeNow:       time.Now,
		timeSince:     time.Since,
	}
}

func (l *Loop) logAndWait(ctx context.Context, err error) {
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

func (l *Loop) signalOrSetStatus(status models.LoopStatus) {
	if l.userTrigger {
		l.userTrigger = false
		select {
		case l.running <- status:
		default: // receiver dropped out - avoid deadlock on events routing when shutting down
		}
	} else {
		l.statusManager.SetStatus(status)
	}
}
