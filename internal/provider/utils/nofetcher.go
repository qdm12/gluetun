package utils

import (
	"context"
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/models"
)

type NoFetcher struct {
	providerName string
}

func NewNoFetcher(providerName string) *NoFetcher {
	return &NoFetcher{
		providerName: providerName,
	}
}

var ErrFetcherNotSupported = errors.New("fetching of servers is not supported")

func (n *NoFetcher) FetchServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
	return nil, fmt.Errorf("%w: for %s", ErrFetcherNotSupported, n.providerName)
}
