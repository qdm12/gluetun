package updater

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/storage"
	"github.com/qdm12/golibs/logging"
)

type Looper interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	RunRestartTicker(ctx context.Context, wg *sync.WaitGroup)
	GetStatus() (status models.LoopStatus)
	SetStatus(status models.LoopStatus) (outcome string, err error)
	GetPeriod() (period time.Duration)
	SetPeriod(period time.Duration)
}

type looper struct {
	state state
	// Objects
	updater       Updater
	storage       storage.Storage
	setAllServers func(allServers models.AllServers)
	logger        logging.Logger
	// Internal channels and locks
	loopLock     sync.Mutex
	start        chan struct{}
	updateTicker chan struct{}
	// Mock functions
	timeNow   func() time.Time
	timeSince func(time.Time) time.Duration
}

func NewLooper(options Options, period time.Duration, currentServers models.AllServers,
	storage storage.Storage, setAllServers func(allServers models.AllServers),
	client *http.Client, logger logging.Logger) Looper {
	loggerWithPrefix := logger.WithPrefix("updater: ")
	return &looper{
		state: state{
			status: constants.Stopped,
			period: period,
		},
		updater:       New(options, client, currentServers, loggerWithPrefix),
		storage:       storage,
		setAllServers: setAllServers,
		logger:        loggerWithPrefix,
		start:         make(chan struct{}),
		updateTicker:  make(chan struct{}),
		timeNow:       time.Now,
		timeSince:     time.Since,
	}
}

func (l *looper) logAndWait(ctx context.Context, err error) {
	l.logger.Error(err)
	const waitTime = 5 * time.Minute
	l.logger.Info("retrying in %s", waitTime)
	timer := time.NewTimer(waitTime)
	select {
	case <-timer.C:
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
	}
}

func (l *looper) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	crashed := false
	select {
	case <-l.start:
	case <-ctx.Done():
		return
	}
	defer l.logger.Warn("loop exited")

	for ctx.Err() == nil {
		if !crashed {
			l.state.setStatusWithLock(constants.Running)
		}

		servers, err := l.updater.UpdateServers(ctx)

		if err != nil {
			crashed = true
			l.state.setStatusWithLock(constants.Crashed)
			if ctx.Err() != nil {
				return
			}
			l.logAndWait(ctx, err)
			continue
		}
		crashed = false

		l.setAllServers(servers)
		if err := l.storage.FlushToFile(servers); err != nil {
			l.logger.Error(err)
		}
		l.logger.Info("Updated servers information")
		l.state.setStatusWithLock(constants.Completed)

		select {
		case <-ctx.Done():
			l.logger.Warn("context canceled: exiting loop")
			return
		case <-l.start:
			l.logger.Info("starting")
		}
	}
}

func (l *looper) RunRestartTicker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	timer := time.NewTimer(time.Hour)
	timer.Stop()
	timerIsStopped := true
	if period := l.GetPeriod(); period > 0 {
		timerIsStopped = false
		timer.Reset(period)
	}
	lastTick := time.Unix(0, 0)
	for {
		select {
		case <-ctx.Done():
			if !timerIsStopped && !timer.Stop() {
				<-timer.C
			}
			return
		case <-timer.C:
			lastTick = l.timeNow()
			l.start <- struct{}{}
			timer.Reset(l.GetPeriod())
		case <-l.updateTicker:
			if !timerIsStopped && !timer.Stop() {
				<-timer.C
			}
			timerIsStopped = true
			period := l.GetPeriod()
			if period == 0 {
				continue
			}
			var waited time.Duration
			if lastTick.UnixNano() > 0 {
				waited = l.timeSince(lastTick)
			}
			leftToWait := period - waited
			timer.Reset(leftToWait)
			timerIsStopped = false
		}
	}
}
