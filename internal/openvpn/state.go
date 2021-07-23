package openvpn

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
)

func newState(status models.LoopStatus,
	settings configuration.OpenVPN, allServers models.AllServers,
	start chan<- struct{}, running <-chan models.LoopStatus,
	stop chan<- struct{}, stopped <-chan struct{}) *state {
	return &state{
		status:     status,
		settings:   settings,
		allServers: allServers,
		start:      start,
		running:    running,
		stop:       stop,
		stopped:    stopped,
	}
}

type state struct {
	loopMu sync.RWMutex

	status   models.LoopStatus
	statusMu sync.RWMutex

	settings   configuration.OpenVPN
	settingsMu sync.RWMutex

	allServers   models.AllServers
	allServersMu sync.RWMutex

	portForwarded   uint16
	portForwardedMu sync.RWMutex

	start   chan<- struct{}
	running <-chan models.LoopStatus
	stop    chan<- struct{}
	stopped <-chan struct{}
}

func (s *state) Lock()   { s.loopMu.Lock() }
func (s *state) Unlock() { s.loopMu.Unlock() }

// SetStatus sets the status thread safely.
// It should only be called by the loop internal code since
// it does not interact with the loop code directly.
func (s *state) SetStatus(status models.LoopStatus) {
	s.statusMu.Lock()
	defer s.statusMu.Unlock()
	s.status = status
}

// GetStatus gets the status thread safely.
func (s *state) GetStatus() (status models.LoopStatus) {
	s.statusMu.RLock()
	defer s.statusMu.RUnlock()
	return s.status
}

func (s *state) GetSettingsAndServers() (settings configuration.OpenVPN,
	allServers models.AllServers) {
	s.settingsMu.RLock()
	s.allServersMu.RLock()
	settings = s.settings
	allServers = s.allServers
	s.settingsMu.RUnlock()
	s.allServersMu.RUnlock()
	return settings, allServers
}

var ErrInvalidStatus = errors.New("invalid status")

// ApplyStatus sends signals to the running loop depending on the
// current status and status requested, such that its next status
// matches the requested one. It is thread safe and a synchronous call
// since it waits to the loop to fully change its status.
func (s *state) ApplyStatus(ctx context.Context, status models.LoopStatus) (
	outcome string, err error) {
	// prevent simultaneous loop changes by restricting
	// multiple SetStatus calls to run sequentially.
	s.loopMu.Lock()
	defer s.loopMu.Unlock()

	// not a read lock as we want to modify it eventually in
	// the code below before any other call.
	s.statusMu.Lock()
	existingStatus := s.status

	switch status {
	case constants.Running:
		if existingStatus != constants.Stopped {
			return "already " + existingStatus.String(), nil
		}

		s.status = constants.Starting
		s.statusMu.Unlock()
		s.start <- struct{}{}

		// Wait for the loop to react to the start signal
		newStatus := constants.Starting // for canceled context
		select {
		case <-ctx.Done():
		case newStatus = <-s.running:
		}
		s.SetStatus(newStatus)

		return newStatus.String(), nil
	case constants.Stopped:
		if existingStatus != constants.Running {
			return "already " + existingStatus.String(), nil
		}

		s.status = constants.Stopping
		s.statusMu.Unlock()
		s.stop <- struct{}{}

		// Wait for the loop to react to the stop signal
		newStatus := constants.Stopping // for canceled context
		select {
		case <-ctx.Done():
		case <-s.stopped:
			newStatus = constants.Stopped
		}
		s.SetStatus(newStatus)

		return newStatus.String(), nil
	default:
		return "", fmt.Errorf("%w: %s: it can only be one of: %s, %s",
			ErrInvalidStatus, status, constants.Running, constants.Stopped)
	}
}

func (s *state) GetSettings() (settings configuration.OpenVPN) {
	s.settingsMu.RLock()
	defer s.settingsMu.RUnlock()
	return s.settings
}

func (s *state) SetSettings(ctx context.Context, settings configuration.OpenVPN) (
	outcome string) {
	s.settingsMu.Lock()
	defer s.settingsMu.Unlock()
	settingsUnchanged := reflect.DeepEqual(s.settings, settings)
	if settingsUnchanged {
		return "settings left unchanged"
	}
	s.settings = settings
	_, _ = s.ApplyStatus(ctx, constants.Stopped)
	outcome, _ = s.ApplyStatus(ctx, constants.Running)
	return outcome
}

func (s *state) GetServers() (servers models.AllServers) {
	s.allServersMu.RLock()
	defer s.allServersMu.RUnlock()
	return s.allServers
}

func (s *state) SetServers(servers models.AllServers) {
	s.allServersMu.Lock()
	defer s.allServersMu.Unlock()
	s.allServers = servers
}

func (s *state) GetPortForwarded() (port uint16) {
	s.portForwardedMu.RLock()
	defer s.portForwardedMu.RUnlock()
	return s.portForwarded
}
