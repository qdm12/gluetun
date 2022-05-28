package windscribe

import (
	"net/http"
)

type Updater struct {
	client *http.Client
	warner Warner
}

type Warner interface {
	Warn(s string)
}

func New(client *http.Client, warner Warner) *Updater {
	return &Updater{
		client: client,
		warner: warner,
	}
}
