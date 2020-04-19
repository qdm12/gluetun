package firewall

import (
	"context"
	"net"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// Configurator allows to change firewall rules and modify network routes
type Configurator interface {
	Version(ctx context.Context) (string, error)
	AcceptAll(ctx context.Context) error
	Clear(ctx context.Context) error
	BlockAll(ctx context.Context) error
	CreateGeneralRules(ctx context.Context) error
	CreateVPNRules(ctx context.Context, dev models.VPNDevice, defaultInterface string, connections []models.OpenVPNConnection) error
	CreateLocalSubnetsRules(ctx context.Context, subnet net.IPNet, extraSubnets []net.IPNet, defaultInterface string) error
	AllowInputTrafficOnPort(ctx context.Context, device models.VPNDevice, port uint16) error
	AllowAnyIncomingOnPort(ctx context.Context, port uint16) error
}

type configurator struct {
	commander command.Commander
	logger    logging.Logger
}

// NewConfigurator creates a new Configurator instance
func NewConfigurator(logger logging.Logger) Configurator {
	return &configurator{
		commander: command.NewCommander(),
		logger:    logger.WithPrefix("firewall configurator: "),
	}
}
