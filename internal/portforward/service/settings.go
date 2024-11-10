package service

import (
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gosettings"
)

type Settings struct {
	Enabled        *bool
	PortForwarder  PortForwarder
	Filepath       string
	UpCommand      string
	DownCommand    string
	Interface      string // needed for PIA, PrivateVPN and ProtonVPN, tun0 for example
	ServerName     string // needed for PIA
	CanPortForward bool   // needed for PIA
	ListeningPort  uint16
	Username       string // needed for PIA
	Password       string // needed for PIA
}

func (s Settings) Copy() (copied Settings) {
	copied.Enabled = gosettings.CopyPointer(s.Enabled)
	copied.PortForwarder = s.PortForwarder
	copied.Filepath = s.Filepath
	copied.UpCommand = s.UpCommand
	copied.DownCommand = s.DownCommand
	copied.Interface = s.Interface
	copied.ServerName = s.ServerName
	copied.CanPortForward = s.CanPortForward
	copied.ListeningPort = s.ListeningPort
	copied.Username = s.Username
	copied.Password = s.Password
	return copied
}

func (s *Settings) OverrideWith(update Settings) {
	s.Enabled = gosettings.OverrideWithPointer(s.Enabled, update.Enabled)
	s.PortForwarder = gosettings.OverrideWithComparable(s.PortForwarder, update.PortForwarder)
	s.Filepath = gosettings.OverrideWithComparable(s.Filepath, update.Filepath)
	s.UpCommand = gosettings.OverrideWithComparable(s.UpCommand, update.UpCommand)
	s.DownCommand = gosettings.OverrideWithComparable(s.DownCommand, update.DownCommand)
	s.Interface = gosettings.OverrideWithComparable(s.Interface, update.Interface)
	s.ServerName = gosettings.OverrideWithComparable(s.ServerName, update.ServerName)
	s.CanPortForward = gosettings.OverrideWithComparable(s.CanPortForward, update.CanPortForward)
	s.ListeningPort = gosettings.OverrideWithComparable(s.ListeningPort, update.ListeningPort)
	s.Username = gosettings.OverrideWithComparable(s.Username, update.Username)
	s.Password = gosettings.OverrideWithComparable(s.Password, update.Password)
}

var (
	ErrPortForwarderNotSet = errors.New("port forwarder not set")
	ErrServerNameNotSet    = errors.New("server name not set")
	ErrUsernameNotSet      = errors.New("username not set")
	ErrPasswordNotSet      = errors.New("password not set")
	ErrFilepathNotSet      = errors.New("file path not set")
	ErrInterfaceNotSet     = errors.New("interface not set")
)

func (s *Settings) Validate(forStartup bool) (err error) {
	// Minimal validation
	if s.Filepath == "" {
		return fmt.Errorf("%w", ErrFilepathNotSet)
	}

	if !forStartup {
		// No additional validation needed if the service
		// is not to be started with the given settings.
		return nil
	}

	// Startup validation requires additional fields set.
	switch {
	case s.PortForwarder == nil:
		return fmt.Errorf("%w", ErrPortForwarderNotSet)
	case s.Interface == "":
		return fmt.Errorf("%w", ErrInterfaceNotSet)
	case s.PortForwarder.Name() == providers.PrivateInternetAccess:
		switch {
		case s.ServerName == "":
			return fmt.Errorf("%w", ErrServerNameNotSet)
		case s.Username == "":
			return fmt.Errorf("%w", ErrUsernameNotSet)
		case s.Password == "":
			return fmt.Errorf("%w", ErrPasswordNotSet)
		}
	}
	return nil
}
