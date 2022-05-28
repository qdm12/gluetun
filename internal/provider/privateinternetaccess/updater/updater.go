package privateinternetaccess

import (
	"net/http"
)

type Updater struct {
	client *http.Client
}

type Warner interface {
	Warn(s string)
}

func New(client *http.Client) *Updater {
	return &Updater{
		client: client,
	}
}
