package cryptostorm

import (
	"math/rand"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/provider/cryptostorm/updater"
)

type Provider struct {
	storage         common.Storage
	randSource      rand.Source
	portForwardPath string
	forwardedPort   uint16 // set after a successful PortForward, used for teardown
	common.Fetcher
}

func New(storage common.Storage, randSource rand.Source,
	client *http.Client, updaterWarner common.Warner,
	parallelResolver common.ParallelResolver,
) *Provider {
	const jsonPortForwardPath = "/gluetun/portforward/cryptostorm.json"
	return &Provider{
		storage:         storage,
		randSource:      randSource,
		portForwardPath: jsonPortForwardPath,
		Fetcher:         updater.New(client, updaterWarner, parallelResolver),
	}
}

func (p *Provider) Name() string {
	return providers.Cryptostorm
}
