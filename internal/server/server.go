package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/dns"
	"github.com/qdm12/gluetun/internal/openvpn"
	"github.com/qdm12/gluetun/internal/updater"
	"github.com/qdm12/golibs/logging"
)

type Server interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
}

type server struct {
	address       string
	logging       bool
	logger        logging.Logger
	openvpnLooper openvpn.Looper
	unboundLooper dns.Looper
	updaterLooper updater.Looper
	lookupIP      func(host string) ([]net.IP, error)
}

func New(address string, logging bool, logger logging.Logger,
	openvpnLooper openvpn.Looper, unboundLooper dns.Looper, updaterLooper updater.Looper) Server {
	return &server{
		address:       address,
		logging:       logging,
		logger:        logger.WithPrefix("http server: "),
		openvpnLooper: openvpnLooper,
		unboundLooper: unboundLooper,
		updaterLooper: updaterLooper,
		lookupIP:      net.LookupIP,
	}
}

func (s *server) Run(ctx context.Context, wg *sync.WaitGroup) {
	server := http.Server{Addr: s.address, Handler: s.makeHandler()}
	go func() {
		defer wg.Done()
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
}

func (s *server) makeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("HTTP %s %s", r.Method, r.RequestURI)
		switch r.Method {
		case http.MethodGet:
			switch r.RequestURI {
			case "/openvpn/actions/restart":
				s.openvpnLooper.Restart()
				w.WriteHeader(http.StatusOK)
			case "/unbound/actions/restart":
				s.unboundLooper.Restart()
				w.WriteHeader(http.StatusOK)
			case "/openvpn/portforwarded":
				s.handleGetPortForwarded(w)
			case "/openvpn/settings":
				s.handleGetOpenvpnSettings(w)
			case "/updater/restart":
				s.updaterLooper.Restart()
				w.WriteHeader(http.StatusOK)
			default:
				routeDoesNotExist(s.logger, w, r)
			}
		default:
			routeDoesNotExist(s.logger, w, r)
		}
	}
}

func routeDoesNotExist(logger logging.Logger, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte(fmt.Sprintf("Nothing here for %s %s", r.Method, r.RequestURI)))
	if err != nil {
		logger.Error(err)
	}
}
