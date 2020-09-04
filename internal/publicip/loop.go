package publicip

import (
	"context"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
)

type Looper interface {
	Run(ctx context.Context)
	RunRestartTicker(ctx context.Context)
	Restart()
	Stop()
	GetPeriod() (period time.Duration)
	SetPeriod(period time.Duration)
}

type looper struct {
	period           time.Duration
	periodMutex      sync.RWMutex
	getter           IPGetter
	logger           logging.Logger
	fileManager      files.FileManager
	ipStatusFilepath models.Filepath
	uid              int
	gid              int
	restart          chan struct{}
	stop             chan struct{}
	updateTicker     chan struct{}
}

func NewLooper(client network.Client, logger logging.Logger, fileManager files.FileManager,
	ipStatusFilepath models.Filepath, period time.Duration, uid, gid int) Looper {
	return &looper{
		period:           period,
		getter:           NewIPGetter(client),
		logger:           logger.WithPrefix("ip getter: "),
		fileManager:      fileManager,
		ipStatusFilepath: ipStatusFilepath,
		uid:              uid,
		gid:              gid,
		restart:          make(chan struct{}),
		stop:             make(chan struct{}),
		updateTicker:     make(chan struct{}),
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
	l.logger.Info("retrying in 5 seconds")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel() // just for the linter
	<-ctx.Done()
}

func (l *looper) Run(ctx context.Context) {
	select {
	case <-l.restart:
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

		ip, err := l.getter.Get()
		if err != nil {
			l.logAndWait(ctx, err)
			continue
		}
		l.logger.Info("Public IP address is %s", ip)
		err = l.fileManager.WriteLinesToFile(
			string(l.ipStatusFilepath),
			[]string{ip.String()},
			files.Ownership(l.uid, l.gid),
			files.Permissions(0600))
		if err != nil {
			l.logAndWait(ctx, err)
			continue
		}
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
