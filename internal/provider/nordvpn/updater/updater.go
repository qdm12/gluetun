package nordvpn

import (
	"net/http"
)

type Updater struct {
	client *http.Client
	warner Warner
}

type Warner interface {
	Warn(message string)
}

func New(client *http.Client, warner Warner) *Updater {
	return &Updater{
		client: client,
		warner: warner,
	}
}
