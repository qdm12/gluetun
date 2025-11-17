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

func buildServerSettings(userSettings settings.DNS,
	filter *mapfilter.Filter, localResolvers []netip.Addr,
	logger Logger) (
	serverSettings server.Settings, err error,
) {
	serverSettings.Logger = logger

	upstreamResolvers := buildProviders(userSettings)

	ipVersion := "ipv4"
	if *userSettings.IPv6 {
		ipVersion = "ipv6"
	}

	var dialer server.Dialer
	switch userSettings.UpstreamType {
	case settings.DNSUpstreamTypeDot:
		dialerSettings := dot.Settings{
			UpstreamResolvers: upstreamResolvers,
			IPVersion:         ipVersion,
		}
		dialer, err = dot.New(dialerSettings)
		if err != nil {
			return server.Settings{}, fmt.Errorf("creating DNS over TLS dialer: %w", err)
		}
	case settings.DNSUpstreamTypeDoh:
		dialerSettings := doh.Settings{
			UpstreamResolvers: upstreamResolvers,
			IPVersion:         ipVersion,
		}
		dialer, err = doh.New(dialerSettings)
		if err != nil {
			return server.Settings{}, fmt.Errorf("creating DNS over HTTPS dialer: %w", err)
		}
	case settings.DNSUpstreamTypePlain:
		dialerSettings := plain.Settings{
			UpstreamResolvers: upstreamResolvers,
			IPVersion:         ipVersion,
		}
		dialer, err = plain.New(dialerSettings)
		if err != nil {
			return server.Settings{}, fmt.Errorf("creating plain DNS dialer: %w", err)
		}
	default:
		panic("unknown upstream type: " + userSettings.UpstreamType)
	}
	serverSettings.Dialer = dialer

	if *userSettings.Caching {
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

	localResolversAddrPorts := make([]netip.AddrPort, len(localResolvers))
	const defaultDNSPort = 53
	for i, addr := range localResolvers {
		localResolversAddrPorts[i] = netip.AddrPortFrom(addr, defaultDNSPort)
	}
	localDNSMiddleware, err := localdns.New(localdns.Settings{
		Resolvers: localResolversAddrPorts, // auto-detected at container start only
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

func buildProviders(userSettings settings.DNS) []provider.Provider {
	if userSettings.UpstreamType == settings.DNSUpstreamTypePlain &&
		len(userSettings.UpstreamPlainAddresses) > 0 {
		providers := make([]provider.Provider, len(userSettings.UpstreamPlainAddresses))
		for i, addrPort := range userSettings.UpstreamPlainAddresses {
			providers[i] = provider.Provider{
				Name: addrPort.String(),
			}
			if addrPort.Addr().Is4() {
				providers[i].Plain.IPv4 = []netip.AddrPort{addrPort}
			} else {
				providers[i].Plain.IPv6 = []netip.AddrPort{addrPort}
			}
		}
	}

	providersData := provider.NewProviders()
	providers := make([]provider.Provider, 0, len(userSettings.Providers)+len(userSettings.UpstreamPlainAddresses))
	for _, providerName := range userSettings.Providers {
		provider, err := providersData.Get(providerName)
		if err != nil {
			panic(err) // this should already had been checked
		}
		providers = append(providers, provider)
	}

	if userSettings.UpstreamType != settings.DNSUpstreamTypePlain {
		return providers
	}

	for _, addrPort := range userSettings.UpstreamPlainAddresses {
		newProvider := provider.Provider{
			Name: addrPort.String(),
		}
		if addrPort.Addr().Is4() {
			newProvider.Plain.IPv4 = []netip.AddrPort{addrPort}
		} else {
			newProvider.Plain.IPv6 = []netip.AddrPort{addrPort}
		}
		providers = append(providers, newProvider)
	}

	return providers
}
