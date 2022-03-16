package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseHardcodedServers(t *testing.T) {
	t.Parallel()

	servers, err := parseHardcodedServers()

	require.NoError(t, err)
	require.NotEmpty(t, len(servers.Cyberghost.Servers))
}
