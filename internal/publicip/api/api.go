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
	String() string
	CanFetchAnyIP() bool
	FetchInfo(ctx context.Context, ip netip.Addr) (
		result models.PublicIP, err error)
}

type Provider string

const (
	Cloudflare  Provider = "cloudflare"
	IfConfigCo  Provider = "ifconfigco"
	IPInfo      Provider = "ipinfo"
	IP2Location Provider = "ip2location"
)

func New(provider Provider, client *http.Client, token string) ( //nolint:ireturn
	a API, err error,
) {
	switch provider {
	case Cloudflare:
		return newCloudflare(client), nil
	case IfConfigCo:
		return newIfConfigCo(client), nil
	case IPInfo:
		return newIPInfo(client, token), nil
	case IP2Location:
		return newIP2Location(client, token), nil
	default:
		panic("provider not valid: " + provider)
	}
}

var ErrProviderNotValid = errors.New("API name is not valid")

func ParseProvider(s string) (provider Provider, err error) {
	switch strings.ToLower(s) {
	case "cloudflare":
		return Cloudflare, nil
	case string(IfConfigCo):
		return IfConfigCo, nil
	case "ipinfo":
		return IPInfo, nil
	case "ip2location":
		return IP2Location, nil
	default:
		return "", fmt.Errorf(`%w: %q can only be "cloudflare", "ifconfigco", "ip2location" or "ipinfo"`,
			ErrProviderNotValid, s)
	}
}
