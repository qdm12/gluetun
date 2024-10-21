package shadowsocks

import (
	"reflect"
	"sync"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
)

type state struct {
	status     models.LoopStatus
	settings   settings.Shadowsocks
	statusMu   sync.RWMutex
	settingsMu sync.RWMutex
}

func (s *state) setStatusWithLock(status models.LoopStatus) {
	s.statusMu.Lock()
	defer s.statusMu.Unlock()
	s.status = status
}

// GetStatus returns the status of the loop for informative purposes.
// In no case it should be used programmatically to avoid any
// TOCTOU race conditions.
func (l *Loop) GetStatus() (status models.LoopStatus) {
	l.state.statusMu.RLock()
	defer l.state.statusMu.RUnlock()
	return l.state.status
}

func (l *Loop) GetSettings() (settings settings.Shadowsocks) {
	l.state.settingsMu.RLock()
	defer l.state.settingsMu.RUnlock()
	return l.state.settings
}

func (l *Loop) UpdateSettings(updateSettings settings.Shadowsocks) (outcome string) {
	l.state.settingsMu.Lock()
	previousSettings := l.state.settings.Copy()
	l.state.settings.OverrideWith(updateSettings)
	settingsUnchanged := reflect.DeepEqual(previousSettings, l.state.settings)
	l.state.settingsMu.Unlock()
	if settingsUnchanged {
		return "settings left unchanged"
	}
	l.refresh <- struct{}{}
	newStatus := <-l.changed
	l.state.statusMu.Lock()
	l.state.status = newStatus
	l.state.statusMu.Unlock()
	return "settings updated (service " + newStatus.String() + ")"
}
