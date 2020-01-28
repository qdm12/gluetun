package dns

import (
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/settings"
)

type Configurator interface {
	MakeUnboundConf(settings settings.DNS) (err error)
	SetLocalNameserver() error
}

type configurator struct {
	logger      logging.Logger
	client      network.Client
	fileManager files.FileManager
}

func NewConfigurator(logger logging.Logger, client network.Client, fileManager files.FileManager) Configurator {
	return &configurator{
		logger:      logger,
		client:      client,
		fileManager: fileManager,
	}
}
