package hidemyass

import (
	"context"
	"net"
	"time"

	"github.com/qdm12/gluetun/internal/updater/resolver"
)

func resolveHosts(ctx context.Context, presolver resolver.Parallel,
	hosts []string, minServers int) (
	hostToIPs map[string][]net.IP, warnings []string, err error) {
	const (
		maxFailRatio    = 0.1
		maxDuration     = 15 * time.Second
		betweenDuration = 2 * time.Second
		maxNoNew        = 2
		maxFails        = 2
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
	return presolver.Resolve(ctx, hosts, settings)
}
