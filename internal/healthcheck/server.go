package healthcheck

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/openvpn"
	"github.com/qdm12/golibs/logging"
)

type Server interface {
	Run(ctx context.Context, done chan<- struct{})
}

type server struct {
	address  string
	logger   logging.Logger
	handler  *handler
	resolver *net.Resolver
	openvpn  openvpnHealth
}

type openvpnHealth struct {
	looper          openvpn.Looper
	healthyWaitTime time.Duration
	healthyTimer    *time.Timer
}

const (
	defaultOpenvpnHealthyWaitTime = 6 * time.Second
	openvpnHealthyWaitTimeAdd     = 5 * time.Second
)

func NewServer(address string, logger logging.Logger,
	openvpnLooper openvpn.Looper) Server {
	return &server{
		address:  address,
		logger:   logger,
		handler:  newHandler(logger),
		resolver: net.DefaultResolver,
		openvpn: openvpnHealth{
			looper:          openvpnLooper,
			healthyWaitTime: defaultOpenvpnHealthyWaitTime,
		},
	}
}

func (s *server) Run(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	loopDone := make(chan struct{})
	go s.runHealthcheckLoop(ctx, loopDone)

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
