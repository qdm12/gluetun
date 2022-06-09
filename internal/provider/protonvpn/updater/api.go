package updater

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
)

var (
	ErrHTTPStatusCodeNotOK = errors.New("HTTP status code not OK")
)

type apiData struct {
	LogicalServers []logicalServer
}

type logicalServer struct {
	Name        string
	ExitCountry string
	Region      *string
	City        *string
	Servers     []physicalServer
}

type physicalServer struct {
	EntryIP net.IP
	ExitIP  net.IP
	Domain  string
	Status  uint8
}

func fetchAPI(ctx context.Context, client *http.Client) (
	data apiData, err error) {
	const url = "https://api.protonmail.ch/vpn/logicals"

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return data, err
	}

	response, err := client.Do(request)
	if err != nil {
		return data, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return data, fmt.Errorf("%w: %d %s", ErrHTTPStatusCodeNotOK,
			response.StatusCode, response.Status)
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&data); err != nil {
		return data, fmt.Errorf("failed unmarshaling response body: %w", err)
	}

	if err := response.Body.Close(); err != nil {
		return data, err
	}

	return data, nil
}
