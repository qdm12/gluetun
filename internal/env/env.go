package env

import (
	"context"
	"os"

	"github.com/qdm12/golibs/logging"
)

type Env interface {
	FatalOnError(err error)
	PrintVersion(ctx context.Context, program string, commandFn func(ctx context.Context) (string, error))
}

type env struct {
	logger logging.Logger
	osExit func(n int)
}

func New(logger logging.Logger) Env {
	return &env{
		logger: logger,
		osExit: os.Exit,
	}
}

func (e *env) FatalOnError(err error) {
	if err != nil {
		e.logger.Error(err)
		e.osExit(1)
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
