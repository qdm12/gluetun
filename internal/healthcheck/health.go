package healthcheck

import (
	"context"
	"fmt"
	"net"
)

func healthCheck(ctx context.Context, resolver *net.Resolver) (err error) {
	// TODO use mullvad API if current provider is Mullvad
	const domainToResolve = "github.com"
	ips, err := resolver.LookupIP(ctx, "ip", domainToResolve)
	switch {
	case err != nil:
		return fmt.Errorf("cannot resolve github.com: %s", err)
	case len(ips) == 0:
		return fmt.Errorf("resolved no IP addresses for %s", domainToResolve)
	default:
		return nil
	}
}
