package openvpn

import (
	"context"
	"net/http"
)

// FetchMultiFiles fetches multiple Openvpn files in parallel and
// parses them to extract each of their host. A mapping from host to
// URL is returned.
func FetchMultiFiles(ctx context.Context, client *http.Client, urls []string,
	failEarly bool) (hostToURL map[string]string, errors []error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	hostToURL = make(map[string]string, len(urls))

	type Result struct {
		url  string
		host string
	}

	results := make(chan Result)
	defer close(results)
	errorsCh := make(chan error)
	defer close(errorsCh)

	for _, url := range urls {
		go func(url string) {
			host, err := FetchFile(ctx, client, url)
			if err != nil {
				errorsCh <- err
				return
			}
			results <- Result{
				url:  url,
				host: host,
			}
		}(url)
	}

	for range urls {
		select {
		case result := <-results:
			hostToURL[result.host] = result.url
		case err := <-errorsCh:
			if !failEarly {
				errors = append(errors, err)
				break
			}

			if len(errors) == 0 {
				errors = []error{err} // keep only the first error
				// stop other operations, this will trigger other errors we ignore
				cancel()
			}
		}
	}

	if len(errors) > 0 && failEarly {
		// we don't care about the result found
		return nil, errors
	}

	return hostToURL, errors
}
