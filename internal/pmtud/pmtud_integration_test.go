//go:build integration

package pmtud

import (
	"context"
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_PathMTUDiscover(t *testing.T) {
	t.Parallel()
	const physicalLinkMTU = 1500
	const timeout = time.Second
	mtu, err := PathMTUDiscover(context.Background(), netip.MustParseAddr("1.1.1.1"),
		physicalLinkMTU, timeout, nil)
	require.NoError(t, err)
	t.Log("MTU found:", mtu)
}
