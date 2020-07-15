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
	Run(ctx context.Context)
	RunRestartTicker(ctx context.Context)
	Restart()
}

type looper struct {
	getter           IPGetter
	logger           logging.Logger
	fileManager      files.FileManager
	ipStatusFilepath models.Filepath
	uid              int
	gid              int
	restart          chan struct{}
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
		restart:          make(chan struct{}),
	}
}

func (l *looper) Restart() { l.restart <- struct{}{} }

func (l *looper) logAndWait(ctx context.Context, err error) {
	l.logger.Error(err)
	l.logger.Info("retrying in 5 seconds")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel() // just for the linter
	<-ctx.Done()
}

func (l *looper) Run(ctx context.Context) {
	select {
	case <-l.restart:
	case <-ctx.Done():
		return
	}
	defer l.logger.Warn("loop exited")

	for ctx.Err() == nil {
		ip, err := l.getter.Get()
		if err != nil {
			l.logAndWait(ctx, err)
			continue
		}
		l.logger.Info("Public IP address is %s", ip)
		err = l.fileManager.WriteLinesToFile(
			string(l.ipStatusFilepath),
			[]string{ip.String()},
			files.Ownership(l.uid, l.gid),
			files.Permissions(0600))
		if err != nil {
			l.logAndWait(ctx, err)
			continue
		}
		select {
		case <-l.restart: // triggered restart
		case <-ctx.Done():
			l.logger.Warn("context canceled: exiting loop")
			return
		}
	}
}

func (l *looper) RunRestartTicker(ctx context.Context) {
	ticker := time.NewTicker(time.Hour)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			l.restart <- struct{}{}
		}
	}
}
