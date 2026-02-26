package updater

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrHTTPStatusCodeNotOK = errors.New("HTTP status code not OK")
	ErrDecodeResponseBody  = errors.New("failed decoding response body")
)

// nodeData represents one entry in the Cryptostorm node list JSON.
type nodeData struct {
	Hostname string `json:"hostname"`
	Country  string `json:"country"`
	City     string `json:"city"`
	IPv4     string `json:"ip"`
	WgPubKey string `json:"wg_pubkey"`
	PortFwd  bool   `json:"port_forward"`
}

// fetchAPI retrieves the Cryptostorm node list.
// Cryptostorm does not publish a formal JSON API; this function fetches
// their publicly available node list. If the upstream format changes,
// update the nodeData struct and parsing logic accordingly.
func fetchAPI(ctx context.Context, client *http.Client) (data []nodeData, err error) {
	const url = "https://cryptostorm.is/wireguard/nodes.json"

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
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDecodeResponseBody, err)
	}

	return data, nil
}
