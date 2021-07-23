package publicip

import (
	"context"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/logging"
)

type Looper interface {
	Run(ctx context.Context, done chan<- struct{})
	RunRestartTicker(ctx context.Context, done chan<- struct{})
	GetStatus() (status models.LoopStatus)
	SetStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
	GetSettings() (settings configuration.PublicIP)
	SetSettings(settings configuration.PublicIP) (outcome string)
	GetPublicIP() (publicIP net.IP)
}

type looper struct {
	state state
	// Objects
	getter IPGetter
	client *http.Client
	logger logging.Logger
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
	backoffTime  time.Duration
	// Mock functions
	timeNow   func() time.Time
	timeSince func(time.Time) time.Duration
}

const defaultBackoffTime = 5 * time.Second

func NewLooper(client *http.Client, logger logging.Logger,
	settings configuration.PublicIP, puid, pgid int) Looper {
	return &looper{
		state: state{
			status:   constants.Stopped,
			settings: settings,
		},
		// Objects
		client:       client,
		getter:       NewIPGetter(client),
		logger:       logger,
		puid:         puid,
		pgid:         pgid,
		start:        make(chan struct{}),
		running:      make(chan models.LoopStatus),
		stop:         make(chan struct{}),
		stopped:      make(chan struct{}),
		updateTicker: make(chan struct{}),
		backoffTime:  defaultBackoffTime,
		timeNow:      time.Now,
		timeSince:    time.Since,
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
			l.backoffTime = defaultBackoffTime
			l.state.setStatusWithLock(constants.Running)
		}

		stayHere := true
		for stayHere {
			select {
			case <-ctx.Done():
				getCancel()
				close(errorCh)
				filepath := l.GetSettings().IPFilepath
				l.logger.Info("Removing ip file " + filepath)
				if err := os.Remove(filepath); err != nil {
					l.logger.Error(err.Error())
				}
				return
			case <-l.start:
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

				message := "Public IP address is " + ip.String()
				result, err := Info(ctx, l.client, ip)
				if err != nil {
					l.logger.Warn(err.Error())
				} else {
					message += " (" + result.Country + ", " + result.Region + ", " + result.City + ")"
				}
				l.logger.Info(message)

				err = persistPublicIP(l.state.settings.IPFilepath,
					ip.String(), l.puid, l.pgid)
				if err != nil {
					l.logger.Error(err.Error())
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

func (l *looper) RunRestartTicker(ctx context.Context, done chan<- struct{}) {
	defer close(done)
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
