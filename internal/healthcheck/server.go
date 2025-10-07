package healthcheck

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
)

type Server struct {
	logger  Logger
	handler *handler
	config  settings.Health
}

func NewServer(config settings.Health, logger Logger) *Server {
	return &Server{
		logger:  logger,
		handler: newHandler(logger),
		config:  config,
	}
}

func (s *Server) SetError(err error) {
	s.handler.setErr(err)
}

type StatusApplier interface {
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
}
