package alpine

import (
	"os/user"

	"github.com/qdm12/golibs/files"
)

type Configurator interface {
	CreateUser(username string, uid int) error
}

type configurator struct {
	fileManager files.FileManager
	lookupUID   func(uid string) (*user.User, error)
	lookupUser  func(username string) (*user.User, error)
}

func NewConfigurator(fileManager files.FileManager) Configurator {
	return &configurator{
		fileManager: fileManager,
		lookupUID:   user.LookupId,
		lookupUser:  user.Lookup,
	}
}
