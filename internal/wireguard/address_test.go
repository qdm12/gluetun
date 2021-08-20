package wireguard

import (
	"fmt"
	"math/rand"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vishvananda/netlink"
)

func Test_addAddresses(t *testing.T) {
	t.Parallel()

	intfName := "test_" + fmt.Sprint(rand.Intn(10000)) //nolint:gosec

	// Add link
	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = intfName
	link := &netlink.Bridge{
		LinkAttrs: linkAttrs,
	}
	err := netlink.LinkAdd(link)
	require.NoError(t, err)

	defer func() {
		err = netlink.LinkDel(link)
		assert.NoError(t, err)
	}()

	addresses := []*net.IPNet{
		{IP: net.IP{1, 2, 3, 4}, Mask: net.IPv4Mask(255, 255, 255, 255)},
		{IP: net.IP{5, 6, 7, 8}, Mask: net.IPv4Mask(255, 255, 255, 255)},
	}

	// Success
	err = addAddresses(link, addresses)
	require.NoError(t, err)

	netlinkAddresses, err := netlink.AddrList(link, netlink.FAMILY_ALL)
	require.NoError(t, err)
	require.Equal(t, len(addresses), len(netlinkAddresses))
	for i, netlinkAddress := range netlinkAddresses {
		ipNet := netlinkAddress.IPNet
		assert.Equal(t, addresses[i], ipNet)
	}

	// Existing address cannot be added
	err = addAddresses(link, addresses)
	require.Error(t, err)
	assert.Equal(t, "file exists: when adding address 1.2.3.4/32 to link test_8081", err.Error())
}
