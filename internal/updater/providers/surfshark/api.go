package surfshark

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrHTTPStatusCodeNotOK   = errors.New("HTTP status code not OK")
	ErrUnmarshalResponseBody = errors.New("failed unmarshaling response body")
)

//nolint:unused
type serverData struct {
	Host     string `json:"connectionName"`
	Country  string `json:"country"`
	Location string `json:"location"`
}

//nolint:unused,deadcode
func fetchAPI(ctx context.Context, client *http.Client) (
	servers []serverData, err error) {
	const url = "https://my.surfshark.com/vpn/api/v4/server/clusters"

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
	if err := decoder.Decode(&servers); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnmarshalResponseBody, err)
	}

	if err := response.Body.Close(); err != nil {
		return nil, err
	}

	return servers, nil
}
