package healthcheck

import (
	"context"
	"net"
	"net/http"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
)

type Server struct {
	logger  Logger
	handler *handler
	dialer  *net.Dialer
	config  settings.Health
	vpn     vpnHealth
	mux     *http.ServeMux
}

func NewServer(config settings.Health,
	logger Logger, vpnLoop StatusApplier,
) *Server {
	s := &Server{
		logger:  logger,
		handler: newHandler(),
		dialer: &net.Dialer{
			Resolver: &net.Resolver{
				PreferGo: true,
			},
		},
		config: config,
		vpn: vpnHealth{
			loop:        vpnLoop,
			healthyWait: *config.VPN.Initial,
		},
		mux: http.NewServeMux(),
	}
	s.mux.Handle("/", s.handler)
	s.mux.HandleFunc("/check/", func(w http.ResponseWriter, r *http.Request) {
		err := s.healthCheck(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	return s
}

type StatusApplier interface {
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
}
