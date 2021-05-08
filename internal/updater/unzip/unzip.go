// Package unzip defines the Unzipper which fetches and extract a zip file
// containing multiple files.
package unzip

import (
	"context"
	"net/http"
)

type Unzipper interface {
	FetchAndExtract(ctx context.Context, url string) (contents map[string][]byte, err error)
}

type unzipper struct {
	client *http.Client
}

func New(client *http.Client) Unzipper {
	return &unzipper{
		client: client,
	}
}
