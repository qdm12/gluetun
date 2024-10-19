package api

import (
	"context"
	"net/netip"

	"github.com/qdm12/gluetun/internal/models"
)

type Fetcher interface {
	String() string
	CanFetchAnyIP() bool
	Token() (token string)
	InfoFetcher
}

type InfoFetcher interface {
	FetchInfo(ctx context.Context, ip netip.Addr) (
		result models.PublicIP, err error)
}

type Warner interface {
	Warn(message string)
}
