package dns

import (
	"context"
	"io"
	"net"

	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
)

type Configurator interface {
	DownloadRootHints(uid, gid int) error
	DownloadRootKey(uid, gid int) error
	MakeUnboundConf(settings settings.DNS, uid, gid int) (err error)
	UseDNSInternally(IP net.IP)
	UseDNSSystemWide(ip net.IP, keepNameserver bool) error
	Start(ctx context.Context, logLevel uint8) (stdout io.ReadCloser, waitFn func() error, err error)
	WaitForUnbound() (err error)
	Version(ctx context.Context) (version string, err error)
}

type configurator struct {
	logger      logging.Logger
	client      network.Client
	fileManager files.FileManager
	commander   command.Commander
	lookupIP    func(host string) ([]net.IP, error)
}

func NewConfigurator(logger logging.Logger, client network.Client, fileManager files.FileManager) Configurator {
	return &configurator{
		logger:      logger.WithPrefix("dns configurator: "),
		client:      client,
		fileManager: fileManager,
		commander:   command.NewCommander(),
		lookupIP:    net.LookupIP,
	}
}
