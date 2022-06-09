package privateinternetaccess

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/constants/openvpn"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/provider/privateinternetaccess/updater"
)

type Provider struct {
	storage    common.Storage
	randSource rand.Source
	timeNow    func() time.Time
	common.Fetcher
	// Port forwarding
	portForwardPath string
	authFilePath    string
}

func New(storage common.Storage, randSource rand.Source,
	timeNow func() time.Time, client *http.Client) *Provider {
	const jsonPortForwardPath = "/gluetun/piaportforward.json"
	return &Provider{
		storage:         storage,
		timeNow:         timeNow,
		randSource:      randSource,
		portForwardPath: jsonPortForwardPath,
		authFilePath:    openvpn.AuthConf,
		Fetcher:         updater.New(client),
	}
}

func (p *Provider) Name() string {
	return providers.PrivateInternetAccess
}
