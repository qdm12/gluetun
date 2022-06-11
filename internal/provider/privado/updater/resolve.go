package updater

import (
	"time"

	"github.com/qdm12/gluetun/internal/updater/resolver"
)

func parallelResolverSettings(hosts []string) (settings resolver.ParallelSettings) {
	const (
		maxFailRatio = 0.1
		maxDuration  = 30 * time.Second
		maxNoNew     = 2
		maxFails     = 2
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
