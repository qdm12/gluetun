package storage

import "fmt"

func panicOnProviderMissingHardcoded(provider string) {
	panic(fmt.Sprintf("provider %s not found in hardcoded servers map; "+
		"did you add the provider key in the embedded servers.json?", provider))
}
