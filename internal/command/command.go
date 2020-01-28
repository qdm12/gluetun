package command

import (
	libcommand "github.com/qdm12/golibs/command"
)

type Command interface {
	VersionOpenVPN() (string, error)
	VersionUnbound() (string, error)
	VersionIptables() (string, error)
	VersionShadowSocks() (string, error)
	VersionTinyProxy() (string, error)
	Unbound() error
}

type command struct {
	command libcommand.Command
}

func NewCommand() Command {
	return &command{
		command: libcommand.NewCommand(),
	}
}
