package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Provider string

const (
	Cloudflare  Provider = "cloudflare"
	IfConfigCo  Provider = "ifconfigco"
	IPInfo      Provider = "ipinfo"
	IP2Location Provider = "ip2location"
)

type NameToken struct {
	Name  string
	Token string
}

func New(nameTokenPairs []NameToken, client *http.Client) (
	fetchers []Fetcher, err error,
) {
	fetchers = make([]Fetcher, len(nameTokenPairs))
	for i, nameTokenPair := range nameTokenPairs {
		provider, err := ParseProvider(nameTokenPair.Name)
		if err != nil {
			return nil, fmt.Errorf("parsing API name: %w", err)
		}
		switch provider {
		case Cloudflare:
			fetchers[i] = newCloudflare(client)
		case IfConfigCo:
			fetchers[i] = newIfConfigCo(client)
		case IPInfo:
			fetchers[i] = newIPInfo(client, nameTokenPair.Token)
		case IP2Location:
			fetchers[i] = newIP2Location(client, nameTokenPair.Token)
		default:
			panic("provider not valid: " + provider)
		}
	}
	return fetchers, nil
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
