package openvpn

func (l *looper) GetPortForwarded() (port uint16) {
	return l.state.GetPortForwarded()
}
