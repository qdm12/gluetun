package service

import (
	"context"
	"net/http"
	"sync"
)

type Service struct {
	// State
	portMutex sync.RWMutex
	ports     []uint16
	// Fixed parameters
	settings Settings
	puid     int
	pgid     int
	// Fixed injected objects
	routing     Routing
	client      *http.Client
	portAllower PortAllower
	logger      Logger
	// Internal channels and locks
	startStopMutex sync.Mutex
	keepPortCancel context.CancelFunc
	keepPortDoneCh <-chan struct{}
}

func New(settings Settings, routing Routing, client *http.Client,
	portAllower PortAllower, logger Logger, puid, pgid int) *Service {
	return &Service{
		// Fixed parameters
		settings: settings,
		puid:     puid,
		pgid:     pgid,
		// Fixed injected objects
		routing:     routing,
		client:      client,
		portAllower: portAllower,
		logger:      logger,
	}
}

func (s *Service) GetPortsForwarded() (ports []uint16) {
	s.portMutex.RLock()
	defer s.portMutex.RUnlock()
	ports = make([]uint16, len(s.ports))
	copy(ports, s.ports)
	return ports
}

func (s *Service) SetPortsForwarded(ctx context.Context, ports []uint16) (err error) {
	for i, port := range s.ports {
		err := s.portAllower.RemoveAllowedPort(ctx, port)
		if err != nil {
			for j := 0; j < i; j++ {
				_ = s.portAllower.SetAllowedPort(ctx, s.ports[j], s.settings.Interface)
			}
			s.logger.Error(err.Error())
			return err
		}
	}

	for i, port := range ports {
		err := s.portAllower.SetAllowedPort(ctx, port, s.settings.Interface)
		if err != nil {
			for j := 0; j < i; j++ {
				_ = s.portAllower.RemoveAllowedPort(ctx, s.ports[j])
			}
			for _, port := range s.ports {
				_ = s.portAllower.SetAllowedPort(ctx, port, s.settings.Interface)
			}
			s.logger.Error(err.Error())
			return err
		}
	}

	err = s.writePortForwardedFile(ports)
	if err != nil {
		_ = s.cleanup()
		return err
	}

	s.portMutex.RLock()
	defer s.portMutex.RUnlock()
	s.ports = make([]uint16, len(ports))
	copy(s.ports, ports)

	s.logger.Info("updated: " + portsToString(s.ports))

	return nil
}
