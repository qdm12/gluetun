package openvpn

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
)

type Runner struct {
	settings configuration.OpenVPN
	starter  command.Starter
	logger   logging.Logger
}

func NewRunner(settings configuration.OpenVPN, starter command.Starter,
	logger logging.Logger) *Runner {
	return &Runner{
		starter:  starter,
		logger:   logger,
		settings: settings,
	}
}

func (r *Runner) Run(ctx context.Context, errCh chan<- error, ready chan<- struct{}) {
	stdoutLines, stderrLines, waitError, err := start(ctx, r.starter, r.settings.Version, r.settings.Flags)
	if err != nil {
		errCh <- err
		return
	}

	streamCtx, streamCancel := context.WithCancel(context.Background())
	streamDone := make(chan struct{})
	go streamLines(streamCtx, streamDone, r.logger,
		stdoutLines, stderrLines, ready)

	select {
	case <-ctx.Done():
		<-waitError
		close(waitError)
		streamCancel()
		<-streamDone
		errCh <- ctx.Err()
	case err := <-waitError:
		close(waitError)
		streamCancel()
		<-streamDone
		errCh <- err
	}
}
