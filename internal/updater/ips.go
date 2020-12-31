package updater

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"

	"github.com/qdm12/golibs/network"
)

func uniqueSortedIPs(ips []net.IP) []net.IP {
	uniqueIPs := make(map[string]struct{}, len(ips))
	for _, ip := range ips {
		key := ip.String()
		uniqueIPs[key] = struct{}{}
	}
	ips = make([]net.IP, 0, len(uniqueIPs))
	for key := range uniqueIPs {
		ip := net.ParseIP(key)
		if ipv4 := ip.To4(); ipv4 != nil {
			ip = ipv4
		}
		ips = append(ips, ip)
	}
	sort.Slice(ips, func(i, j int) bool {
		return bytes.Compare(ips[i], ips[j]) < 0
	})
	return ips
}

var errBadHTTPStatus = errors.New("bad HTTP status received")

type ipInfoData struct {
	Region  string `json:"region"`
	Country string `json:"country"`
	City    string `json:"city"`
}

func getIPInfo(ctx context.Context, client network.Client, ip net.IP) (country, region, city string, err error) {
	const baseURL = "https://ipinfo.io/"
	url := baseURL + ip.String()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", "", err
	}
	b, status, err := client.Do(request)
	if err != nil {
		return "", "", "", err
	} else if status != http.StatusOK {
		return "", "", "", fmt.Errorf("%w: %d", errBadHTTPStatus, status)
	}
	var data ipInfoData
	if err := json.Unmarshal(b, &data); err != nil {
		return "", "", "", err
	}
	country, ok := getCountryCodes()[strings.ToLower(data.Country)]
	if !ok {
		country = data.Country
	}
	return country, data.Region, data.City, nil
}
