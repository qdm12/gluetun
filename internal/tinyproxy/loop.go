package tinyproxy

import (
	"context"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
)

type Looper interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	Restart()
	Start()
	Stop()
	GetSettings() (settings settings.TinyProxy)
	SetSettings(settings settings.TinyProxy)
}

type looper struct {
	conf             Configurator
	firewallConf     firewall.Configurator
	settings         settings.TinyProxy
	settingsMutex    sync.RWMutex
	logger           logging.Logger
	streamMerger     command.StreamMerger
	uid              int
	gid              int
	defaultInterface string
	restart          chan struct{}
	start            chan struct{}
	stop             chan struct{}
}

func (l *looper) logAndWait(ctx context.Context, err error) {
	l.logger.Error(err)
	l.logger.Info("retrying in 1 minute")
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel() // just for the linter
	<-ctx.Done()
}

func NewLooper(conf Configurator, firewallConf firewall.Configurator, settings settings.TinyProxy,
	logger logging.Logger, streamMerger command.StreamMerger, uid, gid int, defaultInterface string) Looper {
	return &looper{
		conf:             conf,
		firewallConf:     firewallConf,
		settings:         settings,
		logger:           logger.WithPrefix("tinyproxy: "),
		streamMerger:     streamMerger,
		uid:              uid,
		gid:              gid,
		defaultInterface: defaultInterface,
		restart:          make(chan struct{}),
		start:            make(chan struct{}),
		stop:             make(chan struct{}),
	}
}

func (l *looper) GetSettings() (settings settings.TinyProxy) {
	l.settingsMutex.RLock()
	defer l.settingsMutex.RUnlock()
	return l.settings
}

func (l *looper) SetSettings(settings settings.TinyProxy) {
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
	wg.Add(1)
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

	var previousPort uint16
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
		err := l.conf.MakeConf(settings.LogLevel, settings.Port, settings.User, settings.Password, l.uid, l.gid)
		if err != nil {
			l.logAndWait(ctx, err)
			continue
		}

		if previousPort > 0 {
			if err := l.firewallConf.RemoveAllowedPort(ctx, previousPort); err != nil {
				l.logger.Error(err)
				continue
			}
		}
		if err := l.firewallConf.SetAllowedPort(ctx, settings.Port, l.defaultInterface); err != nil {
			l.logger.Error(err)
			continue
		}
		previousPort = settings.Port

		tinyproxyCtx, tinyproxyCancel := context.WithCancel(context.Background())
		stream, waitFn, err := l.conf.Start(tinyproxyCtx)
		if err != nil {
			tinyproxyCancel()
			l.logAndWait(ctx, err)
			continue
		}
		go l.streamMerger.Merge(tinyproxyCtx, stream, command.MergeName("tinyproxy"))
		waitError := make(chan error)
		go func() {
			err := waitFn() // blocking
			waitError <- err
		}()
		stayHere := true
		for stayHere {
			select {
			case <-ctx.Done():
				l.logger.Warn("context canceled: exiting loop")
				tinyproxyCancel()
				<-waitError
				close(waitError)
				return
			case <-l.restart: // triggered restart
				l.logger.Info("restarting")
				tinyproxyCancel()
				<-waitError
				close(waitError)
				stayHere = false
			case <-l.start:
				l.logger.Info("already started")
			case <-l.stop:
				l.logger.Info("stopping")
				tinyproxyCancel()
				<-waitError
				close(waitError)
				l.setEnabled(false)
				stayHere = false
			case err := <-waitError: // unexpected error
				tinyproxyCancel()
				close(waitError)
				l.logAndWait(ctx, err)
			}
		}
		tinyproxyCancel() // repetition for linter only
	}
}
