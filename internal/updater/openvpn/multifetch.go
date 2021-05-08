package openvpn

import (
	"context"
	"net/http"
)

// FetchMultiFiles fetches multiple Openvpn files in parallel and
// parses them to extract each of their host. A mapping from host to
// URL is returned.
func FetchMultiFiles(ctx context.Context, client *http.Client, urls []string) (
	hostToURL map[string]string, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	hostToURL = make(map[string]string, len(urls))

	type Result struct {
		url  string
		host string
	}

	results := make(chan Result)
	defer close(results)
	errors := make(chan error)
	defer close(errors)

	for _, url := range urls {
		go func(url string) {
			host, err := FetchFile(ctx, client, url)
			if err != nil {
				errors <- err
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
		case newErr := <-errors:
			if err == nil { // only assign to the first error
				err = newErr
				cancel() // stop other operations, this will trigger other errors we ignore
			}
		case result := <-results:
			hostToURL[result.host] = result.url
		}
	}

	if err != nil {
		return nil, err
	}

	return hostToURL, nil
}
