package healthcheck

import (
	"context"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
)

func (s *server) onUnhealthyOpenvpn(ctx context.Context) {
	s.logger.Info("program has been unhealthy for " +
		s.openvpn.currentHealthyWait.String() + ": restarting OpenVPN")
	_, _ = s.openvpn.looper.ApplyStatus(ctx, constants.Stopped)
	_, _ = s.openvpn.looper.ApplyStatus(ctx, constants.Running)
	s.openvpn.currentHealthyWait += s.openvpn.healthyWaitConfig.Addition
	s.openvpn.healthyTimer = time.NewTimer(s.openvpn.currentHealthyWait)
}
