package common

import (
	"context"
	"errors"
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

var ErrNotEnoughServers = errors.New("not enough servers found")

type Fetcher interface {
	FetchServers(ctx context.Context, minServers int) (servers []models.Server, err error)
}

type ParallelResolver interface {
	Resolve(ctx context.Context, hosts []string, minToFind int) (
		hostToIPs map[string][]net.IP, warnings []string, err error)
}

type Unzipper interface {
	FetchAndExtract(ctx context.Context, url string) (contents map[string][]byte, err error)
}

type Warner interface {
	Warn(s string)
}
