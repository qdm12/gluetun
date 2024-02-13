package ipinfo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
)

type Fetch struct {
	client *http.Client
	token  string
}

func New(client *http.Client, token string) *Fetch {
	return &Fetch{
		client: client,
		token:  token,
	}
}

var (
	ErrTokenNotValid   = errors.New("token is not valid")
	ErrTooManyRequests = errors.New("too many requests sent for this month")
	ErrBadHTTPStatus   = errors.New("bad HTTP status received")
)

// FetchInfo obtains information on the ip address provided
// using the ipinfo.io API. If the ip is the zero value, the public IP address
// of the machine is used as the IP.
func (f *Fetch) FetchInfo(ctx context.Context, ip netip.Addr) (
	result Response, err error) {
	url := "https://ipinfo.io/"
	switch {
	case ip.Is6():
		url = "https://v6.ipinfo.io/" + ip.String()
	case ip.Is4():
		url = "https://ipinfo.io/" + ip.String()
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result, err
	}
	request.Header.Set("Authorization", "Bearer "+f.token)

	response, err := f.client.Do(request)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	if f.token != "" && response.StatusCode == http.StatusUnauthorized {
		return result, fmt.Errorf("%w: %s", ErrTokenNotValid, response.Status)
	}

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
		return result, fmt.Errorf("decoding response: %w", err)
	}

	countryCode := strings.ToLower(result.Country)
	country, ok := constants.CountryCodes()[countryCode]
	if ok {
		result.Country = country
	}

	return result, nil
}
