// Package unzip defines the Unzipper which fetches and extract a zip file
// containing multiple files.
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
