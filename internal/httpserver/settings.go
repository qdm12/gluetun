package httpserver

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gotree"
	"github.com/qdm12/govalid/address"
)

type Settings struct {
	// Address is the server listening address.
	// It defaults to :8000.
	Address string
	// Handler is the HTTP Handler to use.
	// It must be set and cannot be left to nil.
	Handler http.Handler
	// Logger is the logger to use.
	// It must be set and cannot be left to nil.
	Logger Logger
	// ReadHeaderTimeout is the HTTP header read timeout duration
	// of the HTTP server. It defaults to 3 seconds if left unset.
	ReadHeaderTimeout time.Duration
	// ReadTimeout is the HTTP read timeout duration
	// of the HTTP server. It defaults to 3 seconds if left unset.
	ReadTimeout time.Duration
	// ShutdownTimeout is the shutdown timeout duration
	// of the HTTP server. It defaults to 3 seconds if left unset.
	ShutdownTimeout time.Duration
}

func (s *Settings) SetDefaults() {
	s.Address = gosettings.DefaultString(s.Address, ":8000")
	const defaultReadTimeout = 3 * time.Second
	s.ReadHeaderTimeout = gosettings.DefaultNumber(s.ReadHeaderTimeout, defaultReadTimeout)
	s.ReadTimeout = gosettings.DefaultNumber(s.ReadTimeout, defaultReadTimeout)
	const defaultShutdownTimeout = 3 * time.Second
	s.ShutdownTimeout = gosettings.DefaultNumber(s.ShutdownTimeout, defaultShutdownTimeout)
}

func (s Settings) Copy() Settings {
	return Settings{
		Address:           s.Address,
		Handler:           s.Handler,
		Logger:            s.Logger,
		ReadHeaderTimeout: s.ReadHeaderTimeout,
		ReadTimeout:       s.ReadTimeout,
		ShutdownTimeout:   s.ShutdownTimeout,
	}
}

func (s *Settings) MergeWith(other Settings) {
	s.Address = gosettings.MergeWithString(s.Address, other.Address)
	s.Handler = gosettings.MergeWithInterface(s.Handler, other.Handler)
	if s.Logger == nil {
		s.Logger = other.Logger
	}
	s.ReadHeaderTimeout = gosettings.MergeWithNumber(s.ReadHeaderTimeout, other.ReadHeaderTimeout)
	s.ReadTimeout = gosettings.MergeWithNumber(s.ReadTimeout, other.ReadTimeout)
	s.ShutdownTimeout = gosettings.MergeWithNumber(s.ShutdownTimeout, other.ShutdownTimeout)
}

func (s *Settings) OverrideWith(other Settings) {
	s.Address = gosettings.OverrideWithString(s.Address, other.Address)
	s.Handler = gosettings.OverrideWithInterface(s.Handler, other.Handler)
	if other.Logger != nil {
		s.Logger = other.Logger
	}
	s.ReadHeaderTimeout = gosettings.OverrideWithNumber(s.ReadHeaderTimeout, other.ReadHeaderTimeout)
	s.ReadTimeout = gosettings.OverrideWithNumber(s.ReadTimeout, other.ReadTimeout)
	s.ShutdownTimeout = gosettings.OverrideWithNumber(s.ShutdownTimeout, other.ShutdownTimeout)
}

var (
	ErrHandlerIsNotSet           = errors.New("HTTP handler cannot be left unset")
	ErrLoggerIsNotSet            = errors.New("logger cannot be left unset")
	ErrReadHeaderTimeoutTooSmall = errors.New("read header timeout is too small")
	ErrReadTimeoutTooSmall       = errors.New("read timeout is too small")
	ErrShutdownTimeoutTooSmall   = errors.New("shutdown timeout is too small")
)

func (s Settings) Validate() (err error) {
	uid := os.Getuid()
	err = address.Validate(s.Address, address.OptionListening(uid))
	if err != nil {
		return err
	}

	if s.Handler == nil {
		return fmt.Errorf("%w", ErrHandlerIsNotSet)
	}

	if s.Logger == nil {
		return fmt.Errorf("%w", ErrLoggerIsNotSet)
	}

	const minReadTimeout = time.Millisecond
	if s.ReadHeaderTimeout < minReadTimeout {
		return fmt.Errorf("%w: %s must be at least %s",
			ErrReadHeaderTimeoutTooSmall,
			s.ReadHeaderTimeout, minReadTimeout)
	}

	if s.ReadTimeout < minReadTimeout {
		return fmt.Errorf("%w: %s must be at least %s",
			ErrReadTimeoutTooSmall,
			s.ReadTimeout, minReadTimeout)
	}

	const minShutdownTimeout = 5 * time.Millisecond
	if s.ShutdownTimeout < minShutdownTimeout {
		return fmt.Errorf("%w: %s must be at least %s",
			ErrShutdownTimeoutTooSmall,
			s.ShutdownTimeout, minShutdownTimeout)
	}

	return nil
}

func (s Settings) ToLinesNode() (node *gotree.Node) {
	node = gotree.New("HTTP server settings:")
	node.Appendf("Listening address: %s", s.Address)
	node.Appendf("Read header timeout: %s", s.ReadHeaderTimeout)
	node.Appendf("Read timeout: %s", s.ReadTimeout)
	node.Appendf("Shutdown timeout: %s", s.ShutdownTimeout)
	return node
}

func (s Settings) String() string {
	return s.ToLinesNode().String()
}
