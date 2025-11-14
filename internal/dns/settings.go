package dns

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/qdm12/dns/v2/pkg/doh"
	"github.com/qdm12/dns/v2/pkg/dot"
	cachemiddleware "github.com/qdm12/dns/v2/pkg/middlewares/cache"
	"github.com/qdm12/dns/v2/pkg/middlewares/cache/lru"
	filtermiddleware "github.com/qdm12/dns/v2/pkg/middlewares/filter"
	"github.com/qdm12/dns/v2/pkg/middlewares/filter/mapfilter"
	"github.com/qdm12/dns/v2/pkg/middlewares/localdns"
	"github.com/qdm12/dns/v2/pkg/plain"
	"github.com/qdm12/dns/v2/pkg/provider"
	"github.com/qdm12/dns/v2/pkg/server"
	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (l *Loop) GetSettings() (settings settings.DNS) { return l.state.GetSettings() }

func (l *Loop) SetSettings(ctx context.Context, settings settings.DNS) (
	outcome string,
) {
	return l.state.SetSettings(ctx, settings)
}

func buildServerSettings(settings settings.DNS,
	filter *mapfilter.Filter, localResolvers []netip.AddrPort,
	logger Logger) (
	serverSettings server.Settings, err error,
) {
	serverSettings.Logger = logger

	providersData := provider.NewProviders()
	upstreamResolvers := make([]provider.Provider, len(settings.Providers))
	for i := range settings.Providers {
		var err error
		upstreamResolvers[i], err = providersData.Get(settings.Providers[i])
		if err != nil {
			panic(err) // this should already had been checked
		}
	}

	ipVersion := "ipv4"
	if *settings.IPv6 {
		ipVersion = "ipv6"
	}

	var dialer server.Dialer
	switch settings.UpstreamType {
	case "dot":
		dialerSettings := dot.Settings{
			UpstreamResolvers: upstreamResolvers,
			IPVersion:         ipVersion,
		}
		dialer, err = dot.New(dialerSettings)
		if err != nil {
			return server.Settings{}, fmt.Errorf("creating DNS over TLS dialer: %w", err)
		}
	case "doh":
		dialerSettings := doh.Settings{
			UpstreamResolvers: upstreamResolvers,
			IPVersion:         ipVersion,
		}
		dialer, err = doh.New(dialerSettings)
		if err != nil {
			return server.Settings{}, fmt.Errorf("creating DNS over HTTPS dialer: %w", err)
		}
	case "plain":
		dialerSettings := plain.Settings{
			UpstreamResolvers: upstreamResolvers,
			IPVersion:         ipVersion,
		}
		dialer, err = plain.New(dialerSettings)
		if err != nil {
			return server.Settings{}, fmt.Errorf("creating plain DNS dialer: %w", err)
		}
	default:
		panic("unknown upstream type: " + settings.UpstreamType)
	}
	serverSettings.Dialer = dialer

	if *settings.Caching {
		lruCache, err := lru.New(lru.Settings{})
		if err != nil {
			return server.Settings{}, fmt.Errorf("creating LRU cache: %w", err)
		}
		cacheMiddleware, err := cachemiddleware.New(cachemiddleware.Settings{
			Cache: lruCache,
		})
		if err != nil {
			return server.Settings{}, fmt.Errorf("creating cache middleware: %w", err)
		}
		serverSettings.Middlewares = append(serverSettings.Middlewares, cacheMiddleware)
	}

	filterMiddleware, err := filtermiddleware.New(filtermiddleware.Settings{
		Filter: filter,
	})
	if err != nil {
		return server.Settings{}, fmt.Errorf("creating filter middleware: %w", err)
	}
	serverSettings.Middlewares = append(serverSettings.Middlewares, filterMiddleware)

	localDNSMiddleware, err := localdns.New(localdns.Settings{
		Resolvers: localResolvers, // auto-detected at container start only
		Logger:    logger,
	})
	if err != nil {
		return server.Settings{}, fmt.Errorf("creating local DNS middleware: %w", err)
	}
	// Place after cache middleware, since we want to avoid caching for local
	// hostnames that may change regularly.
	// Place after filter middleware to avoid conflicts with the rebinding protection.
	serverSettings.Middlewares = append(serverSettings.Middlewares, localDNSMiddleware)

	return serverSettings, nil
}
