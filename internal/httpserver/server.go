package httpserver

import (
	"fmt"
	"net/http"
	"time"
)

// Server is an HTTP server implementation, which uses
// the HTTP handler provided.
type Server struct {
	address           string
	addressSet        chan struct{}
	handler           http.Handler
	logger            Logger
	readHeaderTimeout time.Duration
	readTimeout       time.Duration
	shutdownTimeout   time.Duration
}

// New creates a new HTTP server with the given settings.
// It returns an error if one of the settings is not valid.
func New(settings Settings) (s *Server, err error) {
	settings.SetDefaults()

	if err = settings.Validate(); err != nil {
		return nil, fmt.Errorf("http server settings validation failed: %w", err)
	}

	return &Server{
		address:           settings.Address,
		addressSet:        make(chan struct{}),
		handler:           settings.Handler,
		logger:            settings.Logger,
		readHeaderTimeout: settings.ReadHeaderTimeout,
		readTimeout:       settings.ReadTimeout,
		shutdownTimeout:   settings.ShutdownTimeout,
	}, nil
}
