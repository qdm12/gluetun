package expressvpn

import (
	"time"

	"github.com/qdm12/gluetun/internal/updater/resolver"
)

func newParallelResolver() *resolver.Parallel {
	const (
		maxFailRatio = 0.1
		maxNoNew     = 1
		maxFails     = 2
	)
	settings := resolver.ParallelSettings{
		MaxFailRatio: maxFailRatio,
		Repeat: resolver.RepeatSettings{
			MaxDuration: time.Second,
			MaxNoNew:    maxNoNew,
			MaxFails:    maxFails,
			SortIPs:     true,
		},
	}
	return resolver.NewParallelResolver(settings)
}
