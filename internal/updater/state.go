package updater

import (
	"fmt"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
)

type state struct {
	status   models.LoopStatus
	period   time.Duration
	statusMu sync.RWMutex
	periodMu sync.RWMutex
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

func (l *looper) GetPeriod() (period time.Duration) {
	l.state.periodMu.RLock()
	defer l.state.periodMu.RUnlock()
	return l.state.period
}

func (l *looper) SetPeriod(period time.Duration) {
	l.state.periodMu.Lock()
	defer l.state.periodMu.Unlock()
	l.state.period = period
	l.updateTicker <- struct{}{}
}
