package dns

import (
	"context"
	"fmt"

	"github.com/qdm12/dns/v2/pkg/dot"
	cachemiddleware "github.com/qdm12/dns/v2/pkg/middlewares/cache"
	"github.com/qdm12/dns/v2/pkg/middlewares/cache/lru"
	filtermiddleware "github.com/qdm12/dns/v2/pkg/middlewares/filter"
	"github.com/qdm12/dns/v2/pkg/middlewares/filter/mapfilter"
	"github.com/qdm12/dns/v2/pkg/provider"
	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (l *Loop) GetSettings() (settings settings.DNS) { return l.state.GetSettings() }

func (l *Loop) SetSettings(ctx context.Context, settings settings.DNS) (
	outcome string,
) {
	return l.state.SetSettings(ctx, settings)
}

func buildDoTSettings(settings settings.DNS,
	filter *mapfilter.Filter, logger Logger) (
	dotSettings dot.ServerSettings, err error,
) {
	var middlewares []dot.Middleware

	if *settings.DoT.Caching {
		lruCache, err := lru.New(lru.Settings{})
		if err != nil {
			return dot.ServerSettings{}, fmt.Errorf("creating LRU cache: %w", err)
		}
		cacheMiddleware, err := cachemiddleware.New(cachemiddleware.Settings{
			Cache: lruCache,
		})
		if err != nil {
			return dot.ServerSettings{}, fmt.Errorf("creating cache middleware: %w", err)
		}
		middlewares = append(middlewares, cacheMiddleware)
	}

	filterMiddleware, err := filtermiddleware.New(filtermiddleware.Settings{
		Filter: filter,
	})
	if err != nil {
		return dot.ServerSettings{}, fmt.Errorf("creating filter middleware: %w", err)
	}
	middlewares = append(middlewares, filterMiddleware)

	providersData := provider.NewProviders()
	providers := make([]provider.Provider, len(settings.DoT.Providers))
	for i := range settings.DoT.Providers {
		var err error
		providers[i], err = providersData.Get(settings.DoT.Providers[i])
		if err != nil {
			panic(err) // this should already had been checked
		}
	}

	ipVersion := "ipv4"
	if *settings.DoT.IPv6 {
		ipVersion = "ipv6"
	}
	return dot.ServerSettings{
		Resolver: dot.ResolverSettings{
			UpstreamResolvers: providers,
			IPVersion:         ipVersion,
			Warner:            logger,
		},
		Middlewares: middlewares,
		Logger:      logger,
	}, nil
}
