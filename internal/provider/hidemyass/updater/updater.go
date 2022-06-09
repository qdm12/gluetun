package updater

import (
	"net/http"

	"github.com/qdm12/gluetun/internal/provider/common"
)

type Updater struct {
	client    *http.Client
	presolver common.ParallelResolver
	warner    Warner
}

type Warner interface {
	Warn(s string)
}

func New(client *http.Client, warner Warner) *Updater {
	return &Updater{
		client:    client,
		presolver: newParallelResolver(),
		warner:    warner,
	}
}
