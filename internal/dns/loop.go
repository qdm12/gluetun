// Package dns defines interfaces to interact with DNS and DNS over TLS.
package dns

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/qdm12/dns/pkg/blacklist"
	"github.com/qdm12/dns/pkg/check"
	"github.com/qdm12/dns/pkg/nameserver"
	"github.com/qdm12/dns/pkg/unbound"
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
)

type Looper interface {
	Run(ctx context.Context, done chan<- struct{})
	RunRestartTicker(ctx context.Context, done chan<- struct{})
	GetStatus() (status models.LoopStatus)
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
	GetSettings() (settings configuration.DNS)
	SetSettings(ctx context.Context, settings configuration.DNS) (
		outcome string)
}

type looper struct {
	state        *state
	conf         unbound.Configurator
	blockBuilder blacklist.Builder
	client       *http.Client
	logger       logging.Logger
	start        <-chan struct{}
	running      chan<- models.LoopStatus
	stop         <-chan struct{}
	stopped      chan<- struct{}
	updateTicker <-chan struct{}
	backoffTime  time.Duration
	timeNow      func() time.Time
	timeSince    func(time.Time) time.Duration
	openFile     os.OpenFileFunc
}

const defaultBackoffTime = 10 * time.Second

func NewLooper(conf unbound.Configurator, settings configuration.DNS, client *http.Client,
	logger logging.Logger, openFile os.OpenFileFunc) Looper {
	start := make(chan struct{})
	running := make(chan models.LoopStatus)
	stop := make(chan struct{})
	stopped := make(chan struct{})
	updateTicker := make(chan struct{})

	state := newState(constants.Stopped, settings, start, running, stop, stopped, updateTicker)

	return &looper{
		state:        state,
		conf:         conf,
		blockBuilder: blacklist.NewBuilder(client),
		client:       client,
		logger:       logger,
		start:        start,
		running:      running,
		stop:         stop,
		stopped:      stopped,
		updateTicker: updateTicker,
		backoffTime:  defaultBackoffTime,
		timeNow:      time.Now,
		timeSince:    time.Since,
		openFile:     openFile,
	}
}

