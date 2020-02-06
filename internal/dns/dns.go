package dns

import (
	"io"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/settings"
)

const logPrefix = "dns configurator"

type Configurator interface {
	DownloadRootHints(uid, gid int) error
	DownloadRootKey(uid, gid int) error
	MakeUnboundConf(settings settings.DNS, uid, gid int) (err error)
	SetLocalNameserver() error
	Start() (stdout io.ReadCloser, err error)
	Version() (version string, err error)
}

type configurator struct {
	logger      logging.Logger
	client      network.Client
	fileManager files.FileManager
	commander   command.Commander
}

func NewConfigurator(logger logging.Logger, client network.Client, fileManager files.FileManager) Configurator {
	return &configurator{
		logger:      logger,
		client:      client,
		fileManager: fileManager,
		commander:   command.NewCommander(),
	}
}
