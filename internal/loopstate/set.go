package loopstate

import "github.com/qdm12/gluetun/internal/models"

// SetStatus sets the status thread safely.
// It should only be called by the loop internal code since
// it does not interact with the loop code directly.
func (s *State) SetStatus(status models.LoopStatus) {
	s.statusMu.Lock()
	defer s.statusMu.Unlock()
	s.status = status
}
