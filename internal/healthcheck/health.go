package healthcheck

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

func (s *Server) runHealthcheckLoop(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	timeoutIndex := 0
	healthcheckTimeouts := []time.Duration{
		2 * time.Second,
		4 * time.Second,
		6 * time.Second,
		8 * time.Second,
		// This can be useful when the connection is under stress
		// See https://github.com/qdm12/gluetun/issues/2270
		10 * time.Second,
	}
	s.vpn.healthyTimer = time.NewTimer(s.vpn.healthyWait)

	for {
		previousErr := s.handler.getErr()

		timeout := healthcheckTimeouts[timeoutIndex]
		healthcheckCtx, healthcheckCancel := context.WithTimeout(
			ctx, timeout)
		err := s.healthCheck(healthcheckCtx)
		healthcheckCancel()

		s.handler.setErr(err)

		switch {
		case previousErr != nil && err == nil: // First success
			s.logger.Info("healthy!")
			timeoutIndex = 0
			s.vpn.healthyTimer.Stop()
			s.vpn.healthyWait = *s.config.VPN.Initial
		case previousErr == nil && err != nil: // First failure
			s.logger.Debug("unhealthy: " + err.Error())
			s.vpn.healthyTimer.Stop()
			s.vpn.healthyTimer = time.NewTimer(s.vpn.healthyWait)
		case previousErr != nil && err != nil: // Nth failure
			if timeoutIndex < len(healthcheckTimeouts)-1 {
				timeoutIndex++
			}
			select {
			case <-s.vpn.healthyTimer.C:
				timeoutIndex = 0 // retry next with the smallest timeout
				s.onUnhealthyVPN(ctx, err.Error())
			default:
			}
		case previousErr == nil && err == nil: // Nth success
			timer := time.NewTimer(s.config.SuccessWait)
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
			}
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
		return fmt.Errorf("dialing: %w", err)
	}

	if strings.HasSuffix(address, ":443") {
		host, _, err := net.SplitHostPort(address)
		if err != nil {
			return fmt.Errorf("splitting host and port: %w", err)
		}
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: host,
		}
		tlsConnection := tls.Client(connection, tlsConfig)
		err = tlsConnection.HandshakeContext(ctx)
		if err != nil {
			return fmt.Errorf("running TLS handshake: %w", err)
		}
	}

	err = connection.Close()
	if err != nil {
		return fmt.Errorf("closing connection: %w", err)
	}

	return nil
}

func makeAddressToDial(address string) (addressToDial string, err error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		addrErr := new(net.AddrError)
		ok := errors.As(err, &addrErr)
		if !ok || addrErr.Err != "missing port in address" {
			return "", fmt.Errorf("splitting host and port from address: %w", err)
		}
		host = address
		const defaultPort = "443"
		port = defaultPort
	}
	address = net.JoinHostPort(host, port)
	return address, nil
}