func (l *looper) logAndWait(ctx context.Context, err error) {
	if err != nil {
		l.logger.Warn(err)
	}
	l.logger.Info("attempting restart in " + l.backoffTime.String())
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

func (l *looper) signalOrSetStatus(userTriggered *bool, status models.LoopStatus) {
	if *userTriggered {
		*userTriggered = false
		l.running <- status
	} else {
		l.state.SetStatus(status)
	}
}

func (l *looper) Run(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	const fallback = false
	l.useUnencryptedDNS(fallback) // TODO remove? Use default DNS by default for Docker resolution?
	// TODO this one is kept if DNS_KEEP_NAMESERVER=on and should be replaced

	select {
	case <-l.start:
	case <-ctx.Done():
		return
	}

	userTriggered := true

	for ctx.Err() == nil {
		// Upper scope variables for Unbound only
		// Their values are to be used if DOT=off
		waitError := make(chan error)
		unboundCancel := func() { waitError <- nil }
		closeStreams := func() {}

		for l.GetSettings().Enabled {
			var err error
			unboundCancel, waitError, closeStreams, err = l.setupUnbound(ctx)
			if err == nil {
				l.backoffTime = defaultBackoffTime
				l.logger.Info("ready")
				l.signalOrSetStatus(&userTriggered, constants.Running)
				break
			}

			l.signalOrSetStatus(&userTriggered, constants.Crashed)

			if ctx.Err() != nil {
				return
			}

			if !errors.Is(err, errUpdateFiles) {
				const fallback = true
				l.useUnencryptedDNS(fallback)
			}
			l.logAndWait(ctx, err)
		}

		if !l.GetSettings().Enabled {
			const fallback = false
			l.useUnencryptedDNS(fallback)
		}

		userTriggered = false

		stayHere := true
		for stayHere {
			select {
			case <-ctx.Done():
				unboundCancel()
				<-waitError
				close(waitError)
				closeStreams()
				return
			case <-l.stop:
				userTriggered = true
				l.logger.Info("stopping")
				const fallback = false
				l.useUnencryptedDNS(fallback)
				unboundCancel()
				<-waitError
				// do not close waitError or the waitError
				// select case will trigger
				closeStreams()
				l.stopped <- struct{}{}
			case <-l.start:
				userTriggered = true
				l.logger.Info("starting")
				stayHere = false
			case err := <-waitError: // unexpected error
				close(waitError)
				closeStreams()

				l.state.Lock() // prevent SetStatus from running in parallel

				unboundCancel()
				l.state.SetStatus(constants.Crashed)
				const fallback = true
				l.useUnencryptedDNS(fallback)
				l.logAndWait(ctx, err)
				stayHere = false

				l.state.Unlock()
			}
		}
	}
}

var errUpdateFiles = errors.New("cannot update files")

// Returning cancel == nil signals we want to re-run setupUnbound
// Returning err == errUpdateFiles signals we should not fall back
// on the plaintext DNS as DOT is still up and running.
func (l *looper) setupUnbound(ctx context.Context) (
	cancel context.CancelFunc, waitError chan error, closeStreams func(), err error) {
	err = l.updateFiles(ctx)
	if err != nil {
		return nil, nil, nil, errUpdateFiles
	}

	settings := l.GetSettings()

	unboundCtx, cancel := context.WithCancel(context.Background())
	stdoutLines, stderrLines, waitError, err := l.conf.Start(unboundCtx, settings.Unbound.VerbosityDetailsLevel)
	if err != nil {
		cancel()
		return nil, nil, nil, err
	}

	collectLinesDone := make(chan struct{})
	go l.collectLines(stdoutLines, stderrLines, collectLinesDone)
	closeStreams = func() {
		close(stdoutLines)
		close(stderrLines)
		<-collectLinesDone
	}

	// use Unbound
	nameserver.UseDNSInternally(net.IP{127, 0, 0, 1})
	err = nameserver.UseDNSSystemWide(l.openFile,
		net.IP{127, 0, 0, 1}, settings.KeepNameserver)
	if err != nil {
		l.logger.Error(err)
	}

	if err := check.WaitForDNS(ctx, net.DefaultResolver); err != nil {
		cancel()
		<-waitError
		close(waitError)
		closeStreams()
		return nil, nil, nil, err
	}

	return cancel, waitError, closeStreams, nil
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
		nameserver.UseDNSInternally(targetIP)
		if err := nameserver.UseDNSSystemWide(l.openFile,
			targetIP, settings.KeepNameserver); err != nil {
			l.logger.Error(err)
		}
		return
	}

	provider := settings.Unbound.Providers[0]
	targetIP = provider.DoT().IPv4[0]
	if fallback {
		l.logger.Info("falling back on plaintext DNS at address " + targetIP.String())
	} else {
		l.logger.Info("using plaintext DNS at address " + targetIP.String())
	}
	nameserver.UseDNSInternally(targetIP)
	if err := nameserver.UseDNSSystemWide(l.openFile, targetIP, settings.KeepNameserver); err != nil {
		l.logger.Error(err)
	}
}

func (l *looper) RunRestartTicker(ctx context.Context, done chan<- struct{}) {
	defer close(done)
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
					l.state.SetStatus(constants.Crashed)
					l.logger.Error(err)
					l.logger.Warn("skipping Unbound restart due to failed files update")
					continue
				}
			}

			_, _ = l.ApplyStatus(ctx, constants.Stopped)
			_, _ = l.ApplyStatus(ctx, constants.Running)

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
	l.logger.Info("downloading DNS over TLS cryptographic files")
	if err := l.conf.SetupFiles(ctx); err != nil {
		return err
	}
	settings := l.GetSettings()

	l.logger.Info("downloading hostnames and IP block lists")
	blockedHostnames, blockedIPs, blockedIPPrefixes, errs := l.blockBuilder.All(
		ctx, settings.BlacklistBuild)
	for _, err := range errs {
		l.logger.Warn(err)
	}

	// TODO change to BlockHostnames() when migrating to qdm12/dns v2
	settings.Unbound.Blacklist.FqdnHostnames = blockedHostnames
	settings.Unbound.Blacklist.IPs = blockedIPs
	settings.Unbound.Blacklist.IPPrefixes = blockedIPPrefixes

	return l.conf.MakeUnboundConf(settings.Unbound)
}

func (l *looper) GetStatus() (status models.LoopStatus) { return l.state.GetStatus() }
func (l *looper) ApplyStatus(ctx context.Context, status models.LoopStatus) (
	outcome string, err error) {
	return l.state.ApplyStatus(ctx, status)
}
func (l *looper) GetSettings() (settings configuration.DNS) { return l.state.GetSettings() }
func (l *looper) SetSettings(ctx context.Context, settings configuration.DNS) (
	outcome string) {
	return l.state.SetSettings(ctx, settings)
}
