package surfshark

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants"
)

// Note: no multi-hop and some servers are missing from their API.
func addServersFromAPI(ctx context.Context, client *http.Client,
	hts hostToServer) (err error) {
	data, err := fetchAPI(ctx, client)
	if err != nil {
		return err
	}

	hostToLocation := constants.SurfsharkHostToLocation()

	const tcp, udp = true, true
	for _, serverData := range data {
		locationData := hostToLocation[serverData.Host] // TODO remove in v4
		retroLoc := locationData.RetroLoc               // empty string if the host has no retro-compatible region
		hts.add(serverData.Host, serverData.Region, serverData.Country,
			serverData.Location, retroLoc, tcp, udp)
	}

	return nil
}

var (
	ErrHTTPStatusCodeNotOK   = errors.New("HTTP status code not OK")
	ErrUnmarshalResponseBody = errors.New("failed unmarshaling response body")
)

type serverData struct {
	Host     string `json:"connectionName"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Location string `json:"location"`
}

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
		return nil, fmt.Errorf("%w: %d %s", ErrHTTPStatusCodeNotOK,
			response.StatusCode, response.Status)
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
