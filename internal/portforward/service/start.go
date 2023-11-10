package service

import (
	"context"
	"fmt"

	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (s *Service) Start(ctx context.Context) (runError <-chan error, err error) {
	s.startStopMutex.Lock()
	defer s.startStopMutex.Unlock()

	if !*s.settings.Enabled {
		return nil, nil //nolint:nilnil
	}

	s.logger.Info("starting")

	gateway, err := s.routing.VPNLocalGatewayIP(s.settings.Interface)
	if err != nil {
		return nil, fmt.Errorf("getting VPN local gateway IP: %w", err)
	}

	obj := utils.PortForwardObjects{
		Logger:     s.logger,
		Gateway:    gateway,
		Client:     s.client,
		ServerName: s.settings.ServerName,
	}
	port, err := s.settings.PortForwarder.PortForward(ctx, obj)
	if err != nil {
		return nil, fmt.Errorf("port forwarding for the first time: %w", err)
	}

	s.logger.Info("port forwarded is " + fmt.Sprint(int(port)))

	err = s.portAllower.SetAllowedPort(ctx, port, s.settings.Interface)
	if err != nil {
		return nil, fmt.Errorf("allowing port in firewall: %w", err)
	}

	if s.settings.ListeningPort != 0 {
		err = s.portAllower.RedirectPort(ctx, s.settings.Interface, port, s.settings.ListeningPort)
		if err != nil {
			return nil, fmt.Errorf("redirecting port in firewall: %w", err)
		}
	}

	err = s.writePortForwardedFile(port)
	if err != nil {
		_ = s.cleanup()
		return nil, fmt.Errorf("writing port file: %w", err)
	}

	s.portMutex.Lock()
	s.port = port
	s.portMutex.Unlock()

	keepPortCtx, keepPortCancel := context.WithCancel(context.Background())
	s.keepPortCancel = keepPortCancel
	runErrorCh := make(chan error)
	keepPortDoneCh := make(chan struct{})
	s.keepPortDoneCh = keepPortDoneCh

	go func(ctx context.Context, portForwarder PortForwarder,
		obj utils.PortForwardObjects, runError chan<- error, doneCh chan<- struct{}) {
		defer close(doneCh)
		err = portForwarder.KeepPortForward(ctx, obj)
		crashed := ctx.Err() == nil
		if !crashed { // stopped by Stop call
			return
		}
		s.startStopMutex.Lock()
		defer s.startStopMutex.Unlock()
		_ = s.cleanup()
		runError <- err
	}(keepPortCtx, s.settings.PortForwarder, obj, runErrorCh, keepPortDoneCh)

	return runErrorCh, nil
}
