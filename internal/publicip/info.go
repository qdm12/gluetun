package publicip

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

type ipInfoData struct {
	Region  string `json:"region"`
	Country string `json:"country"`
	City    string `json:"city"`
}

var ErrBadHTTPStatus = errors.New("bad HTTP status received")

func Info(ctx context.Context, client *http.Client, ip net.IP) ( //nolint:interfacer
	country, region, city string, err error) {
	const baseURL = "https://ipinfo.io/"
	url := baseURL + ip.String()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", "", err
	}

	response, err := client.Do(request)
	if err != nil {
		return "", "", "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("%w: %d", ErrBadHTTPStatus, response.StatusCode)
	}

	decoder := json.NewDecoder(response.Body)
	var data ipInfoData
	if err := decoder.Decode(&data); err != nil {
		return "", "", "", err
	}

	countryCode := strings.ToLower(data.Country)
	country, ok := constants.CountryCodes()[countryCode]
	if !ok {
		country = data.Country
	}
	return country, data.Region, data.City, nil
}
