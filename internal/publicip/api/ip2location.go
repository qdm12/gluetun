package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
)

type ip2Location struct {
	client *http.Client
	token  string
}

func newIP2Location(client *http.Client, token string) *ip2Location {
	return &ip2Location{
		client: client,
		token:  token,
	}
}

func (i *ip2Location) String() string {
	return string(IP2Location)
}

func (i *ip2Location) CanFetchAnyIP() bool {
	return true
}

func (i *ip2Location) Token() string {
	return i.token
}

// FetchInfo obtains information on the ip address provided
// using the api.ip2location.io API. If the ip is the zero value,
// the public IP address of the machine is used as the IP.
func (i *ip2Location) FetchInfo(ctx context.Context, ip netip.Addr) (
	result models.PublicIP, err error,
) {
	url := "https://api.ip2location.io/"
	if ip.IsValid() {
		url += "?ip=" + ip.String()
	}

	if i.token != "" {
		if !strings.Contains(url, "?") {
			url += "?"
		}
		url += "&key=" + i.token
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result, err
	}

	response, err := i.client.Do(request)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	if i.token != "" && response.StatusCode == http.StatusUnauthorized {
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
	var data struct {
		IP          netip.Addr `json:"ip,omitempty"`
		CountryName string     `json:"country_name,omitempty"`
		RegionName  string     `json:"region_name,omitempty"`
		CityName    string     `json:"city_name,omitempty"`
		Latitude    float32    `json:"latitude,omitempty"`
		Longitude   float32    `json:"longitude,omitempty"`
		ZipCode     string     `json:"zip_code,omitempty"`
		// Timezone in the form -07:00
		Timezone string `json:"time_zone,omitempty"`
		As       string `json:"as,omitempty"`
	}
	if err := decoder.Decode(&data); err != nil {
		return result, fmt.Errorf("decoding response: %w", err)
	}

	result = models.PublicIP{
		IP:           data.IP,
		Region:       data.RegionName,
		Country:      data.CountryName,
		City:         data.CityName,
		Hostname:     "", // no hostname
		Location:     fmt.Sprintf("%f,%f", data.Latitude, data.Longitude),
		Organization: data.As,
		PostalCode:   data.ZipCode,
		Timezone:     data.Timezone,
	}
	return result, nil
}
