package routing

import (
	"context"
	"net"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
)

type Routing interface {
	AddRouteVia(ctx context.Context, subnet net.IPNet, defaultGateway net.IP, defaultInterface string) error
	DeleteRouteVia(ctx context.Context, subnet net.IPNet) (err error)
	DefaultRoute() (defaultInterface string, defaultGateway net.IP, defaultSubnet net.IPNet, err error)
	VPNGatewayIP(defaultInterface string) (ip net.IP, err error)
}

type routing struct {
	commander   command.Commander
	logger      logging.Logger
	fileManager files.FileManager
}

// NewConfigurator creates a new Configurator instance
func NewRouting(logger logging.Logger, fileManager files.FileManager) Routing {
	return &routing{
		commander:   command.NewCommander(),
		logger:      logger.WithPrefix("routing: "),
		fileManager: fileManager,
	}
}
