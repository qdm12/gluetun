package publicip

import (
	"context"
	"net"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	publicipmodels "github.com/qdm12/gluetun/internal/publicip/models"
)

type statusManager interface {
	GetStatus() (status models.LoopStatus)
	SetStatus(status models.LoopStatus)
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
}

type stateManager interface {
	GetData() (data publicipmodels.IPInfoData)
	SetData(data publicipmodels.IPInfoData)
	GetSettings() (settings settings.PublicIP)
	SetSettings(ctx context.Context, settings settings.PublicIP) (outcome string)
}

type fetcher interface {
	FetchPublicIP(ctx context.Context) (ip net.IP, err error)
}
