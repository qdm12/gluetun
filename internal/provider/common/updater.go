package common

import (
	"context"
	"errors"
	"net"
)

var ErrNotEnoughServers = errors.New("not enough servers found")

type ParallelResolver interface {
	Resolve(ctx context.Context, hosts []string, minToFind int) (
		hostToIPs map[string][]net.IP, warnings []string, err error)
}
