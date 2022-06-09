package privado

import (
	"net/http"

	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/updater/unzip"
)

type Updater struct {
	client    *http.Client
	unzipper  unzip.Unzipper
	presolver common.ParallelResolver
	warner    Warner
}

type Warner interface {
	Warn(s string)
}

func New(client *http.Client, unzipper unzip.Unzipper,
	warner Warner) *Updater {
	return &Updater{
		client:    client,
		unzipper:  unzipper,
		presolver: newParallelResolver(),
		warner:    warner,
	}
}
