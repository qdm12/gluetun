package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
)

type API interface {
	FetchInfo(ctx context.Context, ip netip.Addr) (
		result models.PublicIP, err error)
}

type Provider string

const (
	IPInfo Provider = "ipinfo"
)

func New(provider Provider, client *http.Client, token string) ( //nolint:ireturn
	a API, err error) {
	switch provider {
	case IPInfo:
		return newIPInfo(client, token), nil
	default:
		panic("provider not valid: " + provider)
	}
}

var (
	ErrProviderNotValid = errors.New("API name is not valid")
)

func ParseProvider(s string) (provider Provider, err error) {
	switch strings.ToLower(s) {
	case "ipinfo":
		return IPInfo, nil
	default:
		return "", fmt.Errorf(`%w: %q can only be "ipinfo" or "ip2location"`,
			ErrProviderNotValid, s)
	}
}
