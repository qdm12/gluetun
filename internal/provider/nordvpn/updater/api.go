package updater

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrHTTPStatusCodeNotOK = errors.New("HTTP status code not OK")
)

func fetchAPI(ctx context.Context, client *http.Client,
	recommended bool, limit uint) (data []serverData, err error) {
	url := "https://api.nordvpn.com/v1/servers/"
	if recommended {
		url += "recommendations"
	}
	url += fmt.Sprintf("?limit=%d", limit) // 0 means no limit

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %s", ErrHTTPStatusCodeNotOK, response.Status)
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("decoding response body: %w", err)
	}

	if err := response.Body.Close(); err != nil {
		return nil, err
	}

	return data, nil
}
