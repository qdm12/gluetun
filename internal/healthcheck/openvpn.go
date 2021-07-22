package healthcheck

import (
	"context"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
)

func (s *server) onUnhealthyOpenvpn(ctx context.Context) {
	s.logger.Info("program has been unhealthy for " +
		s.openvpn.healthyWait.String() + ": restarting OpenVPN")
	_, _ = s.openvpn.looper.ApplyStatus(ctx, constants.Stopped)
	_, _ = s.openvpn.looper.ApplyStatus(ctx, constants.Running)
	s.openvpn.healthyWait += s.config.OpenVPN.Addition
	s.openvpn.healthyTimer = time.NewTimer(s.openvpn.healthyWait)
}
