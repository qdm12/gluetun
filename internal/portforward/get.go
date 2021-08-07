package portforward

import "github.com/qdm12/gluetun/internal/portforward/state"

type Getter = state.PortForwardedGetter

func (l *Loop) GetPortForwarded() (port uint16) {
	return l.state.GetPortForwarded()
}
