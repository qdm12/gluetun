package openvpn

import (
	"context"
	"io"

	"github.com/qdm12/gluetun/internal/os"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
	"golang.org/x/sys/unix"
)

type Configurator interface {
	Version(ctx context.Context) (string, error)
	WriteAuthFile(user, password string, uid, gid int) error
	CheckTUN() error
	CreateTUN() error
	Start(ctx context.Context) (stdout io.ReadCloser, waitFn func() error, err error)
}

type configurator struct {
	logger    logging.Logger
	commander command.Commander
	os        os.OS
	mkDev     func(major uint32, minor uint32) uint64
	mkNod     func(path string, mode uint32, dev int) error
}

func NewConfigurator(logger logging.Logger, os os.OS) Configurator {
	return &configurator{
		logger:    logger.WithPrefix("openvpn configurator: "),
		commander: command.NewCommander(),
		os:        os,
		mkDev:     unix.Mkdev,
		mkNod:     unix.Mknod,
	}
}
