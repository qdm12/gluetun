package updater

import (
	"time"

	"github.com/qdm12/gluetun/internal/updater/resolver"
)

func parallelResolverSettings(hosts []string) (settings resolver.ParallelSettings) {
	const (
		maxFailRatio = 0.1
		maxNoNew     = 1
		maxFails     = 4
		maxDuration  = 3 * time.Second
	)
	return resolver.ParallelSettings{
		Hosts:        hosts,
		MaxFailRatio: maxFailRatio,
		Repeat: resolver.RepeatSettings{
			MaxDuration: maxDuration,
			MaxNoNew:    maxNoNew,
			MaxFails:    maxFails,
			SortIPs:     true,
		},
	}
}
