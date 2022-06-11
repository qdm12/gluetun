package portforward

func (l *Loop) GetPortForwarded() (port uint16) {
	return l.state.GetPortForwarded()
}
