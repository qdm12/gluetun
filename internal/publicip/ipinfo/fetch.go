package ipinfo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
)

type Fetch struct {
	client *http.Client
}

func New(client *http.Client) *Fetch {
	return &Fetch{
		client: client,
	}
}

var (
	ErrTooManyRequests = errors.New("too many requests sent for this month")
	ErrBadHTTPStatus   = errors.New("bad HTTP status received")
)

// FetchInfo obtains information on the ip address provided
// using the ipinfo.io API. If the ip is nil, the public IP address
// of the machine is used as the IP.
func (f *Fetch) FetchInfo(ctx context.Context, ip net.IP) (
	result Response, err error) {
	const baseURL = "https://ipinfo.io/"
	url := baseURL
	if ip != nil {
		url += ip.String()
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result, err
	}

	response, err := f.client.Do(request)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
	case http.StatusTooManyRequests, http.StatusForbidden:
		return result, fmt.Errorf("%w from %s: %d %s",
			ErrTooManyRequests, url, response.StatusCode, response.Status)
	default:
		return result, fmt.Errorf("%w from %s: %d %s",
			ErrBadHTTPStatus, url, response.StatusCode, response.Status)
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&result); err != nil {
		return result, fmt.Errorf("cannot decode response: %w", err)
	}

	countryCode := strings.ToLower(result.Country)
	country, ok := constants.CountryCodes()[countryCode]
	if ok {
		result.Country = country
	}

	return result, nil
}
