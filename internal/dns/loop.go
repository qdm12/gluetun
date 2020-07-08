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
	Run(ctx context.Context, restart <-chan struct{}, wg *sync.WaitGroup)
	RunRestartTicker(ctx context.Context, restart chan<- struct{})
}

type looper struct {
	conf         Configurator
	settings     settings.DNS
	logger       logging.Logger
	streamMerger command.StreamMerger
	uid          int
	gid          int
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
	}
}

func (l *looper) attemptingRestart(err error) {
	l.logger.Warn(err)
	l.logger.Info("attempting restart in 10 seconds")
	time.Sleep(10 * time.Second)
}

func (l *looper) Run(ctx context.Context, restart <-chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	l.fallbackToUnencryptedDNS()
	select {
	case <-restart:
	case <-ctx.Done():
		return
	}
	_, unboundCancel := context.WithCancel(ctx)
	for {
		if !l.settings.Enabled {
			// wait for another restart signal to recheck if it is enabled
			select {
			case <-restart:
			case <-ctx.Done():
				unboundCancel()
				return
			}
		}
		if ctx.Err() == context.Canceled {
			unboundCancel()
			return
		}

		// Setup
		if err := l.conf.DownloadRootHints(l.uid, l.gid); err != nil {
			l.attemptingRestart(err)
			continue
		}
		if err := l.conf.DownloadRootKey(l.uid, l.gid); err != nil {
			l.attemptingRestart(err)
			continue
		}
		if err := l.conf.MakeUnboundConf(l.settings, l.uid, l.gid); err != nil {
			l.attemptingRestart(err)
			continue
		}

		// Start command
		unboundCancel()
		unboundCtx, unboundCancel := context.WithCancel(ctx)
		stream, waitFn, err := l.conf.Start(unboundCtx, l.settings.VerbosityDetailsLevel)
		if err != nil {
			unboundCancel()
			l.fallbackToUnencryptedDNS()
			l.attemptingRestart(err)
		}

		// Started successfully
		go l.streamMerger.Merge(unboundCtx, stream,
			command.MergeName("unbound"), command.MergeColor(constants.ColorUnbound()))
		l.conf.UseDNSInternally(net.IP{127, 0, 0, 1})                         // use Unbound
		if err := l.conf.UseDNSSystemWide(net.IP{127, 0, 0, 1}); err != nil { // use Unbound
			unboundCancel()
			l.fallbackToUnencryptedDNS()
			l.attemptingRestart(err)
		}
		if err := l.conf.WaitForUnbound(); err != nil {
			unboundCancel()
			l.fallbackToUnencryptedDNS()
			l.attemptingRestart(err)
		}
		waitError := make(chan error)
		go func() {
			err := waitFn() // blocking
			if unboundCtx.Err() != context.Canceled {
				waitError <- err
			}
		}()

		// Wait for one of the three cases below
		select {
		case <-ctx.Done():
			l.logger.Warn("context canceled: exiting loop")
			unboundCancel()
			close(waitError)
			return
		case <-restart: // triggered restart
			unboundCancel()
			close(waitError)
			l.logger.Info("restarting")
		case err := <-waitError: // unexpected error
			unboundCancel()
			close(waitError)
			l.fallbackToUnencryptedDNS()
			l.attemptingRestart(err)
		}
	}
}

func (l *looper) fallbackToUnencryptedDNS() {
	// Try with user provided plaintext ip address
	targetIP := l.settings.PlaintextAddress
	if targetIP != nil {
		l.conf.UseDNSInternally(targetIP)
		if err := l.conf.UseDNSSystemWide(targetIP); err != nil {
			l.logger.Error(err)
		}
		return
	}

	// Try with any IPv4 address from the providers chosen
	for _, provider := range l.settings.Providers {
		data := constants.DNSProviderMapping()[provider]
		for _, targetIP = range data.IPs {
			if targetIP.To4() != nil {
				l.conf.UseDNSInternally(targetIP)
				if err := l.conf.UseDNSSystemWide(targetIP); err != nil {
					l.logger.Error(err)
				}
				return
			}
		}
	}

	// No IPv4 address found
	l.logger.Error("no ipv4 DNS address found for providers %s", l.settings.Providers)
}

func (l *looper) RunRestartTicker(ctx context.Context, restart chan<- struct{}) {
	if l.settings.UpdatePeriod == 0 {
		return
	}
	ticker := time.NewTicker(l.settings.UpdatePeriod)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			restart <- struct{}{}
		}
	}
}
