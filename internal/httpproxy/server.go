package httpproxy

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/golibs/logging"
)

type Server interface {
	Run(ctx context.Context, errorCh chan<- error)
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

func (s *server) Run(ctx context.Context, errorCh chan<- error) {
	server := http.Server{Addr: s.address, Handler: s.handler}
	go func() {
		<-ctx.Done()
		const shutdownGraceDuration = 2 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGraceDuration)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("failed shutting down: " + err.Error())
		}
	}()
	s.logger.Info("listening on " + s.address)
	err := server.ListenAndServe()
	s.internalWG.Wait()
	if err != nil && ctx.Err() == nil {
		errorCh <- err
	} else {
		errorCh <- nil
	}
}
