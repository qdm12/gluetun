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

func (s *Server) onUnhealthyVPN(ctx context.Context, lastErrMessage string) {
	s.logger.Info("program has been unhealthy for " +
		s.vpn.healthyWait.String() + ": restarting VPN (healthcheck error: " + lastErrMessage + ")")
	s.logger.Info("👉 See https://github.com/qdm12/gluetun-wiki/blob/main/faq/healthcheck.md")
	s.logger.Info("DO NOT OPEN AN ISSUE UNLESS YOU READ AND TRIED EACH POSSIBLE SOLUTION")
	_, _ = s.vpn.loop.ApplyStatus(ctx, constants.Stopped)
	_, _ = s.vpn.loop.ApplyStatus(ctx, constants.Running)
	s.vpn.healthyWait += *s.config.VPN.Addition
	s.vpn.healthyTimer = time.NewTimer(s.vpn.healthyWait)
}
