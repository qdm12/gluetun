// Package cli defines an interface CLI to run command line operations.
package cli

import (
	"context"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/params"
)

type CLI interface {
	ClientKey(args []string) error
	HealthCheck(ctx context.Context, env params.Env, logger logging.Logger) error
	OpenvpnConfig(logger logging.Logger) error
	Update(ctx context.Context, args []string, logger logging.Logger) error
}

type cli struct {
	repoServersPath string
}

func New() CLI {
	return &cli{
		repoServersPath: "./internal/constants/servers.json",
	}
}
