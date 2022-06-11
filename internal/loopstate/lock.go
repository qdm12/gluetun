package loopstate

func (s *State) Lock()   { s.loopMu.Lock() }
func (s *State) Unlock() { s.loopMu.Unlock() }
