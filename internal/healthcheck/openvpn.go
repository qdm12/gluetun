package healthcheck

import (
	"context"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/vpn"
)

type vpnHealth struct {
	looper       vpn.Looper
	healthyWait  time.Duration
	healthyTimer *time.Timer
}

func (s *Server) onUnhealthyVPN(ctx context.Context) {
	s.logger.Info("program has been unhealthy for " +
		s.vpn.healthyWait.String() + ": restarting VPN")
	_, _ = s.vpn.looper.ApplyStatus(ctx, constants.Stopped)
	_, _ = s.vpn.looper.ApplyStatus(ctx, constants.Running)
	s.vpn.healthyWait += s.config.VPN.Addition
	s.vpn.healthyTimer = time.NewTimer(s.vpn.healthyWait)
}
