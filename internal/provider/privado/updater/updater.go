package privado

import (
	"net/http"

	"github.com/qdm12/gluetun/internal/updater/resolver"
	"github.com/qdm12/gluetun/internal/updater/unzip"
)

type Updater struct {
	client    *http.Client
	unzipper  unzip.Unzipper
	presolver resolver.Parallel
	warner    Warner
}

type Warner interface {
	Warn(s string)
}

func New(client *http.Client, unzipper unzip.Unzipper,
	presolver resolver.Parallel, warner Warner) *Updater {
	return &Updater{
		client:    client,
		unzipper:  unzipper,
		presolver: presolver,
		warner:    warner,
	}
}
