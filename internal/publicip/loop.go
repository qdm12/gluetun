package publicip

import (
	"context"
	"time"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

type Looper interface {
	Run(ctx context.Context, restart <-chan struct{})
	RunRestartTicker(ctx context.Context, restart chan<- struct{})
}

type looper struct {
	getter           IPGetter
	logger           logging.Logger
	fileManager      files.FileManager
	ipStatusFilepath models.Filepath
	uid              int
	gid              int
}

func NewLooper(client network.Client, logger logging.Logger, fileManager files.FileManager,
	ipStatusFilepath models.Filepath, uid, gid int) Looper {
	return &looper{
		getter:           NewIPGetter(client),
		logger:           logger.WithPrefix("ip getter: "),
		fileManager:      fileManager,
		ipStatusFilepath: ipStatusFilepath,
		uid:              uid,
		gid:              gid,
	}
}

func (l *looper) logAndWait(err error) {
	l.logger.Error(err)
	l.logger.Info("retrying in 5 seconds")
	time.Sleep(5 * time.Second)
}

func (l *looper) Run(ctx context.Context, restart <-chan struct{}) {
	select {
	case <-restart:
	case <-ctx.Done():
		return
	}
	for {
		ip, err := l.getter.Get()
		if err != nil {
			l.logAndWait(err)
			continue
		}
		l.logger.Info("Public IP address is %s", ip)
		err = l.fileManager.WriteLinesToFile(
			string(l.ipStatusFilepath),
			[]string{ip.String()},
			files.Ownership(l.uid, l.gid),
			files.Permissions(0600))
		if err != nil {
			l.logAndWait(err)
			continue
		}
		select {
		case <-restart: // triggered restart
		case <-ctx.Done():
			l.logger.Warn("context canceled: exiting loop")
			return
		}
	}
}

func (l *looper) RunRestartTicker(ctx context.Context, restart chan<- struct{}) {
	ticker := time.NewTicker(time.Hour)
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
