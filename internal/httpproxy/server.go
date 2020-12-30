package httpproxy

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/golibs/logging"
)

type Server interface {
	Run(ctx context.Context, wg *sync.WaitGroup, errorCh chan<- error)
}

type server struct {
	address    string
	handler    http.Handler
	logger     logging.Logger
	internalWG *sync.WaitGroup
}

func New(ctx context.Context, address string, logger logging.Logger,
	stealth, verbose bool, username, password string) Server {
	wg := &sync.WaitGroup{}
	return &server{
		address:    address,
		handler:    newHandler(ctx, wg, logger, stealth, verbose, username, password),
		logger:     logger,
		internalWG: wg,
	}
}

func (s *server) Run(ctx context.Context, wg *sync.WaitGroup, errorCh chan<- error) {
	defer wg.Done()
	server := http.Server{Addr: s.address, Handler: s.handler}
	go func() {
		<-ctx.Done()
		s.logger.Warn("shutting down server")
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
	if err != nil && ctx.Err() == nil {
		errorCh <- err
	}
	s.internalWG.Wait()
}
