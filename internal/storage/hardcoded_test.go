package storage

import (
	"testing"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseHardcodedServers(t *testing.T) {
	t.Parallel()

	servers, err := parseHardcodedServers()

	require.NoError(t, err)

	// all providers minus custom
	allProviders := providers.All()
	require.Equal(t, len(allProviders), len(servers.ProviderToServers))
	for _, provider := range allProviders {
		servers, ok := servers.ProviderToServers[provider]
		assert.Truef(t, ok, "for provider %s", provider)
		assert.NotEmptyf(t, servers, "for provider %s", provider)
	}
}
