// Package alpine defines a configurator to interact with the Alpine operating system.
package alpine

import (
	"context"
	"os/user"
)

type Configurator interface {
	CreateUser(username string, uid int) (createdUsername string, err error)
	Version(ctx context.Context) (version string, err error)
}

type configurator struct {
	alpineReleasePath string
	passwdPath        string
	lookupID          func(uid string) (*user.User, error)
	lookup            func(username string) (*user.User, error)
}

func NewConfigurator() Configurator {
	return &configurator{
		alpineReleasePath: "/etc/alpine-release",
		passwdPath:        "/etc/passwd",
		lookupID:          user.LookupId,
		lookup:            user.Lookup,
	}
}
