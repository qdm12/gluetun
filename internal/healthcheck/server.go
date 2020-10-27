package healthcheck

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/golibs/logging"
)

type Server interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
}

type server struct {
	address string
	logger  logging.Logger
	handler http.Handler
}

func NewServer(address string, logger logging.Logger) Server {
	return &server{
		address: address,
		logger:  logger.WithPrefix("healthcheck: "),
		handler: newHandler(logger, &net.Resolver{}),
	}
}

func (s *server) Run(ctx context.Context, wg *sync.WaitGroup) {
	server := http.Server{
		Addr:    s.address,
		Handler: s.handler,
	}
	go func() {
		defer wg.Done()
		<-ctx.Done()
		s.logger.Warn("context canceled: shutting down server")
		defer s.logger.Warn("server shut down")
		const shutdownGraceDuration = 2 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGraceDuration)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("failed shutting down: %s", err)
		}
	}()
	s.logger.Info("listening on %s", s.address)
	err := server.ListenAndServe()
	if err != nil && !errors.Is(ctx.Err(), context.Canceled) {
		s.logger.Error(err)
	}
}
