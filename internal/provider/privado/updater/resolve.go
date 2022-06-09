package privado

import (
	"time"

	"github.com/qdm12/gluetun/internal/updater/resolver"
)

func newParallelResolver() (parallelResolver resolver.Parallel) {
	const (
		maxFailRatio = 0.1
		maxDuration  = 30 * time.Second
		maxNoNew     = 2
		maxFails     = 2
	)
	settings := resolver.ParallelSettings{
		MaxFailRatio: maxFailRatio,
		Repeat: resolver.RepeatSettings{
			MaxDuration: maxDuration,
			MaxNoNew:    maxNoNew,
			MaxFails:    maxFails,
			SortIPs:     true,
		},
	}
	return resolver.NewParallelResolver(settings)
}
