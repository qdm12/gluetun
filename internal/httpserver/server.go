// Package httpserver implements an HTTP server.
package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var _ Interface = (*Server)(nil)

// Interface is the HTTP server composite interface.
type Interface interface {
	Runner
	AddressGetter
}

// Runner is the interface for an HTTP server with a Run method.
type Runner interface {
	Run(ctx context.Context, ready chan<- struct{}, done chan<- struct{})
}

// AddressGetter obtains the address the HTTP server is listening on.
type AddressGetter interface {
	GetAddress() (address string)
}

// Server is an HTTP server implementation, which uses
// the HTTP handler provided.
type Server struct {
	name            string
	address         string
	addressSet      chan struct{}
	handler         http.Handler
	logger          Logger
	shutdownTimeout time.Duration
}

// New creates a new HTTP server with the given settings.
// It returns an error if one of the settings is not valid.
func New(settings Settings) (s *Server, err error) {
	settings.SetDefaults()

	if err = settings.Validate(); err != nil {
		return nil, fmt.Errorf("http server settings validation failed: %w", err)
	}

	return &Server{
		name:            *settings.Name,
		address:         settings.Address,
		addressSet:      make(chan struct{}),
		handler:         settings.Handler,
		logger:          settings.Logger,
		shutdownTimeout: *settings.ShutdownTimeout,
	}, nil
}
