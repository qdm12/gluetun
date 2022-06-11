package healthcheck

import (
	"context"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
)

type vpnHealth struct {
	loop         StatusApplier
	healthyWait  time.Duration
	healthyTimer *time.Timer
}

func (s *Server) onUnhealthyVPN(ctx context.Context) {
	s.logger.Info("program has been unhealthy for " +
		s.vpn.healthyWait.String() + ": restarting VPN")
	_, _ = s.vpn.loop.ApplyStatus(ctx, constants.Stopped)
	_, _ = s.vpn.loop.ApplyStatus(ctx, constants.Running)
	s.vpn.healthyWait += *s.config.VPN.Addition
	s.vpn.healthyTimer = time.NewTimer(s.vpn.healthyWait)
}
