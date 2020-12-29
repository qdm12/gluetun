package publicip

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/os"
	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/golibs/logging"
)

type Looper interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	RunRestartTicker(ctx context.Context, wg *sync.WaitGroup)
	GetStatus() (status models.LoopStatus)
	SetStatus(status models.LoopStatus) (outcome string, err error)
	GetSettings() (settings settings.PublicIP)
	SetSettings(settings settings.PublicIP) (outcome string)
	GetPublicIP() (publicIP net.IP)
}

type looper struct {
	state state
	// Objects
	getter IPGetter
	logger logging.Logger
	os     os.OS
	// Fixed settings
	puid int
	pgid int
	// Internal channels and locks
	loopLock     sync.Mutex
	start        chan struct{}
	running      chan models.LoopStatus
	stop         chan struct{}
	stopped      chan struct{}
	updateTicker chan struct{}
	// Mock functions
	timeNow   func() time.Time
	timeSince func(time.Time) time.Duration
}

func NewLooper(client *http.Client, logger logging.Logger,
	settings settings.PublicIP, puid, pgid int,
	os os.OS) Looper {
	return &looper{
		state: state{
			status:   constants.Stopped,
			settings: settings,
		},
		// Objects
		getter:       NewIPGetter(client),
		logger:       logger.WithPrefix("ip getter: "),
		os:           os,
		puid:         puid,
		pgid:         pgid,
		start:        make(chan struct{}),
		running:      make(chan models.LoopStatus),
		stop:         make(chan struct{}),
		stopped:      make(chan struct{}),
		updateTicker: make(chan struct{}),
		timeNow:      time.Now,
		timeSince:    time.Since,
	}
}

func (l *looper) logAndWait(ctx context.Context, err error) {
	l.logger.Error(err)
	const waitTime = 5 * time.Second
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
		getCtx, getCancel := context.WithCancel(ctx)
		defer getCancel()

		ipCh := make(chan net.IP)
		errorCh := make(chan error)
		go func() {
			ip, err := l.getter.Get(getCtx)
			if err != nil {
				if getCtx.Err() == nil {
					errorCh <- err
				}
				return
			}
			ipCh <- ip
		}()

		if !crashed {
			l.running <- constants.Running
			crashed = false
		} else {
			l.state.setStatusWithLock(constants.Running)
		}

		stayHere := true
		for stayHere {
			select {
			case <-ctx.Done():
				l.logger.Warn("context canceled: exiting loop")
				getCancel()
				close(errorCh)
				filepath := l.GetSettings().IPFilepath
				l.logger.Info("Removing ip file %s", filepath)
				if err := l.os.Remove(string(filepath)); err != nil {
					l.logger.Error(err)
				}
				return
			case <-l.start:
				l.logger.Info("starting")
				getCancel()
				stayHere = false
			case <-l.stop:
				l.logger.Info("stopping")
				getCancel()
				<-errorCh
				l.stopped <- struct{}{}
			case ip := <-ipCh:
				getCancel()
				l.state.setPublicIP(ip)
				l.logger.Info("Public IP address is %s", ip)
				filepath := string(l.state.settings.IPFilepath)
				err := persistPublicIP(l.os.OpenFile, filepath, ip.String(), l.puid, l.pgid)
				if err != nil {
					l.logger.Error(err)
				}
				l.state.setStatusWithLock(constants.Completed)
			case err := <-errorCh:
				getCancel()
				close(ipCh)
				l.state.setStatusWithLock(constants.Crashed)
				l.logAndWait(ctx, err)
				crashed = true
				stayHere = false
			}
		}
		close(errorCh)
	}
}

func (l *looper) RunRestartTicker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	timer := time.NewTimer(time.Hour)
	timer.Stop() // 1 hour, cannot be a race condition
	timerIsStopped := true
	if period := l.GetSettings().Period; period > 0 {
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
			timer.Reset(l.GetSettings().Period)
		case <-l.updateTicker:
			if !timerIsStopped && !timer.Stop() {
				<-timer.C
			}
			timerIsStopped = true
			period := l.GetSettings().Period
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
