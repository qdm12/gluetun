package netlink

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NetLink_IsWireguardSupported(t *testing.T) {
	t.Skip() // TODO unskip once the data race problem with netlink.GenlFamilyList() is fixed

	t.Parallel()
	netLink := &NetLink{}
	ok, err := netLink.IsWireguardSupported()
	require.NoError(t, err)
	if ok { // cannot assert since this depends on kernel
		t.Log("wireguard is supported")
	} else {
		t.Log("wireguard is not supported")
	}
}
