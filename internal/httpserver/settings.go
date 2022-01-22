package httpserver

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gotree"
	"github.com/qdm12/govalid/address"
)

type Settings struct {
	// Name is the server name to use in logs.
	// It defaults to the empty string.
	Name *string
	// Address is the server listening address.
	// It defaults to :8000.
	Address string
	// Handler is the HTTP Handler to use.
	// It must be set and cannot be left to nil.
	Handler http.Handler
	// Logger is the logger to use.
	// It must be set and cannot be left to nil.
	Logger Logger
	// ShutdownTimeout is the shutdown timeout duration
	// of the HTTP server. It defaults to 3 seconds.
	ShutdownTimeout *time.Duration
}

func (s *Settings) SetDefaults() {
	s.Name = helpers.DefaultStringPtr(s.Name, "")
	s.Address = helpers.DefaultString(s.Address, ":8000")
	const defaultShutdownTimeout = 3 * time.Second
	s.ShutdownTimeout = helpers.DefaultDuration(s.ShutdownTimeout, defaultShutdownTimeout)
}

func (s Settings) Copy() Settings {
	return Settings{
		Name:            helpers.CopyStringPtr(s.Name),
		Address:         s.Address,
		Handler:         s.Handler,
		Logger:          s.Logger,
		ShutdownTimeout: helpers.CopyDurationPtr(s.ShutdownTimeout),
	}
}

func (s *Settings) MergeWith(other Settings) {
	s.Name = helpers.MergeWithStringPtr(s.Name, other.Name)
	s.Address = helpers.MergeWithString(s.Address, other.Address)
	s.Handler = helpers.MergeWithHTTPHandler(s.Handler, other.Handler)
	if s.Logger == nil {
		s.Logger = other.Logger
	}
	s.ShutdownTimeout = helpers.MergeWithDuration(s.ShutdownTimeout, other.ShutdownTimeout)
}

func (s *Settings) OverrideWith(other Settings) {
	s.Name = helpers.OverrideWithStringPtr(s.Name, other.Name)
	s.Address = helpers.OverrideWithString(s.Address, other.Address)
	s.Handler = helpers.OverrideWithHTTPHandler(s.Handler, other.Handler)
	if other.Logger != nil {
		s.Logger = other.Logger
	}
	s.ShutdownTimeout = helpers.OverrideWithDuration(s.ShutdownTimeout, other.ShutdownTimeout)
}

var (
	ErrHandlerIsNotSet         = errors.New("HTTP handler cannot be left unset")
	ErrLoggerIsNotSet          = errors.New("logger cannot be left unset")
	ErrShutdownTimeoutTooSmall = errors.New("shutdown timeout is too small")
)

func (s Settings) Validate() (err error) {
	uid := os.Getuid()
	_, err = address.Validate(s.Address, address.OptionListening(uid))
	if err != nil {
		return err
	}

	if s.Handler == nil {
		return ErrHandlerIsNotSet
	}

	if s.Logger == nil {
		return ErrLoggerIsNotSet
	}

	const minShutdownTimeout = 5 * time.Millisecond
	if *s.ShutdownTimeout < minShutdownTimeout {
		return fmt.Errorf("%w: %s must be at least %s",
			ErrShutdownTimeoutTooSmall,
			*s.ShutdownTimeout, minShutdownTimeout)
	}

	return nil
}

func (s Settings) ToLinesNode() (node *gotree.Node) {
	node = gotree.New("%s HTTP server settings:", strings.Title(*s.Name))
	node.Appendf("Listening address: %s", s.Address)
	node.Appendf("Shutdown timeout: %s", *s.ShutdownTimeout)
	return node
}

func (s Settings) String() string {
	return s.ToLinesNode().String()
}
