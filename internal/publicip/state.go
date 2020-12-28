package publicip

import (
	"fmt"
	"net"
	"reflect"
	"sync"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/settings"
)

type state struct {
	status     models.LoopStatus
	settings   settings.PublicIP
	ip         net.IP
	statusMu   sync.RWMutex
	settingsMu sync.RWMutex
	ipMu       sync.RWMutex
}

func (s *state) setStatusWithLock(status models.LoopStatus) {
	s.statusMu.Lock()
	defer s.statusMu.Unlock()
	s.status = status
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
		case constants.Stopped, constants.Stopping, constants.Starting, constants.Crashed:
			return fmt.Sprintf("already %s", existingStatus), nil
		}
		l.loopLock.Lock()
		defer l.loopLock.Unlock()
		l.state.status = constants.Stopping
		l.state.statusMu.Unlock()
		l.stop <- struct{}{}
		<-l.stopped
		l.state.statusMu.Lock()
		l.state.status = status
		return status.String(), nil
	default:
		return "", fmt.Errorf("status %q can only be %q or %q",
			status, constants.Running, constants.Stopped)
	}
}

func (l *looper) GetSettings() (settings settings.PublicIP) {
	l.state.settingsMu.RLock()
	defer l.state.settingsMu.RUnlock()
	return l.state.settings
}

func (l *looper) SetSettings(settings settings.PublicIP) (outcome string) {
	l.state.settingsMu.Lock()
	defer l.state.settingsMu.Unlock()
	settingsUnchanged := reflect.DeepEqual(settings, l.state.settings)
	if settingsUnchanged {
		return "settings left unchanged"
	}
	periodChanged := l.state.settings.Period != settings.Period
	l.state.settings = settings
	if periodChanged {
		l.updateTicker <- struct{}{}
		// TODO blocking
	}
	return "settings updated"
}

func (l *looper) GetPublicIP() (publicIP net.IP) {
	l.state.ipMu.RLock()
	defer l.state.ipMu.RUnlock()
	publicIP = make(net.IP, len(l.state.ip))
	copy(publicIP, l.state.ip)
	return publicIP
}

func (s *state) setPublicIP(publicIP net.IP) {
	s.ipMu.Lock()
	defer s.ipMu.Unlock()
	s.ip = make(net.IP, len(publicIP))
	copy(s.ip, publicIP)
}
