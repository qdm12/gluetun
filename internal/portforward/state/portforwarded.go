package state

type PortForwardedGetterSetter interface {
	PortForwardedGetter
	SetPortForwarded(port uint16)
}

type PortForwardedGetter interface {
	GetPortForwarded() (port uint16)
}

// GetPortForwarded is used by the control HTTP server
// to obtain the port currently forwarded.
func (s *State) GetPortForwarded() (port uint16) {
	s.portForwardedMu.RLock()
	defer s.portForwardedMu.RUnlock()
	return s.portForwarded
}

// SetPortForwarded is only used from within the OpenVPN loop
// to set the port forwarded.
func (s *State) SetPortForwarded(port uint16) {
	s.portForwardedMu.Lock()
	defer s.portForwardedMu.Unlock()
	s.portForwarded = port
}
