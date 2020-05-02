package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/golibs/logging"
)

type Server interface {
	SetOpenVPNRestart(f func())
	Run(ctx context.Context) error
}

type server struct {
	address                 string
	logger                  logging.Logger
	restartOpenvpn          func()
	restartOpenvpnSet       context.Context
	restartOpenvpnSetSignal func()
	sync.RWMutex
}

func New(address string, logger logging.Logger) Server {
	restartOpenvpnSet, restartOpenvpnSetSignal := context.WithCancel(context.Background())
	return &server{
		address:                 address,
		logger:                  logger.WithPrefix("http server: "),
		restartOpenvpnSet:       restartOpenvpnSet,
		restartOpenvpnSetSignal: restartOpenvpnSetSignal,
	}
}

func (s *server) Run(ctx context.Context) error {
	if s.restartOpenvpnSet.Err() == nil {
		s.logger.Warn("restartOpenvpn function is not set, waiting...")
		<-s.restartOpenvpnSet.Done()
	}
	server := http.Server{Addr: s.address, Handler: s.makeHandler()}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("failed shutting down: %s", err)
		}
	}()
	s.logger.Info("listening on %s", s.address)
	return server.ListenAndServe()
}

func (s *server) SetOpenVPNRestart(f func()) {
	s.Lock()
	defer s.Unlock()
	s.restartOpenvpn = f
	if s.restartOpenvpnSet.Err() == nil {
		s.restartOpenvpnSetSignal()
	}
}

func (s *server) makeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("HTTP %s %s", r.Method, r.RequestURI)
		switch r.Method {
		case http.MethodGet:
			switch r.RequestURI {
			case "/openvpn/actions/restart":
				s.RLock()
				defer s.RUnlock()
				if s.restartOpenvpn == nil {
					functionNotSet("restartOpenvpn", s.logger, w)
					return
				}
				s.restartOpenvpn()
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

func functionNotSet(functionName string, logger logging.Logger, w http.ResponseWriter) {
	logger.Error("function %s is not set", functionName)
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte(fmt.Sprintf("%s function is not set", functionName)))
	if err != nil {
		logger.Error(err)
	}
}
