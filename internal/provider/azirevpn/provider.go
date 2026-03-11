package azirevpn

import (
	"math/rand"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/azirevpn/updater"
	"github.com/qdm12/gluetun/internal/provider/common"
)

type Provider struct {
	storage    common.Storage
	randSource rand.Source
	common.Fetcher

	client *http.Client
	token  string

	dataPath string
}

func New(storage common.Storage, randSource rand.Source, client *http.Client,
	updaterWarner common.Warner, token string,
) *Provider {
	const jsonDataPath = "/tmp/gluetun/azirevpn_data.json"
	return &Provider{
		storage:    storage,
		randSource: randSource,
		Fetcher:    updater.New(client, updaterWarner, token),
		client:     client,
		token:      token,
		dataPath:   jsonDataPath,
	}
}

func (p *Provider) Name() string {
	return providers.Azirevpn
}
