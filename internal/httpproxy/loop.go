package httpproxy

import (
	"context"
	"fmt"
	"sync"

	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/golibs/logging"
)

type Looper interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	Restart()
	Start()
	Stop()
	GetSettings() (settings settings.HTTPProxy)
	SetSettings(settings settings.HTTPProxy)
}

type looper struct {
	settings      settings.HTTPProxy
	settingsMutex sync.RWMutex
	logger        logging.Logger
	restart       chan struct{}
	start         chan struct{}
	stop          chan struct{}
}

func NewLooper(logger logging.Logger, settings settings.HTTPProxy) Looper {
	return &looper{
		settings: settings,
		logger:   logger.WithPrefix("http proxy: "),
		restart:  make(chan struct{}),
		start:    make(chan struct{}),
		stop:     make(chan struct{}),
	}
}

func (l *looper) GetSettings() (settings settings.HTTPProxy) {
	l.settingsMutex.RLock()
	defer l.settingsMutex.RUnlock()
	return l.settings
}

func (l *looper) SetSettings(settings settings.HTTPProxy) {
	l.settingsMutex.Lock()
	defer l.settingsMutex.Unlock()
	l.settings = settings
}

func (l *looper) isEnabled() bool {
	l.settingsMutex.RLock()
	defer l.settingsMutex.RUnlock()
	return l.settings.Enabled
}

func (l *looper) setEnabled(enabled bool) {
	l.settingsMutex.Lock()
	defer l.settingsMutex.Unlock()
	l.settings.Enabled = enabled
}

func (l *looper) Restart() { l.restart <- struct{}{} }
func (l *looper) Start()   { l.start <- struct{}{} }
func (l *looper) Stop()    { l.stop <- struct{}{} }

func (l *looper) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	waitForStart := true
	for waitForStart {
		select {
		case <-l.stop:
			l.logger.Info("not started yet")
		case <-l.start:
			waitForStart = false
		case <-l.restart:
			waitForStart = false
		case <-ctx.Done():
			return
		}
	}
	defer l.logger.Warn("loop exited")

	for ctx.Err() == nil {
		for !l.isEnabled() {
			// wait for a signal to re-enable
			select {
			case <-l.stop:
				l.logger.Info("already disabled")
			case <-l.restart:
				l.setEnabled(true)
			case <-l.start:
				l.setEnabled(true)
			case <-ctx.Done():
				return
			}
		}

		settings := l.GetSettings()
		address := fmt.Sprintf("0.0.0.0:%d", settings.Port)

		server := New(ctx, address, l.logger, settings.Stealth, settings.Log, settings.User, settings.Password)

		runCtx, runCancel := context.WithCancel(context.Background())
		runWg := &sync.WaitGroup{}
		runWg.Add(1)
		// TODO crashed channel
		go server.Run(runCtx, runWg)

		stayHere := true
		for stayHere {
			select {
			case <-ctx.Done():
				l.logger.Warn("context canceled: exiting loop")
				runCancel()
				runWg.Wait()
				return
			case <-l.restart: // triggered restart
				l.logger.Info("restarting")
				runCancel()
				runWg.Wait()
				stayHere = false
			case <-l.start:
				l.logger.Info("already started")
			case <-l.stop:
				l.logger.Info("stopping")
				runCancel()
				runWg.Wait()
				l.setEnabled(false)
				stayHere = false
			}
		}
		runCancel() // repetition for linter only
	}
}
