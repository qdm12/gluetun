package hidemyass

import (
	"time"

	"github.com/qdm12/gluetun/internal/updater/resolver"
)

func newParallelResolver() *resolver.Parallel {
	const (
		maxFailRatio    = 0.1
		maxDuration     = 15 * time.Second
		betweenDuration = 2 * time.Second
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
