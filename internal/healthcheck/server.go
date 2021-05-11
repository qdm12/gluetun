package healthcheck

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/qdm12/golibs/logging"
)

type Server interface {
	Run(ctx context.Context, healthy chan<- bool, done chan<- struct{})
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

func (s *server) Run(ctx context.Context, healthy chan<- bool, done chan<- struct{}) {
	defer close(done)

	loopDone := make(chan struct{})
	go s.runHealthcheckLoop(ctx, healthy, loopDone)

	server := http.Server{
		Addr:    s.address,
		Handler: s.handler,
	}
	serverDone := make(chan struct{})
	go func() {
		defer close(serverDone)
		<-ctx.Done()
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

	<-loopDone
	<-serverDone
}
