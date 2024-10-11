package publicip

import (
	"context"
	"net/netip"

	"github.com/qdm12/gluetun/internal/models"
)

type Fetcher interface {
	String() string
	FetchInfo(ctx context.Context, ip netip.Addr) (
		result models.PublicIP, err error)
}

type Logger interface {
	Info(s string)
	Warn(s string)
	Error(s string)
}
