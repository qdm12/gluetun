package httpproxy

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/golibs/logging"
)

type Server interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
}

type server struct {
	address    string
	handler    http.Handler
	logger     logging.Logger
	internalWG *sync.WaitGroup
}

func New(ctx context.Context, address string,
	logger logging.Logger, client *http.Client,
	stealth, verbose bool) Server {
	proxyLogger := logger.WithPrefix("http proxy: ")
	wg := &sync.WaitGroup{}
	return &server{
		address:    address,
		handler:    newHandler(ctx, wg, client, proxyLogger, stealth, verbose),
		logger:     proxyLogger,
		internalWG: wg,
	}
}

func (s *server) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	server := http.Server{Addr: s.address, Handler: s.handler}
	go func() {
		<-ctx.Done()
		s.logger.Warn("context canceled: exiting loop")
		defer s.logger.Warn("loop exited")
		const shutdownGraceDuration = 2 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGraceDuration)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("failed shutting down: %s", err)
		}
	}()
	s.logger.Info("listening on %s", s.address)
	err := server.ListenAndServe()
	if err != nil && ctx.Err() != context.Canceled {
		s.logger.Error(err)
	}
	s.internalWG.Wait()
}
