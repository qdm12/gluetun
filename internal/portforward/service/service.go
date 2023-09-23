package service

import (
	"context"
	"net/http"
	"sync"
)

type Service struct {
	// State
	portMutex sync.RWMutex
	port      uint16
	// Fixed parameters
	settings Settings
	puid     int
	pgid     int
	// Fixed injected objets
	client      *http.Client
	portAllower PortAllower
	logger      Logger
	// Internal channels and locks
	startStopMutex sync.Mutex
	keepPortCancel context.CancelFunc
	keepPortDoneCh <-chan struct{}
}

func New(settings Settings, client *http.Client,
	portAllower PortAllower, logger Logger, puid, pgid int) *Service {
	return &Service{
		// Fixed parameters
		settings: settings,
		puid:     puid,
		pgid:     pgid,
		// Fixed injected objets
		client:      client,
		portAllower: portAllower,
		logger:      logger,
	}
}

func (s *Service) GetPortForwarded() (port uint16) {
	s.portMutex.RLock()
	defer s.portMutex.RUnlock()
	return s.port
}
