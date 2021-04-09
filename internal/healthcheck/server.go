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
	Run(ctx context.Context, healthy chan<- bool, wg *sync.WaitGroup)
}

type server struct {
	address  string
	logger   logging.Logger
	handler  *handler
	resolver *net.Resolver
}

func NewServer(address string, logger logging.Logger) Server {
	healthcheckLogger := logger.NewChild(logging.SetPrefix("healthcheck: "))
	return &server{
		address:  address,
		logger:   healthcheckLogger,
		handler:  newHandler(healthcheckLogger),
		resolver: net.DefaultResolver,
	}
}

func (s *server) Run(ctx context.Context, healthy chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()

	internalWg := &sync.WaitGroup{}
	internalWg.Add(1)
	go s.runHealthcheckLoop(ctx, healthy, internalWg)

	server := http.Server{
		Addr:    s.address,
		Handler: s.handler,
	}
	internalWg.Add(1)
	go func() {
		defer internalWg.Done()
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

	internalWg.Wait()
}
