// Package server defines an interface to run the HTTP control server.
package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/dns"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/portforward"
	"github.com/qdm12/gluetun/internal/publicip"
	"github.com/qdm12/gluetun/internal/updater"
	"github.com/qdm12/gluetun/internal/vpn"
	"github.com/qdm12/golibs/logging"
)

type Server interface {
	Run(ctx context.Context, done chan<- struct{})
}

type server struct {
	address string
	logger  logging.Logger
	handler http.Handler
}

func New(ctx context.Context, address string, logEnabled bool, logger logging.Logger,
	buildInfo models.BuildInformation, openvpnLooper vpn.Looper,
	pfGetter portforward.Getter, unboundLooper dns.Looper,
	updaterLooper updater.Looper, publicIPLooper publicip.Looper) Server {
	handler := newHandler(ctx, logger, logEnabled, buildInfo,
		openvpnLooper, pfGetter, unboundLooper, updaterLooper, publicIPLooper)
	return &server{
		address: address,
		logger:  logger,
		handler: handler,
	}
}

func (s *server) Run(ctx context.Context, done chan<- struct{}) {
	defer close(done)
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
	if err != nil && errors.Is(ctx.Err(), context.Canceled) {
		s.logger.Error(err.Error())
	}
}
