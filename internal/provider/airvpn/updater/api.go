package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/qdm12/gluetun/internal/provider/common"
)

type apiData struct {
	Servers []apiServer `json:"servers"`
}

type apiServer struct {
	PublicName  string `json:"public_name"`
	CountryName string `json:"country_name"`
	CountryCode string `json:"country_code"`
	Location    string `json:"location"`
	Continent   string `json:"continent"`
	IPv4In1     net.IP `json:"ip_v4_in1"`
	IPv4In2     net.IP `json:"ip_v4_in2"`
	IPv4In3     net.IP `json:"ip_v4_in3"`
	IPv4In4     net.IP `json:"ip_v4_in4"`
	IPv6In1     net.IP `json:"ip_v6_in1"`
	IPv6In2     net.IP `json:"ip_v6_in2"`
	IPv6In3     net.IP `json:"ip_v6_in3"`
	IPv6In4     net.IP `json:"ip_v6_in4"`
	Health      string `json:"health"`
}

func fetchAPI(ctx context.Context, client *http.Client) (
	data apiData, err error) {
	const url = "https://airvpn.org/api/status/"

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return data, fmt.Errorf("creating HTTP request: %w", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return data, fmt.Errorf("doing HTTP request: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return data, fmt.Errorf("%w: %d %s",
			common.ErrHTTPStatusCodeNotOK, response.StatusCode, response.Status)
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&data); err != nil {
		_ = response.Body.Close()
		return data, fmt.Errorf("unmarshaling response body: %w", err)
	}

	if err := response.Body.Close(); err != nil {
		return data, fmt.Errorf("closing response body: %w", err)
	}

	return data, nil
}
