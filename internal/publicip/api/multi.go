package api

import (
	"context"
	"net/netip"

	"github.com/qdm12/gluetun/internal/models"
)

// FetchMultiInfo obtains the public IP address information for every IP
// addresses provided and returns a slice of results with the corresponding
// order as to the IP addresses slice order.
// If an error is encountered, all the operations are canceled and
// an error is returned, so the results returned should be considered
// incomplete in this case.
func FetchMultiInfo(ctx context.Context, fetcher InfoFetcher, ips []netip.Addr) (
	results []models.PublicIP, err error,
) {
	ctx, cancel := context.WithCancel(ctx)

	type asyncResult struct {
		index  int
		result models.PublicIP
		err    error
	}
	resultsCh := make(chan asyncResult)

	for i, ip := range ips {
		go func(index int, ip netip.Addr) {
			aResult := asyncResult{
				index: index,
			}
			aResult.result, aResult.err = fetcher.FetchInfo(ctx, ip)
			resultsCh <- aResult
		}(i, ip)
	}

	results = make([]models.PublicIP, len(ips))
	for range len(ips) {
		aResult := <-resultsCh
		if aResult.err != nil {
			if err == nil {
				// Cancel on the first error encountered
				err = aResult.err
				cancel()
			}
			continue // ignore errors after the first one
		}

		results[aResult.index] = aResult.result
	}

	close(resultsCh)
	cancel()

	return results, err
}
