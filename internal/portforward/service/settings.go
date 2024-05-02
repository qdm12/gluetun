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
	Interface      string // needed for PIA and ProtonVPN, tun0 for example
	ServerName     string // needed for PIA
	CanPortForward bool   // needed for PIA
	ListeningPort  uint16
}

func (s Settings) Copy() (copied Settings) {
	copied.Enabled = gosettings.CopyPointer(s.Enabled)
	copied.PortForwarder = s.PortForwarder
	copied.Filepath = s.Filepath
	copied.Interface = s.Interface
	copied.ServerName = s.ServerName
	copied.CanPortForward = s.CanPortForward
	copied.ListeningPort = s.ListeningPort
	return copied
}

func (s *Settings) OverrideWith(update Settings) {
	s.Enabled = gosettings.OverrideWithPointer(s.Enabled, update.Enabled)
	s.PortForwarder = gosettings.OverrideWithComparable(s.PortForwarder, update.PortForwarder)
	s.Filepath = gosettings.OverrideWithComparable(s.Filepath, update.Filepath)
	s.Interface = gosettings.OverrideWithComparable(s.Interface, update.Interface)
	s.ServerName = gosettings.OverrideWithComparable(s.ServerName, update.ServerName)
	s.CanPortForward = gosettings.OverrideWithComparable(s.CanPortForward, update.CanPortForward)
	s.ListeningPort = gosettings.OverrideWithComparable(s.ListeningPort, update.ListeningPort)
}

var (
	ErrPortForwarderNotSet = errors.New("port forwarder not set")
	ErrServerNameNotSet    = errors.New("server name not set")
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
	case s.PortForwarder.Name() == providers.PrivateInternetAccess && s.ServerName == "":
		return fmt.Errorf("%w", ErrServerNameNotSet)
	}
	return nil
}
