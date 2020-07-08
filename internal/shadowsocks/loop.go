package shadowsocks

import (
	"context"
	"sync"
	"time"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/firewall"
	"github.com/qdm12/private-internet-access-docker/internal/settings"
)

type Looper interface {
	Run(ctx context.Context, restart <-chan struct{}, wg *sync.WaitGroup)
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
}

func (l *looper) logAndWait(err error) {
	l.logger.Error(err)
	l.logger.Info("retrying in 1 minute")
	time.Sleep(time.Minute)
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
	}
}

func (l *looper) Run(ctx context.Context, restart <-chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	select {
	case <-restart:
	case <-ctx.Done():
		return
	}
	for {
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
			l.logAndWait(err)
			continue
		}
		err = l.firewallConf.AllowAnyIncomingOnPort(ctx, l.settings.Port)
		// TODO remove firewall rule on exit below
		if err != nil {
			l.logger.Error(err)
		}
		shadowsocksCtx, shadowsocksCancel := context.WithCancel(ctx)
		stdout, stderr, waitFn, err := l.conf.Start(ctx, "0.0.0.0", l.settings.Port, l.settings.Password, l.settings.Log)
		if err != nil {
			shadowsocksCancel()
			l.logAndWait(err)
			continue
		}
		go l.streamMerger.Merge(shadowsocksCtx, stdout,
			command.MergeName("shadowsocks"), command.MergeColor(constants.ColorShadowsocks()))
		go l.streamMerger.Merge(shadowsocksCtx, stderr,
			command.MergeName("shadowsocks error"), command.MergeColor(constants.ColorShadowsocksError()))
		waitError := make(chan error)
		go func() {
			err := waitFn() // blocking
			if shadowsocksCtx.Err() != context.Canceled {
				waitError <- err
			}
		}()
		select {
		case <-ctx.Done():
			l.logger.Warn("context canceled: exiting loop")
			shadowsocksCancel()
			close(waitError)
			return
		case <-restart: // triggered restart
			l.logger.Info("restarting")
			shadowsocksCancel()
			close(waitError)
		case err := <-waitError: // unexpected error
			l.logger.Warn(err)
			l.logger.Info("restarting")
			shadowsocksCancel()
			close(waitError)
			time.Sleep(time.Second)
		}
	}
}
