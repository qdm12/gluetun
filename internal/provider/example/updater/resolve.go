package updater

import (
	"time"

	"github.com/qdm12/gluetun/internal/updater/resolver"
)

// TODO: remove this file if the parallel resolver is not used
// by the updater.
func parallelResolverSettings(hosts []string) (settings resolver.ParallelSettings) {
	// TODO: adapt these constant values below to make the resolution
	// as fast and as reliable as possible.
	const (
		maxFailRatio    = 0.1
		maxDuration     = 20 * time.Second
		betweenDuration = time.Second
		maxNoNew        = 2
		maxFails        = 2
	)
	return resolver.ParallelSettings{
		Hosts:        hosts,
		MaxFailRatio: maxFailRatio,
		Repeat: resolver.RepeatSettings{
			MaxDuration:     maxDuration,
			BetweenDuration: betweenDuration,
			MaxNoNew:        maxNoNew,
			MaxFails:        maxFails,
			SortIPs:         true,
		},
	}
}
