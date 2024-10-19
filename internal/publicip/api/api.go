package api

import (
	"errors"
	"fmt"
	"maps"
	"net/http"
	"slices"
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
	possibleProviders := []Provider{
		Cloudflare,
		IfConfigCo,
		IP2Location,
		IPInfo,
	}
	stringToProvider := make(map[string]Provider, len(possibleProviders))
	for _, provider := range possibleProviders {
		stringToProvider[string(provider)] = provider
	}
	provider, ok := stringToProvider[strings.ToLower(s)]
	if ok {
		return provider, nil
	}

	providerStrings := slices.Sorted(maps.Keys(stringToProvider))
	for i := range providerStrings {
		providerStrings[i] = `"` + providerStrings[i] + `"`
	}

	return "", fmt.Errorf(`%w: %q can only be %s`,
		ErrProviderNotValid, s, orStrings(providerStrings))
}

func orStrings(strings []string) (result string) {
	return joinStrings(strings, "or")
}

func joinStrings(strings []string, lastJoin string) (result string) {
	if len(strings) == 0 {
		return ""
	}

	result = strings[0]
	for i := 1; i < len(strings); i++ {
		if i < len(strings)-1 {
			result += ", " + strings[i]
		} else {
			result += " " + lastJoin + " " + strings[i]
		}
	}

	return result
}
