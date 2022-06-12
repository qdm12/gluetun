package publicip

import (
	"context"
	"net"

	"github.com/qdm12/gluetun/internal/publicip/ipinfo"
)

type Fetcher interface {
	FetchInfo(ctx context.Context, ip net.IP) (
		result ipinfo.Response, err error)
}
