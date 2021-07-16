package dns

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

type state struct {
	status     models.LoopStatus
	settings   configuration.DNS
	statusMu   sync.RWMutex
	settingsMu sync.RWMutex
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

var ErrInvalidStatus = errors.New("invalid status")

func (l *looper) SetStatus(ctx context.Context, status models.LoopStatus) (
	outcome string, err error) {
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
		newStatus := constants.Starting // for canceled context
		select {
		case <-ctx.Done():
		case newStatus = <-l.running:
		}

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

		newStatus := constants.Stopping // for canceled context
		select {
		case <-ctx.Done():
		case <-l.stopped:
			newStatus = constants.Stopped
		}
		l.state.statusMu.Lock()
		l.state.status = newStatus
		return status.String(), nil
	default:
		return "", fmt.Errorf("%w: %s: it can only be one of: %s, %s",
			ErrInvalidStatus, status, constants.Running, constants.Stopped)
	}
}

func (l *looper) GetSettings() (settings configuration.DNS) {
	l.state.settingsMu.RLock()
	defer l.state.settingsMu.RUnlock()
	return l.state.settings
}

func (l *looper) SetSettings(ctx context.Context, settings configuration.DNS) (
	outcome string) {
	l.state.settingsMu.Lock()
	settingsUnchanged := reflect.DeepEqual(l.state.settings, settings)
	if settingsUnchanged {
		l.state.settingsMu.Unlock()
		return "settings left unchanged"
	}
	tempSettings := l.state.settings
	tempSettings.UpdatePeriod = settings.UpdatePeriod
	onlyUpdatePeriodChanged := reflect.DeepEqual(tempSettings, settings)
	l.state.settings = settings
	l.state.settingsMu.Unlock()
	if onlyUpdatePeriodChanged {
		l.updateTicker <- struct{}{}
		return "update period changed"
	}
	_, _ = l.SetStatus(ctx, constants.Stopped)
	if settings.Enabled {
		outcome, _ = l.SetStatus(ctx, constants.Running)
	}
	return outcome
}
