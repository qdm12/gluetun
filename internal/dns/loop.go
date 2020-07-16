package dns

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/settings"
)

type Looper interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	RunRestartTicker(ctx context.Context)
	Restart()
	Start()
	Stop()
	GetSettings() (settings settings.DNS)
	SetSettings(settings settings.DNS)
}

type looper struct {
	conf             Configurator
	settings         settings.DNS
	settingsMutex    sync.RWMutex
	logger           logging.Logger
	streamMerger     command.StreamMerger
	uid              int
	gid              int
	restart          chan struct{}
	start            chan struct{}
	stop             chan struct{}
	updateTicker     chan struct{}
	tickerReady      bool
	tickerReadyMutex sync.Mutex
}

func NewLooper(conf Configurator, settings settings.DNS, logger logging.Logger,
	streamMerger command.StreamMerger, uid, gid int) Looper {
	return &looper{
		conf:         conf,
		settings:     settings,
		logger:       logger.WithPrefix("dns over tls: "),
		uid:          uid,
		gid:          gid,
		streamMerger: streamMerger,
		restart:      make(chan struct{}),
		start:        make(chan struct{}),
		stop:         make(chan struct{}),
		updateTicker: make(chan struct{}),
	}
}

func (l *looper) Restart() { l.restart <- struct{}{} }
func (l *looper) Start()   { l.start <- struct{}{} }
func (l *looper) Stop()    { l.stop <- struct{}{} }

func (l *looper) GetSettings() (settings settings.DNS) {
	l.settingsMutex.RLock()
	defer l.settingsMutex.RUnlock()
	return l.settings
}

func (l *looper) SetSettings(settings settings.DNS) {
	l.settingsMutex.Lock()
	defer l.settingsMutex.Unlock()
	updatePeriodDiffers := l.settings.UpdatePeriod != settings.UpdatePeriod
	l.settings = settings
	l.settingsMutex.Unlock()
	if updatePeriodDiffers {
		l.updateTicker <- struct{}{}
	}
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

func (l *looper) logAndWait(ctx context.Context, err error) {
	l.logger.Warn(err)
	l.logger.Info("attempting restart in 10 seconds")
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	<-ctx.Done()
}

func (l *looper) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	l.fallbackToUnencryptedDNS()
	waitForStart := true
	for waitForStart {
		select {
		case <-l.stop:
			l.logger.Info("not started yet")
		case <-l.restart:
			waitForStart = false
		case <-l.start:
			waitForStart = false
		case <-ctx.Done():
			return
		}
	}
	defer l.logger.Warn("loop exited")

	var unboundCtx context.Context
	var unboundCancel context.CancelFunc = func() {}
	var waitError chan error
	triggeredRestart := false
	l.setEnabled(true)
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
				unboundCancel()
				return
			}
		}

		settings := l.GetSettings()

		// Setup
		if err := l.conf.DownloadRootHints(l.uid, l.gid); err != nil {
			l.logAndWait(ctx, err)
			continue
		}
		if err := l.conf.DownloadRootKey(l.uid, l.gid); err != nil {
			l.logAndWait(ctx, err)
			continue
		}
		if err := l.conf.MakeUnboundConf(settings, l.uid, l.gid); err != nil {
			l.logAndWait(ctx, err)
			continue
		}

		if triggeredRestart {
			triggeredRestart = false
			unboundCancel()
			<-waitError
			close(waitError)
		}
		unboundCtx, unboundCancel = context.WithCancel(context.Background())
		stream, waitFn, err := l.conf.Start(unboundCtx, settings.VerbosityDetailsLevel)
		if err != nil {
			unboundCancel()
			l.fallbackToUnencryptedDNS()
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
			unboundCancel()
			l.fallbackToUnencryptedDNS()
			l.logAndWait(ctx, err)
			continue
		}
		waitError = make(chan error)
		go func() {
			err := waitFn() // blocking
			waitError <- err
		}()

		stayHere := true
		for stayHere {
			select {
			case <-ctx.Done():
				l.logger.Warn("context canceled: exiting loop")
				unboundCancel()
				<-waitError
				close(waitError)
				return
			case <-l.restart: // triggered restart
				l.logger.Info("restarting")
				// unboundCancel occurs next loop run when the setup is complete
				triggeredRestart = true
				stayHere = false
			case <-l.start:
				l.logger.Info("already started")
			case <-l.stop:
				l.logger.Info("stopping")
				unboundCancel()
				<-waitError
				close(waitError)
				l.setEnabled(false)
				stayHere = false
			case err := <-waitError: // unexpected error
				close(waitError)
				unboundCancel()
				l.fallbackToUnencryptedDNS()
				l.logAndWait(ctx, err)
				stayHere = false
			}
		}
	}
	unboundCancel()
}

func (l *looper) fallbackToUnencryptedDNS() {
	settings := l.GetSettings()

	// Try with user provided plaintext ip address
	targetIP := settings.PlaintextAddress
	if targetIP != nil {
		l.logger.Info("falling back on plaintext DNS at address %s", targetIP)
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
				l.logger.Info("falling back on plaintext DNS at address %s", targetIP)
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

func (l *looper) RunRestartTicker(ctx context.Context) {
	l.tickerReadyMutex.Lock()
	l.tickerReady = true
	l.tickerReadyMutex.Unlock()
	var ticker *time.Ticker = nil
	settings := l.GetSettings()
	if settings.UpdatePeriod > 0 {
		ticker = time.NewTicker(settings.UpdatePeriod)
	}
	for {
		select {
		case <-ctx.Done():
			if ticker != nil {
				ticker.Stop()
			}
			return
		case <-ticker.C:
			l.restart <- struct{}{}
		case <-l.updateTicker:
			if ticker != nil {
				ticker.Stop()
			}
			period := l.GetSettings().UpdatePeriod
			ticker = time.NewTicker(period)
		}
	}
}
