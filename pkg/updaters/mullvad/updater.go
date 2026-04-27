package mullvad

import (
	"net/http"
)

type Updater struct {
	client *http.Client
}

func New(client *http.Client) *Updater {
	return &Updater{
		client: client,
	}
}

func (u *Updater) Version() uint16 {
	return 4 //nolint:mnd
}
