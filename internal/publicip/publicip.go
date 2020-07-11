package publicip

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strings"

	"github.com/qdm12/golibs/network"
)

type IPGetter interface {
	Get() (ip net.IP, err error)
}

type ipGetter struct {
	client   network.Client
	randIntn func(n int) int
}

func NewIPGetter(client network.Client) IPGetter {
	return &ipGetter{
		client:   client,
		randIntn: rand.Intn,
	}
}

func (i *ipGetter) Get() (ip net.IP, err error) {
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
	content, status, err := i.client.GetContent(url, network.UseRandomUserAgent())
	if err != nil {
		return nil, err
	} else if status != http.StatusOK {
		return nil, fmt.Errorf("received unexpected status code %d from %s", status, url)
	}
	s := strings.ReplaceAll(string(content), "\n", "")
	ip = net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("cannot parse IP address from %q", s)
	}
	return ip, nil
}
