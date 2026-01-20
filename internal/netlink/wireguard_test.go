//go:build linux || darwin

package netlink

import (
	"testing"
)

func Test_NetLink_IsWireguardSupported(t *testing.T) {
	t.Parallel()

	netLink := &NetLink{
		debugLogger: &noopLogger{},
	}
	ok := netLink.IsWireguardSupported()
	if ok { // cannot assert since this depends on kernel
		t.Log("wireguard is supported")
	} else {
		t.Log("wireguard is not supported")
	}
}
