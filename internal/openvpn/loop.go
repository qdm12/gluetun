package openvpn

import (
	"context"
	"sync"
	"time"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/settings"
)

type Looper interface {
	Run(ctx context.Context, restart <-chan struct{}, wg *sync.WaitGroup)
}

type looper struct {
	conf         Configurator
	settings     settings.OpenVPN
	logger       logging.Logger
	streamMerger command.StreamMerger
	fatalOnError func(err error)
	uid          int
	gid          int
}

func NewLooper(conf Configurator, settings settings.OpenVPN, logger logging.Logger,
	streamMerger command.StreamMerger, fatalOnError func(err error), uid, gid int) Looper {
	return &looper{
		conf:         conf,
		settings:     settings,
		logger:       logger.WithPrefix("openvpn: "),
		streamMerger: streamMerger,
		fatalOnError: fatalOnError,
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
		openvpnCtx, openvpnCancel := context.WithCancel(ctx)
		err := l.conf.WriteAuthFile(
			l.settings.User,
			l.settings.Password,
			l.uid,
			l.gid,
		)
		l.fatalOnError(err)
		stream, waitFn, err := l.conf.Start(openvpnCtx)
		l.fatalOnError(err)
		go l.streamMerger.Merge(openvpnCtx, stream,
			command.MergeName("openvpn"), command.MergeColor(constants.ColorOpenvpn()))
		waitError := make(chan error)
		go func() {
			err := waitFn() // blocking
			waitError <- err
		}()
		select {
		case <-ctx.Done():
			l.logger.Warn("context canceled: exiting loop")
			openvpnCancel()
			<-waitError
			close(waitError)
			return
		case <-restart: // triggered restart
			l.logger.Info("restarting")
			openvpnCancel()
			close(waitError)
		case err := <-waitError: // unexpected error
			l.logger.Warn(err)
			l.logger.Info("restarting")
			openvpnCancel()
			close(waitError)
			time.Sleep(time.Second)
		}
	}
}
