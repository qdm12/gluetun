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

type Checker struct {
	targetAddress string
	dialer        *net.Dialer

	// Periodic service
	logger DebugLogger

	// Internal periodic service signals
	stop context.CancelFunc
	done <-chan struct{}
}

func NewChecker(tlsDialAddress string, logger DebugLogger) *Checker {
	return &Checker{
		targetAddress: tlsDialAddress,
		dialer: &net.Dialer{
			Resolver: &net.Resolver{
				PreferGo: true,
			},
		},
		logger: logger,
	}
}

func (c *Checker) Start(ctx context.Context) (runError <-chan error, err error) {
	err = c.check(ctx)
	if err != nil {
		return nil, fmt.Errorf("startup healthcheck: %w", err)
	}
	c.logger.Debug("initial healthcheck successful")

	ready := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	c.stop = cancel
	done := make(chan struct{})
	c.done = done
	const period = 5 * time.Minute
	timer := time.NewTimer(period)
	runErrorCh := make(chan error)
	runError = runErrorCh
	go func() {
		defer close(done)
		close(ready)
		for {
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
				err := c.check(ctx)
				if err != nil {
					runErrorCh <- fmt.Errorf("periodic healthcheck: %w", err)
					return
				}
				c.logger.Debug("healthcheck successful")
				timer.Reset(period)
			}
		}
	}()
	<-ready
	return runError, nil
}

func (c *Checker) Stop() error {
	c.stop()
	<-c.done
	return nil
}

func (c *Checker) check(ctx context.Context) error {
	// 10s timeout in case the connection is under stress
	// See https://github.com/qdm12/gluetun/issues/2270
	const timeout = 10 * time.Second
	ctx, healthcheckCancel := context.WithTimeout(ctx, timeout)
	defer healthcheckCancel()

	// TODO use mullvad API if current provider is Mullvad

	address, err := makeAddressToDial(c.targetAddress)
	if err != nil {
		return err
	}

	const dialNetwork = "tcp4"
	connection, err := c.dialer.DialContext(ctx, dialNetwork, address)
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
