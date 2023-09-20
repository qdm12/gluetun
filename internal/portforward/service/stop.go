package service

import (
	"context"
	"fmt"
	"os"
)

func (s *Service) Stop() (err error) {
	s.startStopMutex.Lock()
	defer s.startStopMutex.Unlock()

	if s.port == 0 {
		return nil // already not running
	}

	s.logger.Info("stopping")

	s.keepPortCancel()
	<-s.keepPortDoneCh

	return s.cleanup()
}

func (s *Service) cleanup() (err error) {
	err = s.portAllower.RemoveAllowedPort(context.Background(), s.port)
	if err != nil {
		return fmt.Errorf("blocking previous port in firewall: %w", err)
	}

	s.port = 0

	filepath := *s.settings.Settings.Filepath
	s.logger.Info("removing port file " + filepath)
	err = os.Remove(filepath)
	if err != nil {
		return fmt.Errorf("removing port file: %w", err)
	}

	return nil
}
