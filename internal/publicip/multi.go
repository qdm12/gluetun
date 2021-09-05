package publicip

import (
	"context"
	"net"
	"net/http"

	"github.com/qdm12/gluetun/internal/publicip/models"
)

// MultiInfo obtains the public IP address information for every IP
// addresses provided and returns a slice of results with the corresponding
// order as to the IP addresses slice order.
// If an error is encountered, all the operations are canceled and
// an error is returned, so the results returned should be considered
// incomplete in this case.
func MultiInfo(ctx context.Context, client *http.Client, ips []net.IP) (
	results []models.IPInfoData, err error) {
	ctx, cancel := context.WithCancel(ctx)

	type asyncResult struct {
		index  int
		result models.IPInfoData
		err    error
	}
	resultsCh := make(chan asyncResult)

	for i, ip := range ips {
		go func(index int, ip net.IP) {
			aResult := asyncResult{
				index: index,
			}
			aResult.result, aResult.err = Info(ctx, client, ip)
			resultsCh <- aResult
		}(i, ip)
	}

	results = make([]models.IPInfoData, len(ips))
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
