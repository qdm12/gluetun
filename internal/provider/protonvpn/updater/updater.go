package updater

import (
	"net/http"

	"github.com/qdm12/gluetun/internal/provider/common"
)

type Updater struct {
	client   *http.Client
	email    string
	password string
	warner   common.Warner
}

func New(client *http.Client, warner common.Warner, email, password string) *Updater {
	return &Updater{
		client:   client,
		email:    email,
		password: password,
		warner:   warner,
	}
}
