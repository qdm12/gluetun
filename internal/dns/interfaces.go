package dns

import (
	"context"

	"github.com/qdm12/dns/pkg/unbound"
)

type Configurator interface {
	SetupFiles(ctx context.Context) error
	MakeUnboundConf(settings unbound.Settings) (err error)
	Start(ctx context.Context, verbosityDetailsLevel uint8) (
		stdoutLines, stderrLines chan string, waitError chan error, err error)
	Version(ctx context.Context) (version string, err error)
}
