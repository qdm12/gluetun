package ipinfo

import (
	"context"
	"net"
)

// FetchMultiInfo obtains the public IP address information for every IP
// addresses provided and returns a slice of results with the corresponding
// order as to the IP addresses slice order.
// If an error is encountered, all the operations are canceled and
// an error is returned, so the results returned should be considered
// incomplete in this case.
func (f *Fetch) FetchMultiInfo(ctx context.Context, ips []net.IP) (
	results []Response, err error) {
	ctx, cancel := context.WithCancel(ctx)

	type asyncResult struct {
		index  int
		result Response
		err    error
	}
	resultsCh := make(chan asyncResult)

	for i, ip := range ips {
		go func(index int, ip net.IP) {
			aResult := asyncResult{
				index: index,
			}
			aResult.result, aResult.err = f.FetchInfo(ctx, ip)
			resultsCh <- aResult
		}(i, ip)
	}

	results = make([]Response, len(ips))
	for i := 0; i < len(ips); i++ {
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
