package updater

import (
	"net/http"

	"github.com/qdm12/gluetun/internal/provider/common"
)

type Updater struct {
	client   *http.Client
	username string
	password string
	warner   common.Warner
}

func New(client *http.Client, warner common.Warner, username, password string) *Updater {
	return &Updater{
		client:   client,
		username: username,
		password: password,
		warner:   warner,
	}
}
