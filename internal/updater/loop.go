package updater

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/storage"
)

type Looper interface {
	Run(ctx context.Context, done chan<- struct{})
	RunRestartTicker(ctx context.Context, done chan<- struct{})
	GetStatus() (status models.LoopStatus)
	SetStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
	GetSettings() (settings settings.Updater)
	SetSettings(settings settings.Updater) (outcome string)
}

type looper struct {
	state state
	// Objects
	updater       Updater
	flusher       storage.Flusher
	setAllServers func(allServers models.AllServers)
	logger        infoErrorer
	// Internal channels and locks
	loopLock     sync.Mutex
	start        chan struct{}
	running      chan models.LoopStatus
	stop         chan struct{}
	stopped      chan struct{}
	updateTicker chan struct{}
	backoffTime  time.Duration
	// Mock functions
	timeNow   func() time.Time
	timeSince func(time.Time) time.Duration
}

const defaultBackoffTime = 5 * time.Second

func NewLooper(settings settings.Updater, currentServers models.AllServers,
	flusher storage.Flusher, setAllServers func(allServers models.AllServers),
	client *http.Client, logger Logger) Looper {
	return &looper{
		state: state{
			status:   constants.Stopped,
			settings: settings,
		},
		updater:       New(settings, client, currentServers, logger),
		flusher:       flusher,
		setAllServers: setAllServers,
		logger:        logger,
		start:         make(chan struct{}),
		running:       make(chan models.LoopStatus),
		stop:          make(chan struct{}),
		stopped:       make(chan struct{}),
		updateTicker:  make(chan struct{}),
		timeNow:       time.Now,
		timeSince:     time.Since,
		backoffTime:   defaultBackoffTime,
	}
}

func (l *looper) logAndWait(ctx context.Context, err error) {
	if err != nil {
		l.logger.Error(err.Error())
	}
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

func (l *looper) Run(ctx context.Context, done chan<- struct{}) {
	defer close(done)
	crashed := false
	select {
	case <-l.start:
	case <-ctx.Done():
		return
	}

	for ctx.Err() == nil {
		updateCtx, updateCancel := context.WithCancel(ctx)

		serversCh := make(chan models.AllServers)
		errorCh := make(chan error)
		runWg := &sync.WaitGroup{}
		runWg.Add(1)
		go func() {
			defer runWg.Done()
			servers, err := l.updater.UpdateServers(updateCtx)
			if err != nil {
				if updateCtx.Err() == nil {
					errorCh <- err
				}
				return
			}
			serversCh <- servers
		}()

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
				updateCancel()
				runWg.Wait()
				close(errorCh)
				return
			case <-l.start:
				l.logger.Info("starting")
				updateCancel()
				runWg.Wait()
				stayHere = false
			case <-l.stop:
				l.logger.Info("stopping")
				updateCancel()
				runWg.Wait()
				l.stopped <- struct{}{}
			case servers := <-serversCh:
				l.setAllServers(servers)
				if err := l.flusher.FlushToFile(servers); err != nil {
					l.logger.Error(err.Error())
				}
				runWg.Wait()
				l.state.setStatusWithLock(constants.Completed)
				l.logger.Info("Updated servers information")
			case err := <-errorCh:
				close(serversCh)
				runWg.Wait()
				l.state.setStatusWithLock(constants.Crashed)
				l.logAndWait(ctx, err)
				crashed = true
				stayHere = false
			}
		}
		updateCancel()
		close(errorCh)
	}
}

func (l *looper) RunRestartTicker(ctx context.Context, done chan<- struct{}) {
	defer close(done)
	timer := time.NewTimer(time.Hour)
	timer.Stop()
	timerIsStopped := true
	if period := *l.GetSettings().Period; period > 0 {
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
			timer.Reset(*l.GetSettings().Period)
		case <-l.updateTicker:
			if !timerIsStopped && !timer.Stop() {
				<-timer.C
			}
			timerIsStopped = true
			period := *l.GetSettings().Period
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
