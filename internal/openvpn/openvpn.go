package openvpn

import (
	"context"
	"io"
	"os"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
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
	fileManager files.FileManager
	logger      logging.Logger
	commander   command.Commander
	openFile    func(name string, flag int, perm os.FileMode) (*os.File, error)
	mkDev       func(major uint32, minor uint32) uint64
	mkNod       func(path string, mode uint32, dev int) error
}

func NewConfigurator(logger logging.Logger, fileManager files.FileManager) Configurator {
	return &configurator{
		fileManager: fileManager,
		logger:      logger.WithPrefix("openvpn configurator: "),
		commander:   command.NewCommander(),
		openFile:    os.OpenFile,
		mkDev:       unix.Mkdev,
		mkNod:       unix.Mknod,
	}
}
