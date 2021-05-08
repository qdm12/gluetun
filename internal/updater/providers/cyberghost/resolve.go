package cyberghost

import (
	"context"
	"net"
	"time"

	"github.com/qdm12/gluetun/internal/updater/resolver"
)

func resolveHosts(ctx context.Context, presolver resolver.Parallel,
	possibleHosts []string, minServers int) (
	hostToIPs map[string][]net.IP, err error) {
	const (
		maxFailRatio    = 1
		maxDuration     = 10 * time.Second
		betweenDuration = 500 * time.Millisecond
		maxNoNew        = 2
		maxFails        = 10
	)
	settings := resolver.ParallelSettings{
		MaxFailRatio: maxFailRatio,
		MinFound:     minServers,
		Repeat: resolver.RepeatSettings{
			MaxDuration:     maxDuration,
			BetweenDuration: betweenDuration,
			MaxNoNew:        maxNoNew,
			MaxFails:        maxFails,
			SortIPs:         true,
		},
	}
	hostToIPs, _, err = presolver.Resolve(ctx, possibleHosts, settings)
	if err != nil {
		return nil, err
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return hostToIPs, nil
}
