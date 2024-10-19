package api

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"strings"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/models"
)

type ResilientFetcher struct {
	fetchers         []Fetcher
	logger           Warner
	fetcherToBanTime map[Fetcher]time.Time
	mutex            sync.RWMutex
	timeNow          func() time.Time
}

// NewResilient creates a 'resilient' fetcher given multiple fetchers.
// For example, it can handle bans and move on to another fetcher if one fails.
func NewResilient(fetchers []Fetcher, logger Warner) *ResilientFetcher {
	return &ResilientFetcher{
		fetchers:         fetchers,
		logger:           logger,
		fetcherToBanTime: make(map[Fetcher]time.Time, len(fetchers)),
		timeNow:          time.Now,
	}
}

func (r *ResilientFetcher) isBanned(fetcher Fetcher) (banned bool) {
	banTime, banned := r.fetcherToBanTime[fetcher]
	if !banned {
		return false
	}
	const banDuration = 30 * 24 * time.Hour
	banExpiryTime := banTime.Add(banDuration)
	now := r.timeNow()
	if now.After(banExpiryTime) {
		delete(r.fetcherToBanTime, fetcher)
		return false
	}
	return true
}

func (r *ResilientFetcher) String() string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, fetcher := range r.fetchers {
		if r.isBanned(fetcher) {
			continue
		}
		return fetcher.String()
	}
	return "<all-banned>"
}

func (r *ResilientFetcher) Token() string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, fetcher := range r.fetchers {
		if r.isBanned(fetcher) {
			continue
		}
		return fetcher.Token()
	}
	return "<all-banned>"
}

// CanFetchAnyIP returns true if any of the fetchers
// can fetch any IP address and is not banned.
func (r *ResilientFetcher) CanFetchAnyIP() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, fetcher := range r.fetchers {
		if !fetcher.CanFetchAnyIP() || r.isBanned(fetcher) {
			continue
		}
		return true
	}
	return false
}

var ErrFetchersAllRateLimited = errors.New("all fetchers are rate limited")

// FetchInfo obtains information on the ip address provided.
// If the ip is the zero value, the public IP address of the machine
// is used as the IP.
// If a fetcher gets banned, the next one is tried – until all have been exhausted.
// Fetchers still within their banned period are skipped.
// If an error unrelated to being banned is encountered, it is returned and more
// fetchers are tried.
func (r *ResilientFetcher) FetchInfo(ctx context.Context, ip netip.Addr) (
	result models.PublicIP, err error,
) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, fetcher := range r.fetchers {
		if r.isBanned(fetcher) ||
			(ip.IsValid() && !fetcher.CanFetchAnyIP()) {
			continue
		}

		result, err = fetcher.FetchInfo(ctx, ip)
		if err == nil || !errors.Is(err, ErrTooManyRequests) {
			return result, err
		}

		// Fetcher is banned
		r.fetcherToBanTime[fetcher] = r.timeNow()
		r.logger.Warn(fetcher.String() + ": " + err.Error())
	}

	fetcherNames := make([]string, len(r.fetchers))
	for i, fetcher := range r.fetchers {
		fetcherNames[i] = fetcher.String()
	}

	return result, fmt.Errorf("%w (%s)",
		ErrFetchersAllRateLimited,
		strings.Join(fetcherNames, ", "))
}

func (r *ResilientFetcher) UpdateFetchers(fetchers []Fetcher) {
	newFetcherNameToFetcher := make(map[string]Fetcher, len(fetchers))
	for _, fetcher := range fetchers {
		newFetcherNameToFetcher[fetcher.String()] = fetcher
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	newFetcherToBanTime := make(map[Fetcher]time.Time, len(r.fetcherToBanTime))
	for bannedFetcher, banTime := range r.fetcherToBanTime {
		if !r.isBanned(bannedFetcher) {
			// fetcher is no longer in its ban period.
			continue
		}
		bannedName := bannedFetcher.String()
		newFetcher, isNewFetcher := newFetcherNameToFetcher[bannedName]
		if isNewFetcher && newFetcher.Token() == bannedFetcher.Token() {
			newFetcherToBanTime[newFetcher] = banTime
		}
	}

	r.fetchers = fetchers
	r.fetcherToBanTime = newFetcherToBanTime
}
