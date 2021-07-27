package windscribe

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
	ErrHTTPStatusCodeNotOK   = errors.New("HTTP status code not OK")
	ErrUnmarshalResponseBody = errors.New("failed unmarshaling response body")
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

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&data); err != nil {
		return data, fmt.Errorf("%w: %s", ErrUnmarshalResponseBody, err)
	}

	return data, nil
}
