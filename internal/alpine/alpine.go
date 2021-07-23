// Package alpine defines a configurator to interact with the Alpine operating system.
package alpine

import (
	"os/user"
)

var _ Alpiner = (*Alpine)(nil)

type Alpiner interface {
	UserCreater
	VersionGetter
}

type Alpine struct {
	alpineReleasePath string
	passwdPath        string
	lookupID          func(uid string) (*user.User, error)
	lookup            func(username string) (*user.User, error)
}

func New() *Alpine {
	return &Alpine{
		alpineReleasePath: "/etc/alpine-release",
		passwdPath:        "/etc/passwd",
		lookupID:          user.LookupId,
		lookup:            user.Lookup,
	}
}
