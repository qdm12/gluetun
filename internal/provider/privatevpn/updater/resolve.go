package updater

import (
	"time"

	"github.com/qdm12/gluetun/internal/updater/resolver"
)

func newParallelResolver() (parallelResolver *resolver.Parallel) {
	const (
		maxFailRatio    = 0.1
		maxDuration     = 6 * time.Second
		betweenDuration = time.Second
		maxNoNew        = 2
		maxFails        = 2
	)
	settings := resolver.ParallelSettings{
		MaxFailRatio: maxFailRatio,
		Repeat: resolver.RepeatSettings{
			MaxDuration:     maxDuration,
			BetweenDuration: betweenDuration,
			MaxNoNew:        maxNoNew,
			MaxFails:        maxFails,
			SortIPs:         true,
		},
	}
	return resolver.NewParallelResolver(settings)
}
