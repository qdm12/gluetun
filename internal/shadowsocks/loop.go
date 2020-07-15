package shadowsocks

import (
	"context"
	"sync"
	"time"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/private-internet-access-docker/internal/firewall"
	"github.com/qdm12/private-internet-access-docker/internal/settings"
)

type Looper interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	Restart()
}

type looper struct {
	conf         Configurator
	firewallConf firewall.Configurator
	settings     settings.ShadowSocks
	dnsSettings  settings.DNS
	logger       logging.Logger
	streamMerger command.StreamMerger
	uid          int
	gid          int
	restart      chan struct{}
}

func (l *looper) logAndWait(ctx context.Context, err error) {
	l.logger.Error(err)
	l.logger.Info("retrying in 1 minute")
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel() // just for the linter
	<-ctx.Done()
}

func NewLooper(conf Configurator, firewallConf firewall.Configurator, settings settings.ShadowSocks, dnsSettings settings.DNS,
	logger logging.Logger, streamMerger command.StreamMerger, uid, gid int) Looper {
	return &looper{
		conf:         conf,
		firewallConf: firewallConf,
		settings:     settings,
		dnsSettings:  dnsSettings,
		logger:       logger.WithPrefix("shadowsocks: "),
		streamMerger: streamMerger,
		uid:          uid,
		gid:          gid,
		restart:      make(chan struct{}),
	}
}

func (l *looper) Restart() { l.restart <- struct{}{} }

func (l *looper) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	select {
	case <-l.restart:
	case <-ctx.Done():
		return
	}
	defer l.logger.Warn("loop exited")

	var previousPort uint16
	for ctx.Err() == nil {
		nameserver := l.dnsSettings.PlaintextAddress.String()
		if l.dnsSettings.Enabled {
			nameserver = "127.0.0.1"
		}
		err := l.conf.MakeConf(
			l.settings.Port,
			l.settings.Password,
			l.settings.Method,
			nameserver,
			l.uid,
			l.gid)
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
		if err := l.firewallConf.SetAllowedPort(ctx, l.settings.Port); err != nil {
			l.logger.Error(err)
			continue
		}
		previousPort = l.settings.Port

		shadowsocksCtx, shadowsocksCancel := context.WithCancel(context.Background())
		stdout, stderr, waitFn, err := l.conf.Start(shadowsocksCtx, "0.0.0.0", l.settings.Port, l.settings.Password, l.settings.Log)
		if err != nil {
			shadowsocksCancel()
			l.logAndWait(ctx, err)
			continue
		}
		go l.streamMerger.Merge(shadowsocksCtx, stdout, command.MergeName("shadowsocks"))
		go l.streamMerger.Merge(shadowsocksCtx, stderr, command.MergeName("shadowsocks error"))
		waitError := make(chan error)
		go func() {
			err := waitFn() // blocking
			waitError <- err
		}()
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
		case err := <-waitError: // unexpected error
			shadowsocksCancel()
			close(waitError)
			l.logAndWait(ctx, err)
		}
	}
}
