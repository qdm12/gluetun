package env

import (
	"context"

	"github.com/qdm12/golibs/logging"
)

type Env interface {
	FatalOnError(err error)
	PrintVersion(ctx context.Context, program string, commandFn func(ctx context.Context) (string, error))
}

type env struct {
	logger        logging.Logger
	cancelContext func()
}

func New(logger logging.Logger, cancelContext context.CancelFunc) Env {
	return &env{
		logger:        logger,
		cancelContext: cancelContext,
	}
}

func (e *env) FatalOnError(err error) {
	if err != nil {
		e.logger.Error(err)
		e.cancelContext()
	}
}

func (e *env) PrintVersion(ctx context.Context, program string, commandFn func(ctx context.Context) (string, error)) {
	version, err := commandFn(ctx)
	if err != nil {
		e.logger.Error(err)
	} else {
		e.logger.Info("%s version: %s", program, version)
	}
}
