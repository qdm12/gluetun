package updater

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/storage"
	"github.com/qdm12/golibs/logging"
)

type Looper interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	RunRestartTicker(ctx context.Context)
	Restart()
	Stop()
	GetPeriod() (period time.Duration)
	SetPeriod(period time.Duration)
}

type looper struct {
	period        time.Duration
	periodMutex   sync.RWMutex
	updater       Updater
	storage       storage.Storage
	setAllServers func(allServers models.AllServers)
	logger        logging.Logger
	restart       chan struct{}
	stop          chan struct{}
	updateTicker  chan struct{}
}

func NewLooper(options Options, period time.Duration, currentServers models.AllServers,
	storage storage.Storage, setAllServers func(allServers models.AllServers),
	client *http.Client, logger logging.Logger) Looper {
	loggerWithPrefix := logger.WithPrefix("updater: ")
	return &looper{
		period:        period,
		updater:       New(options, client, currentServers, loggerWithPrefix),
		storage:       storage,
		setAllServers: setAllServers,
		logger:        loggerWithPrefix,
		restart:       make(chan struct{}),
		stop:          make(chan struct{}),
		updateTicker:  make(chan struct{}),
	}
}

func (l *looper) Restart() { l.restart <- struct{}{} }
func (l *looper) Stop()    { l.stop <- struct{}{} }

func (l *looper) GetPeriod() (period time.Duration) {
	l.periodMutex.RLock()
	defer l.periodMutex.RUnlock()
	return l.period
}

func (l *looper) SetPeriod(period time.Duration) {
	l.periodMutex.Lock()
	l.period = period
	l.periodMutex.Unlock()
	l.updateTicker <- struct{}{}
}

func (l *looper) logAndWait(ctx context.Context, err error) {
	l.logger.Error(err)
	l.logger.Info("retrying in 5 minutes")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel() // just for the linter
	<-ctx.Done()
}

func (l *looper) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-l.restart:
		l.logger.Info("starting...")
	case <-ctx.Done():
		return
	}
	defer l.logger.Warn("loop exited")

	enabled := true

	for ctx.Err() == nil {
		for !enabled {
			// wait for a signal to re-enable
			select {
			case <-l.stop:
				l.logger.Info("already disabled")
			case <-l.restart:
				enabled = true
			case <-ctx.Done():
				return
			}
		}

		// Enabled and has a period set

		servers, err := l.updater.UpdateServers(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			l.logAndWait(ctx, err)
			continue
		}
		l.setAllServers(servers)
		if err := l.storage.FlushToFile(servers); err != nil {
			l.logger.Error(err)
		}
		l.logger.Info("Updated servers information")

		select {
		case <-l.restart: // triggered restart
		case <-l.stop:
			enabled = false
		case <-ctx.Done():
			return
		}
	}
}

func (l *looper) RunRestartTicker(ctx context.Context) {
	ticker := time.NewTicker(time.Hour)
	period := l.GetPeriod()
	if period > 0 {
		ticker = time.NewTicker(period)
	} else {
		ticker.Stop()
	}
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			l.restart <- struct{}{}
		case <-l.updateTicker:
			ticker.Stop()
			ticker = time.NewTicker(l.GetPeriod())
		}
	}
}
