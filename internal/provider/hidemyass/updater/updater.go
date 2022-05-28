package hidemyass

import (
	"net/http"

	"github.com/qdm12/gluetun/internal/updater/resolver"
)

type Updater struct {
	client    *http.Client
	presolver resolver.Parallel
	warner    Warner
}

type Warner interface {
	Warn(s string)
}

func New(client *http.Client, presolver resolver.Parallel,
	warner Warner) *Updater {
	return &Updater{
		client:    client,
		presolver: presolver,
		warner:    warner,
	}
}
