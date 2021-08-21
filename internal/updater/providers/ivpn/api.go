package ivpn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	errBuildRequest          = errors.New("cannot build HTTP request")
	errDoRequest             = errors.New("failed doing HTTP request")
	errHTTPStatusCodeNotOK   = errors.New("HTTP status code not OK")
	errUnmarshalResponseBody = errors.New("failed unmarshaling response body")
	errCloseBody             = errors.New("failed closing HTTP body")
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
}

type apiHostnames struct {
	OpenVPN string `json:"openvpn"`
}

func fetchAPI(ctx context.Context, client *http.Client) (
	data apiData, err error) {
	const url = "https://api.ivpn.net/v4/servers/stats"

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return data, fmt.Errorf("%w: %s", errBuildRequest, err)
	}

	response, err := client.Do(request)
	if err != nil {
		return data, fmt.Errorf("%w: %s", errDoRequest, err)
	}

	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return data, fmt.Errorf("%w: %d %s",
			errHTTPStatusCodeNotOK, response.StatusCode, response.Status)
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&data); err != nil {
		_ = response.Body.Close()
		return data, fmt.Errorf("%w: %s", errUnmarshalResponseBody, err)
	}

	if err := response.Body.Close(); err != nil {
		return data, fmt.Errorf("%w: %s", errCloseBody, err)
	}

	return data, nil
}
