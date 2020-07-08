package tinyproxy

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
	settings     settings.TinyProxy
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

func NewLooper(conf Configurator, firewallConf firewall.Configurator, settings settings.TinyProxy,
	logger logging.Logger, streamMerger command.StreamMerger, uid, gid int) Looper {
	return &looper{
		conf:         conf,
		firewallConf: firewallConf,
		settings:     settings,
		logger:       logger.WithPrefix("tinyproxy: "),
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
		err := l.conf.MakeConf(
			l.settings.LogLevel,
			l.settings.Port,
			l.settings.User,
			l.settings.Password,
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
		tinyproxyCtx, tinyproxyCancel := context.WithCancel(ctx)
		stream, waitFn, err := l.conf.Start(tinyproxyCtx)
		if err != nil {
			tinyproxyCancel()
			l.logAndWait(err)
			continue
		}
		go l.streamMerger.Merge(tinyproxyCtx, stream,
			command.MergeName("tinyproxy"), command.MergeColor(constants.ColorTinyproxy()))
		waitError := make(chan error)
		go func() {
			err := waitFn() // blocking
			if tinyproxyCtx.Err() != context.Canceled {
				waitError <- err
			}
		}()
		select {
		case <-ctx.Done():
			l.logger.Warn("context canceled: exiting loop")
			tinyproxyCancel()
			close(waitError)
			return
		case <-restart: // triggered restart
			l.logger.Info("restarting")
			tinyproxyCancel()
			close(waitError)
		case err := <-waitError: // unexpected error
			l.logger.Warn(err)
			l.logger.Info("restarting")
			tinyproxyCancel()
			close(waitError)
			time.Sleep(time.Second)
		}
	}
}
