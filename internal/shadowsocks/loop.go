package shadowsocks

import (
	"context"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	shadowsockslib "github.com/qdm12/ss-server/pkg/tcpudp"
)

type Loop struct {
	state state
	// Other objects
	logger Logger
	// Internal channels and locks
	loopLock      sync.Mutex
	running       chan models.LoopStatus
	stop, stopped chan struct{}
	start         chan struct{}
	backoffTime   time.Duration
}

func (l *Loop) logAndWait(ctx context.Context, err error) {
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

const defaultBackoffTime = 10 * time.Second

func NewLoop(settings settings.Shadowsocks, logger Logger) *Loop {
	return &Loop{
		state: state{
			status:   constants.Stopped,
			settings: settings,
		},
		logger:      logger,
		start:       make(chan struct{}),
		running:     make(chan models.LoopStatus),
		stop:        make(chan struct{}),
		stopped:     make(chan struct{}),
		backoffTime: defaultBackoffTime,
	}
}

func (l *Loop) Run(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	crashed := false

	if *l.GetSettings().Enabled {
		go func() {
			_, _ = l.SetStatus(ctx, constants.Running)
		}()
	}

	select {
	case <-l.start:
	case <-ctx.Done():
		return
	}

	for ctx.Err() == nil {
		settings := l.GetSettings()
		server, err := shadowsockslib.NewServer(settings.Settings, l.logger)
		if err != nil {
			crashed = true
			l.logAndWait(ctx, err)
			continue
		}

		shadowsocksCtx, shadowsocksCancel := context.WithCancel(ctx)

		waitError := make(chan error)
		go func() {
			waitError <- server.Listen(shadowsocksCtx)
		}()
		if err != nil {
			crashed = true
			shadowsocksCancel()
			l.logAndWait(ctx, err)
			continue
		}

		isStableTimer := time.NewTimer(time.Second)

		stayHere := true
		for stayHere {
			select {
			case <-ctx.Done():
				shadowsocksCancel()
				<-waitError
				close(waitError)
				return
			case <-isStableTimer.C:
				if !crashed {
					l.running <- constants.Running
					crashed = false
				} else {
					l.backoffTime = defaultBackoffTime
					l.state.setStatusWithLock(constants.Running)
				}
			case <-l.start:
				l.logger.Info("starting")
				shadowsocksCancel()
				<-waitError
				close(waitError)
				stayHere = false
			case <-l.stop:
				l.logger.Info("stopping")
				shadowsocksCancel()
				<-waitError
				close(waitError)
				l.stopped <- struct{}{}
			case err := <-waitError: // unexpected error
				shadowsocksCancel()
				close(waitError)
				if ctx.Err() != nil {
					return
				}
				l.state.setStatusWithLock(constants.Crashed)
				l.logAndWait(ctx, err)
				crashed = true
				stayHere = false
			}
		}
		shadowsocksCancel() // repetition for linter only
		if !isStableTimer.Stop() {
			<-isStableTimer.C
		}
	}
}
