package alpine

import (
	"os/user"

	"github.com/qdm12/gluetun/internal/os"
)

type Configurator interface {
	CreateUser(username string, uid int) (createdUsername string, err error)
}

type configurator struct {
	openFile   os.OpenFileFunc
	lookupUID  func(uid string) (*user.User, error)
	lookupUser func(username string) (*user.User, error)
}

func NewConfigurator(openFile os.OpenFileFunc) Configurator {
	return &configurator{
		openFile:   openFile,
		lookupUID:  user.LookupId,
		lookupUser: user.Lookup,
	}
}
