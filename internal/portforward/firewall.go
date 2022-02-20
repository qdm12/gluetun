package portforward

import "context"

// firewallBlockPort obtains the state port thread safely and blocks
// it in the firewall if it is not the zero value (0).
func (l *Loop) firewallBlockPort(ctx context.Context) {
	port := l.state.GetPortForwarded()
	if port == 0 {
		return
	}

	err := l.portAllower.RemoveAllowedPort(ctx, port)
	if err != nil {
		l.logger.Error("cannot block previous port in firewall: " + err.Error())
	}
}

// firewallAllowPort obtains the state port thread safely and allows
// it in the firewall if it is not the zero value (0).
func (l *Loop) firewallAllowPort(ctx context.Context) {
	port := l.state.GetPortForwarded()
	if port == 0 {
		return
	}

	startData := l.state.GetStartData()
	err := l.portAllower.SetAllowedPort(ctx, port, startData.Interface)
	if err != nil {
		l.logger.Error("cannot allow port: " + err.Error())
	}
}
