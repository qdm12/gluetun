package healthcheck

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

func (s *Server) runHealthcheckLoop(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	s.vpn.healthyTimer = time.NewTimer(s.vpn.healthyWait)

	for {
		previousErr := s.handler.getErr()

		const healthcheckTimeout = 3 * time.Second
		healthcheckCtx, healthcheckCancel := context.WithTimeout(
			ctx, healthcheckTimeout)
		err := s.healthCheck(healthcheckCtx)
		healthcheckCancel()

		s.handler.setErr(err)

		if previousErr != nil && err == nil {
			s.logger.Info("healthy!")
			s.vpn.healthyTimer.Stop()
			s.vpn.healthyWait = *s.config.VPN.Initial
		} else if previousErr == nil && err != nil {
			s.logger.Info("unhealthy: " + err.Error() + "(see https://github.com/qdm12/gluetun/wiki/Healthcheck)")
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

func (s *Server) healthCheck(ctx context.Context) (err error) {
	// TODO use mullvad API if current provider is Mullvad

	address, err := makeAddressToDial(s.config.TargetAddress)
	if err != nil {
		return err
	}

	const dialNetwork = "tcp4"
	connection, err := s.dialer.DialContext(ctx, dialNetwork, address)
	if err != nil {
		return fmt.Errorf("cannot dial: %w", err)
	}

	err = connection.Close()
	if err != nil {
		return fmt.Errorf("cannot close connection: %w", err)
	}

	return nil
}

func makeAddressToDial(address string) (addressToDial string, err error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		addrErr := new(net.AddrError)
		ok := errors.As(err, &addrErr)
		if !ok || addrErr.Err != "missing port in address" {
			return "", fmt.Errorf("cannot split host and port from address: %w", err)
		}
		host = address
		const defaultPort = "443"
		port = defaultPort
	}
	address = net.JoinHostPort(host, port)
	return address, nil
}
