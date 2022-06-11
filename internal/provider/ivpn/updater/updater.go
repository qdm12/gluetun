package updater

import (
	"net/http"

	"github.com/qdm12/gluetun/internal/provider/common"
)

type Updater struct {
	client    *http.Client
	presolver common.ParallelResolver
	warner    common.Warner
}

func New(client *http.Client, warner common.Warner) *Updater {
	return &Updater{
		client:    client,
		presolver: newParallelResolver(),
		warner:    warner,
	}
}
