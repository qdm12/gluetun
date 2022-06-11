package unzip

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrHTTPStatusCodeNotOK = errors.New("HTTP status code not OK")
)

func (u *Unzipper) FetchAndExtract(ctx context.Context, url string) (
	contents map[string][]byte, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := u.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %s: %d %s", ErrHTTPStatusCodeNotOK,
			url, response.StatusCode, response.Status)
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if err := response.Body.Close(); err != nil {
		return nil, err
	}

	return zipExtractAll(b)
}
