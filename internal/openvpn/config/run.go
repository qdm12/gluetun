package config

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/golibs/logging"
)

type Runner interface {
	Run(ctx context.Context, errCh chan<- error, ready chan<- struct{},
		logger logging.Logger, settings configuration.OpenVPN)
}

func (c *Configurator) Run(ctx context.Context, errCh chan<- error,
	ready chan<- struct{}, logger logging.Logger, settings configuration.OpenVPN) {
	stdoutLines, stderrLines, waitError, err := c.start(ctx, settings.Version, settings.Flags)
	if err != nil {
		errCh <- err
		return
	}

	streamCtx, streamCancel := context.WithCancel(context.Background())
	streamDone := make(chan struct{})
	go streamLines(streamCtx, streamDone, logger,
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
