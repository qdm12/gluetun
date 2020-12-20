package openvpn

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/settings"
)

type state struct {
	status          models.LoopStatus
	settings        settings.OpenVPN
	allServers      models.AllServers
	portForwarded   uint16
	statusMu        sync.RWMutex
	settingsMu      sync.RWMutex
	allServersMu    sync.RWMutex
	portForwardedMu sync.RWMutex
}

func (s *state) setStatusWithLock(status models.LoopStatus) {
	s.statusMu.Lock()
	defer s.statusMu.Unlock()
	s.status = status
}

func (s *state) getSettingsAndServers() (settings settings.OpenVPN, allServers models.AllServers) {
	s.settingsMu.RLock()
	s.allServersMu.RLock()
	settings = s.settings
	allServers = s.allServers
	s.settingsMu.RLock()
	s.allServersMu.RLock()
	return settings, allServers
}

func (l *looper) GetStatus() (status models.LoopStatus) {
	l.state.statusMu.RLock()
	defer l.state.statusMu.RUnlock()
	return l.state.status
}

func (l *looper) SetStatus(status models.LoopStatus) (outcome string, err error) {
	l.state.statusMu.Lock()
	defer l.state.statusMu.Unlock()
	existingStatus := l.state.status

	switch status {
	case constants.Running:
		switch existingStatus {
		case constants.Starting, constants.Running, constants.Stopping, constants.Crashed:
			return fmt.Sprintf("already %s", existingStatus), nil
		}
		l.loopLock.Lock()
		defer l.loopLock.Unlock()
		l.state.status = constants.Starting
		l.state.statusMu.Unlock()
		l.start <- struct{}{}
		newStatus := <-l.running
		l.state.statusMu.Lock()
		l.state.status = newStatus
		return newStatus.String(), nil
	case constants.Stopped:
		switch existingStatus {
		case constants.Starting, constants.Stopping, constants.Stopped, constants.Crashed:
			return fmt.Sprintf("already %s", existingStatus), nil
		}
		l.loopLock.Lock()
		defer l.loopLock.Unlock()
		l.state.status = constants.Stopping
		l.state.statusMu.Unlock()
		l.stop <- struct{}{}
		<-l.stopped
		l.state.statusMu.Lock()
		l.state.status = constants.Stopped
		return status.String(), nil
	default:
		return "", fmt.Errorf("status %q can only be %q or %q",
			status, constants.Running, constants.Stopped)
	}
}

func (l *looper) GetSettings() (settings settings.OpenVPN) {
	l.state.settingsMu.RLock()
	defer l.state.settingsMu.RUnlock()
	return l.state.settings
}

func (l *looper) SetSettings(settings settings.OpenVPN) (outcome string) {
	l.state.settingsMu.Lock()
	settingsUnchanged := reflect.DeepEqual(l.state.settings, settings)
	if settingsUnchanged {
		l.state.settingsMu.Unlock()
		return "settings left unchanged"
	}
	l.state.settings = settings
	_, _ = l.SetStatus(constants.Stopped)
	outcome, _ = l.SetStatus(constants.Running)
	return outcome
}

func (l *looper) GetServers() (servers models.AllServers) {
	l.state.allServersMu.RLock()
	defer l.state.allServersMu.RUnlock()
	return l.state.allServers
}

func (l *looper) SetServers(servers models.AllServers) {
	l.state.allServersMu.Lock()
	defer l.state.allServersMu.Unlock()
	l.state.allServers = servers
}

func (l *looper) GetPortForwarded() (port uint16) {
	l.state.portForwardedMu.RLock()
	defer l.state.portForwardedMu.RUnlock()
	return port
}
