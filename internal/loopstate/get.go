package loopstate

import "github.com/qdm12/gluetun/internal/models"

// GetStatus gets the status thread safely.
func (s *State) GetStatus() (status models.LoopStatus) {
	s.statusMu.RLock()
	defer s.statusMu.RUnlock()
	return s.status
}
