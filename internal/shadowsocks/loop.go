package shadowsocks

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/golibs/logging"
	shadowsockslib "github.com/qdm12/ss-server/pkg"
)

type Looper interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	Restart()
	Start()
	Stop()
	GetSettings() (settings settings.ShadowSocks)
	SetSettings(settings settings.ShadowSocks)
}

type looper struct {
	firewallConf     firewall.Configurator
	settings         settings.ShadowSocks
	settingsMutex    sync.RWMutex
	logger           logging.Logger
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

func NewLooper(firewallConf firewall.Configurator, settings settings.ShadowSocks,
	logger logging.Logger, defaultInterface string) Looper {
	return &looper{
		firewallConf:     firewallConf,
		settings:         settings,
		logger:           logger.WithPrefix("shadowsocks: "),
		defaultInterface: defaultInterface,
		restart:          make(chan struct{}),
		start:            make(chan struct{}),
		stop:             make(chan struct{}),
	}
}

func (l *looper) Restart() { l.restart <- struct{}{} }
func (l *looper) Start()   { l.start <- struct{}{} }
func (l *looper) Stop()    { l.stop <- struct{}{} }

func (l *looper) GetSettings() (settings settings.ShadowSocks) {
	l.settingsMutex.RLock()
	defer l.settingsMutex.RUnlock()
	return l.settings
}

func (l *looper) SetSettings(settings settings.ShadowSocks) {
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

	l.setEnabled(true)

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
		server, err := shadowsockslib.NewServer(settings.Method, settings.Password, adaptLogger(l.logger, settings.Log))
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

		shadowsocksCtx, shadowsocksCancel := context.WithCancel(context.Background())

		waitError := make(chan error)
		go func() {
			waitError <- server.Listen(shadowsocksCtx, fmt.Sprintf("0.0.0.0:%d", settings.Port))
		}()
		if err != nil {
			shadowsocksCancel()
			l.logAndWait(ctx, err)
			continue
		}

		stayHere := true
		for stayHere {
			select {
			case <-ctx.Done():
				l.logger.Warn("context canceled: exiting loop")
				shadowsocksCancel()
				<-waitError
				close(waitError)
				return
			case <-l.restart: // triggered restart
				l.logger.Info("restarting")
				shadowsocksCancel()
				<-waitError
				close(waitError)
				stayHere = false
			case <-l.start:
				l.logger.Info("already started")
			case <-l.stop:
				l.logger.Info("stopping")
				shadowsocksCancel()
				<-waitError
				close(waitError)
				l.setEnabled(false)
				stayHere = false
			case err := <-waitError: // unexpected error
				shadowsocksCancel()
				close(waitError)
				l.logAndWait(ctx, err)
			}
		}
		shadowsocksCancel() // repetition for linter only
	}
}
