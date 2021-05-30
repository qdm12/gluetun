// Package publicip defines interfaces to get your public IP address.
package publicip

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strings"
)

type IPGetter interface {
	Get(ctx context.Context) (ip net.IP, err error)
}

type ipGetter struct {
	client   *http.Client
	randIntn func(n int) int
}

func NewIPGetter(client *http.Client) IPGetter {
	return &ipGetter{
		client:   client,
		randIntn: rand.Intn,
	}
}

var ErrParseIP = errors.New("cannot parse IP address")

func (i *ipGetter) Get(ctx context.Context) (ip net.IP, err error) {
	urls := []string{
		"https://ifconfig.me/ip",
		"http://ip1.dynupdate.no-ip.com:8245",
		"http://ip1.dynupdate.no-ip.com",
		"https://api.ipify.org",
		"https://diagnostic.opendns.com/myip",
		"https://domains.google.com/checkip",
		"https://ifconfig.io/ip",
		"https://ipinfo.io/ip",
	}
	url := urls[i.randIntn(len(urls))]

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := i.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w from %s: %s", ErrBadStatusCode, url, response.Status)
	}

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCannotReadBody, err)
	}

	s := strings.ReplaceAll(string(content), "\n", "")
	ip = net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("%w: %s", ErrParseIP, s)
	}
	return ip, nil
}
