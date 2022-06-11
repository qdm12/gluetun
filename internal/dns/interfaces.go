package dns

import (
	"context"

	"github.com/qdm12/dns/pkg/unbound"
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
)

type Configurator interface {
	SetupFiles(ctx context.Context) error
	MakeUnboundConf(settings unbound.Settings) (err error)
	Start(ctx context.Context, verbosityDetailsLevel uint8) (
		stdoutLines, stderrLines chan string, waitError chan error, err error)
	Version(ctx context.Context) (version string, err error)
}

type statusManager interface {
	GetStatus() (status models.LoopStatus)
	SetStatus(status models.LoopStatus)
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
}

type stateManager interface {
	GetSettings() (settings settings.DNS)
	SetSettings(ctx context.Context, settings settings.DNS) (outcome string)
}
