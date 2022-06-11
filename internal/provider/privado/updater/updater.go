package updater

import (
	"net/http"

	"github.com/qdm12/gluetun/internal/provider/common"
)

type Updater struct {
	client    *http.Client
	unzipper  common.Unzipper
	presolver common.ParallelResolver
	warner    common.Warner
}

func New(client *http.Client, unzipper common.Unzipper,
	warner common.Warner) *Updater {
	return &Updater{
		client:    client,
		unzipper:  unzipper,
		presolver: newParallelResolver(),
		warner:    warner,
	}
}
