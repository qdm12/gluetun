package updater

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrHTTPStatusCodeNotOK = errors.New("HTTP status code not OK")
)

type apiData struct {
	Data []regionData `json:"data"`
}

type regionData struct {
	Region string      `json:"name"`
	Groups []groupData `json:"groups"`
}

type groupData struct {
	City     string       `json:"city"`
	Nodes    []serverData `json:"nodes"`
	OvpnX509 string       `json:"ovpn_x509"`
	WgPubKey string       `json:"wg_pubkey"`
}

type serverData struct {
	Hostname string `json:"hostname"`
	IP       net.IP `json:"ip"`
	IP2      net.IP `json:"ip2"`
	IP3      net.IP `json:"ip3"`
}

func fetchAPI(ctx context.Context, client *http.Client) (
	data apiData, err error) {
	const baseURL = "https://assets.windscribe.com/serverlist/mob-v2/1/"
	cacheBreaker := time.Now().Unix()
	url := baseURL + strconv.Itoa(int(cacheBreaker))

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

	return data, nil
}
