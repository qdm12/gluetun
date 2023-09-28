package service

import (
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gosettings"
)

type Settings struct {
	Enabled       *bool
	PortForwarder PortForwarder
	Filepath      string
	Interface     string // needed for PIA and ProtonVPN, tun0 for example
	ServerName    string // needed for PIA
}

func (s Settings) Copy() (copied Settings) {
	copied.Enabled = gosettings.CopyPointer(s.Enabled)
	copied.PortForwarder = s.PortForwarder
	copied.Filepath = s.Filepath
	copied.Interface = s.Interface
	copied.ServerName = s.ServerName
	return copied
}

func (s *Settings) OverrideWith(update Settings) {
	s.Enabled = gosettings.OverrideWithPointer(s.Enabled, update.Enabled)
	s.PortForwarder = gosettings.OverrideWithInterface(s.PortForwarder, update.PortForwarder)
	s.Filepath = gosettings.OverrideWithString(s.Filepath, update.Filepath)
	s.Interface = gosettings.OverrideWithString(s.Interface, update.Interface)
	s.ServerName = gosettings.OverrideWithString(s.ServerName, update.ServerName)
}

var (
	ErrServerNameNotSet = errors.New("server name not set")
	ErrFilepathNotSet   = errors.New("file path not set")
	ErrInterfaceNotSet  = errors.New("interface not set")
)

func (s *Settings) Validate() (err error) {
	switch {
	// Port forwarder can be nil when the loop updates
	// to stop the service.
	case s.Filepath == "":
		return fmt.Errorf("%w", ErrFilepathNotSet)
	case s.Interface == "":
		return fmt.Errorf("%w", ErrInterfaceNotSet)
	case s.PortForwarder.Name() == providers.PrivateInternetAccess && s.ServerName == "":
		return fmt.Errorf("%w", ErrServerNameNotSet)
	}
	return nil
}
