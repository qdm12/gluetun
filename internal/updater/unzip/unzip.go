package unzip

import (
	"net/http"
)

type Unzipper struct {
	client *http.Client
}

func New(client *http.Client) *Unzipper {
	return &Unzipper{
		client: client,
	}
}
