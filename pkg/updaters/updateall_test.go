package updaters

import (
	"net/http"
	"testing"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/stretchr/testify/assert"
)

func Test_buildFetchers(t *testing.T) {
	t.Parallel()

	// Verify that all providers have a fetcher
	client := (*http.Client)(nil)
	logger := (Logger)(nil)
	parallelResolver := (common.ParallelResolver)(nil)
	unzipper := (common.Unzipper)(nil)
	ipFetcher := (common.IPFetcher)(nil)

	fetchers := buildFetchers(client, parallelResolver, unzipper, ipFetcher,
		logger, "protonEmail", "protonPassword")

	allProviders := providers.All()
	assert.Equal(t, len(allProviders), len(fetchers), "number of fetchers should match number of providers")
	for _, providerName := range allProviders {
		assert.Contains(t, fetchers, providerName, "fetcher for provider %s should be present", providerName)
	}
}
