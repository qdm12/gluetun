package privateinternetaccess

import (
	"math/rand"
	"time"

	"github.com/qdm12/gluetun/internal/constants/openvpn"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/common"
)

type Provider struct {
	storage    common.Storage
	randSource rand.Source
	timeNow    func() time.Time
	// Port forwarding
	portForwardPath string
	authFilePath    string
}

func New(storage common.Storage, randSource rand.Source,
	timeNow func() time.Time) *Provider {
	const jsonPortForwardPath = "/gluetun/piaportforward.json"
	return &Provider{
		storage:         storage,
		timeNow:         timeNow,
		randSource:      randSource,
		portForwardPath: jsonPortForwardPath,
		authFilePath:    openvpn.AuthConf,
	}
}

func (p *Provider) Name() string {
	return providers.PrivateInternetAccess
}
