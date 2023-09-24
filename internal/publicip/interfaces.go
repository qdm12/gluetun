package publicip

import (
	"context"
	"net/netip"

	"github.com/qdm12/gluetun/internal/publicip/ipinfo"
)

type Fetcher interface {
	FetchInfo(ctx context.Context, ip netip.Addr) (
		result ipinfo.Response, err error)
}

type Logger interface {
	Info(s string)
	Warn(s string)
	Error(s string)
}
