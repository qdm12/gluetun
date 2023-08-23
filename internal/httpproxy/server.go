package httpproxy

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/qdm12/goservices"
	"github.com/qdm12/goservices/httpserver"
)

type Server struct {
	httpServer    *httpserver.Server
	handlerCtx    context.Context //nolint:containedctx
	handlerCancel context.CancelFunc
	handlerWg     *sync.WaitGroup

	// Server settings
	httpServerSettings httpserver.Settings

	// Handler settings
	logger   Logger
	stealth  bool
	verbose  bool
	username string
	password string
}

func ptrTo[T any](x T) *T { return &x }

func New(address string, logger Logger,
	stealth, verbose bool, username, password string,
	readHeaderTimeout, readTimeout time.Duration,
) (server *Server, err error) {
	return &Server{
		handlerWg: &sync.WaitGroup{},
		httpServerSettings: httpserver.Settings{
			// Handler is set when calling Start and reset when Stop is called
			Handler:           nil,
			Name:              ptrTo("proxy"),
			Address:           ptrTo(address),
			ReadTimeout:       readTimeout,
			ReadHeaderTimeout: readHeaderTimeout,
			Logger:            logger,
		},
		logger:   logger,
		stealth:  stealth,
		verbose:  verbose,
		username: username,
		password: password,
	}, nil
}

func (s *Server) Start(ctx context.Context) (
	runError <-chan error, err error,
) {
	if s.httpServer != nil {
		return nil, fmt.Errorf("%w", goservices.ErrAlreadyStarted)
	}

	s.handlerCtx, s.handlerCancel = context.WithCancel(context.Background())
	s.httpServerSettings.Handler = newHandler(s.handlerCtx, s.handlerWg,
		s.logger, s.stealth, s.verbose, s.username, s.password)
	s.httpServer, err = httpserver.New(s.httpServerSettings)
	if err != nil {
		return nil, fmt.Errorf("creating http server: %w", err)
	}

	return s.httpServer.Start(ctx)
}

func (s *Server) Stop() (err error) {
	if s.httpServer == nil {
		return fmt.Errorf("%w", goservices.ErrAlreadyStopped)
	}
	s.handlerCancel()
	err = s.httpServer.Stop()
	s.handlerWg.Wait()
	s.httpServer = nil // signal the server is down
	return err
}
