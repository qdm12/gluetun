package healthcheck

import (
	"context"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/openvpn"
)

type openvpnHealth struct {
	looper       openvpn.Looper
	healthyWait  time.Duration
	healthyTimer *time.Timer
}

func (s *Server) onUnhealthyOpenvpn(ctx context.Context) {
	s.logger.Info("program has been unhealthy for " +
		s.openvpn.healthyWait.String() + ": restarting OpenVPN")
	_, _ = s.openvpn.looper.ApplyStatus(ctx, constants.Stopped)
	_, _ = s.openvpn.looper.ApplyStatus(ctx, constants.Running)
	s.openvpn.healthyWait += s.config.OpenVPN.Addition
	s.openvpn.healthyTimer = time.NewTimer(s.openvpn.healthyWait)
}
