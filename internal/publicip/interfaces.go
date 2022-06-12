package publicip

import (
	"context"
	"net"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/publicip/ipinfo"
)

type statusManager interface {
	GetStatus() (status models.LoopStatus)
	SetStatus(status models.LoopStatus)
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
}

type stateManager interface {
	GetData() (data ipinfo.Response)
	SetData(data ipinfo.Response)
	GetSettings() (settings settings.PublicIP)
	SetSettings(ctx context.Context, settings settings.PublicIP) (outcome string)
}

type Fetcher interface {
	FetchInfo(ctx context.Context, ip net.IP) (
		result ipinfo.Response, err error)
}
