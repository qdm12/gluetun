package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"

	"github.com/qdm12/gluetun/internal/models"
)

type ifConfigCo struct {
	client *http.Client
}

func newIfConfigCo(client *http.Client) *ifConfigCo {
	return &ifConfigCo{
		client: client,
	}
}

func (i *ifConfigCo) String() string {
	return string(IfConfigCo)
}

func (i *ifConfigCo) CanFetchAnyIP() bool {
	return true
}

// FetchInfo obtains information on the ip address provided
// using the ifconfig.co/json API. If the ip is the zero value,
// the public IP address of the machine is used as the IP.
func (i *ifConfigCo) FetchInfo(ctx context.Context, ip netip.Addr) (
	result models.PublicIP, err error,
) {
	url := "https://ifconfig.co/json"
	if ip.IsValid() {
		url += "?ip=" + ip.String()
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

	switch response.StatusCode {
	case http.StatusOK:
	case http.StatusTooManyRequests:
		return result, fmt.Errorf("%w from %s: %s",
			ErrTooManyRequests, url, response.Status)
	default:
		return result, fmt.Errorf("%w from %s: %s",
			ErrBadHTTPStatus, url, response.Status)
	}

	decoder := json.NewDecoder(response.Body)
	var data struct {
		IP         netip.Addr `json:"ip,omitempty"`
		Country    string     `json:"country,omitempty"`
		RegionName string     `json:"region_name,omitempty"`
		ZipCode    string     `json:"zip_code,omitempty"`
		City       string     `json:"city,omitempty"`
		Latitude   float32    `json:"latitude,omitempty"`
		Longitude  float32    `json:"longitude,omitempty"`
		Hostname   string     `json:"hostname,omitempty"`
		// Timezone in the form America/Montreal
		Timezone string `json:"time_zone,omitempty"`
		AsnOrg   string `json:"asn_org,omitempty"`
	}
	err = decoder.Decode(&data)
	if err != nil {
		return result, fmt.Errorf("decoding response: %w", err)
	}

	result = models.PublicIP{
		IP:           data.IP,
		Region:       data.RegionName,
		Country:      data.Country,
		City:         data.City,
		Hostname:     data.Hostname,
		Location:     fmt.Sprintf("%f,%f", data.Latitude, data.Longitude),
		Organization: data.AsnOrg,
		PostalCode:   data.ZipCode,
		Timezone:     data.Timezone,
	}
	return result, nil
}
