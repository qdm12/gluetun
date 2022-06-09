package updater

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	errHTTPStatusCodeNotOK = errors.New("HTTP status code not OK")
)

type apiData struct {
	Servers []apiServer `json:"servers"`
}

type apiServer struct {
	Hostnames apiHostnames `json:"hostnames"`
	IsActive  bool         `json:"is_active"`
	Country   string       `json:"country"`
	City      string       `json:"city"`
	ISP       string       `json:"isp"`
	WgPubKey  string       `json:"wg_public_key"`
}

type apiHostnames struct {
	OpenVPN   string `json:"openvpn"`
	Wireguard string `json:"wireguard"`
}

func fetchAPI(ctx context.Context, client *http.Client) (
	data apiData, err error) {
	const url = "https://api.ivpn.net/v4/servers/stats"

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return data, err
	}

	response, err := client.Do(request)
	if err != nil {
		return data, err
	}

	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return data, fmt.Errorf("%w: %d %s",
			errHTTPStatusCodeNotOK, response.StatusCode, response.Status)
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&data); err != nil {
		_ = response.Body.Close()
		return data, fmt.Errorf("failed unmarshaling response body: %w", err)
	}

	if err := response.Body.Close(); err != nil {
		return data, fmt.Errorf("cannot close response body: %w", err)
	}

	return data, nil
}
