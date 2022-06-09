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

type serverData struct {
	Domain    string `json:"domain"`
	IPAddress string `json:"ip_address"`
	Name      string `json:"name"`
	Country   string `json:"country"`
	Features  struct {
		UDP bool `json:"openvpn_udp"`
		TCP bool `json:"openvpn_tcp"`
	} `json:"features"`
}

func fetchAPI(ctx context.Context, client *http.Client) (data []serverData, err error) {
	const url = "https://nordvpn.com/api/server"

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
		return nil, fmt.Errorf("failed unmarshaling response body: %w", err)
	}

	if err := response.Body.Close(); err != nil {
		return nil, err
	}

	return data, nil
}
