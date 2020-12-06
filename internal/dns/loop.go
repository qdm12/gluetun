package dns

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
)

type Looper interface {
	Run(ctx context.Context, wg *sync.WaitGroup, signalDNSReady func())
	RunRestartTicker(ctx context.Context, wg *sync.WaitGroup)
	GetStatus() (status models.LoopStatus)
	SetStatus(status models.LoopStatus) (outcome string, err error)
	GetSettings() (settings settings.DNS)
	SetSettings(settings settings.DNS) (outcome string)
}

type looper struct {
	state        state
	conf         Configurator
	logger       logging.Logger
	streamMerger command.StreamMerger
	uid          int
	gid          int
	loopLock     sync.Mutex
	start        chan struct{}
	running      chan models.LoopStatus
	stop         chan struct{}
	stopped      chan struct{}
	updateTicker chan struct{}
	timeNow      func() time.Time
	timeSince    func(time.Time) time.Duration
}

func NewLooper(conf Configurator, settings settings.DNS, logger logging.Logger,
	streamMerger command.StreamMerger, uid, gid int) Looper {
	return &looper{
		state: state{
			status:   constants.Stopped,
			settings: settings,
		},
		conf:         conf,
		logger:       logger.WithPrefix("dns over tls: "),
		uid:          uid,
		gid:          gid,
		streamMerger: streamMerger,
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
	l.logger.Warn(err)
	l.logger.Info("attempting restart in 10 seconds")
	const waitDuration = 10 * time.Second
	timer := time.NewTimer(waitDuration)
	select {
	case <-timer.C:
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
	}
}

func (l *looper) Run(ctx context.Context, wg *sync.WaitGroup, signalDNSReady func()) {
	defer wg.Done()

	const fallback = false
	l.useUnencryptedDNS(fallback) // TODO remove? Use default DNS by default for Docker resolution?

	select {
	case <-l.start:
	case <-ctx.Done():
		return
	}

	defer l.logger.Warn("loop exited")

	for ctx.Err() == nil {
		err := l.updateFiles(ctx)
		if err == nil {
			break
		}
		l.state.setStatusWithLock(constants.Crashed)
		l.logAndWait(ctx, err)
	}

	crashed := false

	for ctx.Err() == nil {
		settings := l.GetSettings()

		unboundCtx, unboundCancel := context.WithCancel(context.Background())
		stream, waitFn, err := l.conf.Start(unboundCtx, settings.VerbosityDetailsLevel)
		if err != nil {
			unboundCancel()
			if !crashed {
				l.running <- constants.Crashed
			}
			crashed = true
			const fallback = true
			l.useUnencryptedDNS(fallback)
			l.logAndWait(ctx, err)
			continue
		}

		// Started successfully
		go l.streamMerger.Merge(unboundCtx, stream, command.MergeName("unbound"))

		l.conf.UseDNSInternally(net.IP{127, 0, 0, 1})                                                  // use Unbound
		if err := l.conf.UseDNSSystemWide(net.IP{127, 0, 0, 1}, settings.KeepNameserver); err != nil { // use Unbound
			l.logger.Error(err)
		}

		if err := l.conf.WaitForUnbound(); err != nil {
			if !crashed {
				l.running <- constants.Crashed
				crashed = true
			}
			unboundCancel()
			const fallback = true
			l.useUnencryptedDNS(fallback)
			l.logAndWait(ctx, err)
			continue
		}

		waitError := make(chan error)
		go func() {
			err := waitFn() // blocking
			waitError <- err
		}()

		l.logger.Info("DNS over TLS is ready")
		if !crashed {
			l.running <- constants.Running
			crashed = false
		} else {
			l.state.setStatusWithLock(constants.Running)
		}
		signalDNSReady()

		stayHere := true
		for stayHere {
			select {
			case <-ctx.Done():
				l.logger.Warn("context canceled: exiting loop")
				unboundCancel()
				<-waitError
				close(waitError)
				return
			case <-l.stop:
				l.logger.Info("stopping")
				const fallback = false
				l.useUnencryptedDNS(fallback)
				unboundCancel()
				<-waitError
				l.stopped <- struct{}{}
			case <-l.start:
				l.logger.Info("starting")
				stayHere = false
			case err := <-waitError: // unexpected error
				unboundCancel()
				l.state.setStatusWithLock(constants.Crashed)
				const fallback = true
				l.useUnencryptedDNS(fallback)
				l.logAndWait(ctx, err)
				stayHere = false
			}
		}
		close(waitError)
		unboundCancel()
	}
}

func (l *looper) useUnencryptedDNS(fallback bool) {
	settings := l.GetSettings()

	// Try with user provided plaintext ip address
	targetIP := settings.PlaintextAddress
	if targetIP != nil {
		if fallback {
			l.logger.Info("falling back on plaintext DNS at address %s", targetIP)
		} else {
			l.logger.Info("using plaintext DNS at address %s", targetIP)
		}
		l.conf.UseDNSInternally(targetIP)
		if err := l.conf.UseDNSSystemWide(targetIP, settings.KeepNameserver); err != nil {
			l.logger.Error(err)
		}
		return
	}

	// Try with any IPv4 address from the providers chosen
	for _, provider := range settings.Providers {
		data := constants.DNSProviderMapping()[provider]
		for _, targetIP = range data.IPs {
			if targetIP.To4() != nil {
				if fallback {
					l.logger.Info("falling back on plaintext DNS at address %s", targetIP)
				} else {
					l.logger.Info("using plaintext DNS at address %s", targetIP)
				}
				l.conf.UseDNSInternally(targetIP)
				if err := l.conf.UseDNSSystemWide(targetIP, settings.KeepNameserver); err != nil {
					l.logger.Error(err)
				}
				return
			}
		}
	}

	// No IPv4 address found
	l.logger.Error("no ipv4 DNS address found for providers %s", settings.Providers)
}

func (l *looper) RunRestartTicker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	// Timer that acts as a ticker
	timer := time.NewTimer(time.Hour)
	timer.Stop()
	timerIsStopped := true
	settings := l.GetSettings()
	if settings.UpdatePeriod > 0 {
		timer.Reset(settings.UpdatePeriod)
		timerIsStopped = false
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

			status := l.GetStatus()
			if status == constants.Running {
				if err := l.updateFiles(ctx); err != nil {
					l.state.setStatusWithLock(constants.Crashed)
					l.logger.Error(err)
					l.logger.Warn("skipping Unbound restart due to failed files update")
					continue
				}
			}

			_, _ = l.SetStatus(constants.Stopped)
			_, _ = l.SetStatus(constants.Running)

			settings := l.GetSettings()
			timer.Reset(settings.UpdatePeriod)
		case <-l.updateTicker:
			if !timer.Stop() {
				<-timer.C
			}
			timerIsStopped = true
			settings := l.GetSettings()
			newUpdatePeriod := settings.UpdatePeriod
			if newUpdatePeriod == 0 {
				continue
			}
			var waited time.Duration
			if lastTick.UnixNano() != 0 {
				waited = l.timeSince(lastTick)
			}
			leftToWait := newUpdatePeriod - waited
			timer.Reset(leftToWait)
			timerIsStopped = false
		}
	}
}

func (l *looper) updateFiles(ctx context.Context) (err error) {
	if err := l.conf.DownloadRootHints(ctx, l.uid, l.gid); err != nil {
		return err
	}
	if err := l.conf.DownloadRootKey(ctx, l.uid, l.gid); err != nil {
		return err
	}
	settings := l.GetSettings()
	if err := l.conf.MakeUnboundConf(ctx, settings, l.uid, l.gid); err != nil {
		return err
	}
	return nil
}
