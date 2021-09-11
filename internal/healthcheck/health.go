// Package healthcheck defines the client and server side of the built in healthcheck.
package healthcheck

import (
	"context"
	"time"
)

func (s *Server) runHealthcheckLoop(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	s.vpn.healthyTimer = time.NewTimer(s.vpn.healthyWait)

	for {
		previousErr := s.handler.getErr()

		err := healthCheck(ctx, s.pinger)
		s.handler.setErr(err)

		if previousErr != nil && err == nil {
			s.logger.Info("healthy!")
			s.vpn.healthyTimer.Stop()
			s.vpn.healthyWait = s.config.VPN.Initial
		} else if previousErr == nil && err != nil {
			s.logger.Info("unhealthy: " + err.Error())
			s.vpn.healthyTimer.Stop()
			s.vpn.healthyTimer = time.NewTimer(s.vpn.healthyWait)
		}

		if err != nil { // try again after 1 second
			timer := time.NewTimer(time.Second)
			select {
			case <-ctx.Done():
				if !timer.Stop() {
					<-timer.C
				}
				return
			case <-timer.C:
			case <-s.vpn.healthyTimer.C:
				s.onUnhealthyVPN(ctx)
			}
			continue
		}

		// Success, check again in 5 seconds
		const period = 5 * time.Second
		timer := time.NewTimer(period)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return
		case <-timer.C:
		}
	}
}

func healthCheck(ctx context.Context, pinger Pinger) (err error) {
	// TODO use mullvad API if current provider is Mullvad
	// If we run without root, you need to run this on the gluetun binary:
	// setcap cap_net_raw=+ep /path/to/your/compiled/binary
	// Alternatively, we could have a separate binary just for healthcheck to
	// reduce attack surface.
	errCh := make(chan error)
	go func() {
		errCh <- pinger.Run()
	}()

	select {
	case <-ctx.Done():
		pinger.Stop()
		<-errCh
		return ctx.Err()
	case err = <-errCh:
		return err
	}
}
