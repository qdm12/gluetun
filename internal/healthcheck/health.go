// Package healthcheck defines the client and server side of the built in healthcheck.
package healthcheck

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

func (s *server) runHealthcheckLoop(ctx context.Context, healthy chan<- bool, done chan<- struct{}) {
	defer close(done)
	for {
		previousErr := s.handler.getErr()

		err := healthCheck(ctx, s.resolver)
		s.handler.setErr(err)

		// Notify the healthy channel, or not if it's already full
		select {
		case healthy <- err == nil:
		default:
		}

		if previousErr != nil && err == nil {
			s.logger.Info("healthy!")
		} else if previousErr == nil && err != nil {
			s.logger.Info("unhealthy: " + err.Error())
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

var (
	errNoIPResolved = errors.New("no IP address resolved")
)

func healthCheck(ctx context.Context, resolver *net.Resolver) (err error) {
	// TODO use mullvad API if current provider is Mullvad
	const domainToResolve = "github.com"
	ips, err := resolver.LookupIP(ctx, "ip", domainToResolve)
	switch {
	case err != nil:
		return err
	case len(ips) == 0:
		return fmt.Errorf("%w for %s", errNoIPResolved, domainToResolve)
	default:
		return nil
	}
}
