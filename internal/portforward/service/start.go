package service

import (
	"context"
	"fmt"
)

func (s *Service) Start(ctx context.Context) (runError <-chan error, err error) {
	s.startStopMutex.Lock()
	defer s.startStopMutex.Unlock()

	if !*s.settings.Settings.Enabled {
		return nil, nil //nolint:nilnil
	}

	s.logger.Info("starting")
	port, err := s.settings.PortForwarder.PortForward(ctx, s.client, s.logger,
		s.settings.Gateway, s.settings.ServerName)
	if err != nil {
		return nil, fmt.Errorf("port forwarding for the first time: %w", err)
	}

	s.logger.Info("port forwarded is " + fmt.Sprint(int(port)))

	err = s.portAllower.SetAllowedPort(ctx, port, s.settings.Interface)
	if err != nil {
		return nil, fmt.Errorf("allowing port in firewall: %w", err)
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

	go func(ctx context.Context, settings Settings, port uint16,
		runError chan<- error, doneCh chan<- struct{}) {
		defer close(doneCh)
		err = settings.PortForwarder.KeepPortForward(ctx, port,
			settings.Gateway, settings.ServerName, s.logger)
		crashed := ctx.Err() == nil
		if !crashed { // stopped by Stop call
			return
		}
		_ = s.cleanup()
		runError <- err
	}(keepPortCtx, s.settings, port, runErrorCh, keepPortDoneCh)

	return runErrorCh, nil
}
