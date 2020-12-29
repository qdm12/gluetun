package dns

import (
	"context"
	"io"
	"net"
	"net/http"

	"github.com/qdm12/gluetun/internal/os"
	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
)

type Configurator interface {
	DownloadRootHints(ctx context.Context, uid, gid int) error
	DownloadRootKey(ctx context.Context, uid, gid int) error
	MakeUnboundConf(ctx context.Context, settings settings.DNS, username string, uid, gid int) (err error)
	UseDNSInternally(IP net.IP)
	UseDNSSystemWide(ip net.IP, keepNameserver bool) error
	Start(ctx context.Context, logLevel uint8) (stdout io.ReadCloser, waitFn func() error, err error)
	WaitForUnbound() (err error)
	Version(ctx context.Context) (version string, err error)
}

type configurator struct {
	logger    logging.Logger
	client    *http.Client
	openFile  os.OpenFileFunc
	commander command.Commander
	lookupIP  func(host string) ([]net.IP, error)
}

func NewConfigurator(logger logging.Logger, httpClient *http.Client,
	openFile os.OpenFileFunc) Configurator {
	return &configurator{
		logger:    logger.WithPrefix("dns configurator: "),
		client:    httpClient,
		openFile:  openFile,
		commander: command.NewCommander(),
		lookupIP:  net.LookupIP,
	}
}
