package healthcheck

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/healthcheck/dns"
	"github.com/qdm12/gluetun/internal/healthcheck/icmp"
)

type Checker struct {
	tlsDialAddr string
	dialer      *net.Dialer
	echoer      *icmp.Echoer
	dnsClient   *dns.Client
	logger      Logger
	icmpTarget  netip.Addr
	configMutex sync.Mutex

	icmpNotPermitted bool
	smallCheckName   string

	// Internal periodic service signals
	stop context.CancelFunc
	done <-chan struct{}
}

func NewChecker(logger Logger) *Checker {
	return &Checker{
		dialer: &net.Dialer{
			Resolver: &net.Resolver{
				PreferGo: true,
			},
		},
		echoer:    icmp.NewEchoer(logger),
		dnsClient: dns.New(),
		logger:    logger,
	}
}

// SetConfig sets the TCP+TLS dial address and the ICMP echo IP address
// to target by the [Checker].
// This function MUST be called before calling [Checker.Start].
func (c *Checker) SetConfig(tlsDialAddr string, icmpTarget netip.Addr) {
	c.configMutex.Lock()
	defer c.configMutex.Unlock()
	c.tlsDialAddr = tlsDialAddr
	c.icmpTarget = icmpTarget
}

// Start starts the checker by first running a blocking 2s-timed TCP+TLS check,
// and, on success, starts the periodic checks in a separate goroutine:
// - a "small" ICMP echo check every 15 seconds
// - a "full" TCP+TLS check every 5 minutes
// It returns a channel `runError` that receives an error (nil or not) when a periodic check is performed.
// It returns an error if the initial TCP+TLS check fails.
// The Checker has to be ultimately stopped by calling [Checker.Stop].
func (c *Checker) Start(ctx context.Context) (runError <-chan error, err error) {
	if c.tlsDialAddr == "" || c.icmpTarget.IsUnspecified() {
		panic("call Checker.SetConfig with non empty values before Checker.Start")
	}

	// connection isn't under load yet when the checker starts, so a short
	// 6 seconds timeout suffices and provides quick enough feedback that
	// the new connection is not working.
	const timeout = 6 * time.Second
	tcpTLSCheckCtx, tcpTLSCheckCancel := context.WithTimeout(ctx, timeout)
	err = tcpTLSCheck(tcpTLSCheckCtx, c.dialer, c.tlsDialAddr)
	tcpTLSCheckCancel()
	if err != nil {
		return nil, fmt.Errorf("startup check: %w", err)
	}

	ready := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	c.stop = cancel
	done := make(chan struct{})
	c.done = done
	c.smallCheckName = "ICMP echo"
	const smallCheckPeriod = time.Minute
	smallCheckTimer := time.NewTimer(smallCheckPeriod)
	const fullCheckPeriod = 5 * time.Minute
	fullCheckTimer := time.NewTimer(fullCheckPeriod)
	runErrorCh := make(chan error)
	runError = runErrorCh
	go func() {
		defer close(done)
		close(ready)
		for {
			select {
			case <-ctx.Done():
				fullCheckTimer.Stop()
				smallCheckTimer.Stop()
				return
			case <-smallCheckTimer.C:
				err := c.smallPeriodicCheck(ctx)
				if err != nil {
					err = fmt.Errorf("small periodic check: %w", err)
				}
				runErrorCh <- err
				smallCheckTimer.Reset(smallCheckPeriod)
			case <-fullCheckTimer.C:
				err := c.fullPeriodicCheck(ctx)
				if err != nil {
					err = fmt.Errorf("full periodic check: %w", err)
				}
				runErrorCh <- err
				fullCheckTimer.Reset(fullCheckPeriod)
			}
		}
	}()
	<-ready
	return runError, nil
}

func (c *Checker) Stop() error {
	c.stop()
	<-c.done
	c.icmpTarget = netip.Addr{}
	return nil
}

func (c *Checker) smallPeriodicCheck(ctx context.Context) error {
	c.configMutex.Lock()
	ip := c.icmpTarget
	c.configMutex.Unlock()
	const maxTries = 3
	const timeout = 10 * time.Second
	const extraTryTime = 10 * time.Second // 10s added for each subsequent retry
	check := func(ctx context.Context) error {
		if c.icmpNotPermitted {
			return c.dnsClient.Check(ctx)
		}
		err := c.echoer.Echo(ctx, ip)
		if errors.Is(err, icmp.ErrNotPermitted) {
			c.icmpNotPermitted = true
			c.smallCheckName = "plain DNS over UDP"
			c.logger.Warnf("%s; permanently falling back to plaintext DNS checks.", err)
			return c.dnsClient.Check(ctx)
		}
		return err
	}
	return withRetries(ctx, maxTries, timeout, extraTryTime, c.logger, c.smallCheckName, check)
}

func (c *Checker) fullPeriodicCheck(ctx context.Context) error {
	const maxTries = 2
	// 20s timeout in case the connection is under stress
	// See https://github.com/qdm12/gluetun/issues/2270
	const timeout = 20 * time.Second
	const extraTryTime = 10 * time.Second // 10s added for each subsequent retry
	check := func(ctx context.Context) error {
		return tcpTLSCheck(ctx, c.dialer, c.tlsDialAddr)
	}
	return withRetries(ctx, maxTries, timeout, extraTryTime, c.logger, "TCP+TLS dial", check)
}

func tcpTLSCheck(ctx context.Context, dialer *net.Dialer, targetAddress string) error {
	// TODO use mullvad API if current provider is Mullvad

	address, err := makeAddressToDial(targetAddress)
	if err != nil {
		return err
	}

	const dialNetwork = "tcp4"
	connection, err := dialer.DialContext(ctx, dialNetwork, address)
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

var ErrAllCheckTriesFailed = errors.New("all check tries failed")

func withRetries(ctx context.Context, maxTries uint, tryTimeout, extraTryTime time.Duration,
	warner Logger, checkName string, check func(ctx context.Context) error,
) error {
	try := uint(0)
	for {
		timeout := tryTimeout + time.Duration(try)*extraTryTime //nolint:gosec
		checkCtx, cancel := context.WithTimeout(ctx, timeout)
		err := check(checkCtx)
		cancel()
		switch {
		case err == nil:
			return nil
		case ctx.Err() != nil:
			return fmt.Errorf("%s context error: %w", checkName, ctx.Err())
		default:
			warner.Warnf("%s attempt %d/%d failed: %v", checkName, try+1, maxTries, err)
			try++
			if try == maxTries {
				return fmt.Errorf("%w: %s: after %d attempts", ErrAllCheckTriesFailed, checkName, maxTries)
			}
		}
	}
}
