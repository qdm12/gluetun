package api

import (
	"context"
	"net/netip"

	"github.com/qdm12/gluetun/internal/models"
)

type Fetcher interface {
	FetchInfo(ctx context.Context, ip netip.Addr) (
		result models.PublicIP, err error)
}
