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

func (l *looper) SetStatus(status models.LoopStatus) (message string) {
	l.state.statusMu.Lock()
	defer l.state.statusMu.Unlock()
	existingStatus := l.state.status

	switch status {
	case constants.Running:
		switch existingStatus {
		case constants.Starting:
			return "already starting"
		case constants.Running:
			return "already running"
		}
		l.loopLock.Lock()
		defer l.loopLock.Unlock()
		l.state.status = constants.Starting
		l.start <- existingStatus
		<-l.running
		l.state.status = constants.Running
		return "running"
	case constants.Stopped:
		switch existingStatus {
		case constants.Stopping:
			return "already stopping"
		case constants.Stopped:
			return "already stopped"
		}
		l.loopLock.Lock()
		defer l.loopLock.Unlock()
		l.state.status = constants.Stopping
		l.stop <- struct{}{}
		<-l.stopped
		l.state.status = constants.Stopped
		return "stopped"
	default:
		return fmt.Sprintf("status %q can only be %q or %q",
			status, constants.Running, constants.Stopped)
	}
}

func (l *looper) GetSettings() (settings settings.OpenVPN) {
	l.state.settingsMu.RLock()
	defer l.state.settingsMu.RUnlock()
	return l.state.settings
}

func (l *looper) SetSettings(settings settings.OpenVPN) {
	l.state.settingsMu.Lock()
	settingsChanged := !reflect.DeepEqual(l.state.settings, settings)
	if settingsChanged {
		l.state.settings = settings
		l.loopLock.Lock()
		defer l.loopLock.Unlock()
		l.state.statusMu.Lock()
		defer l.state.statusMu.Unlock()
		existingStatus := l.state.status
		l.state.status = constants.Starting
		l.start <- existingStatus
		l.state.settingsMu.Unlock() // unlocks to allow the loop to read settings at restart
		<-l.running
		l.state.status = constants.Running
		return
	}
	l.state.settingsMu.Unlock()
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
