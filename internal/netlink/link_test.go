//go:build linux

package netlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NetLink_LinkList(t *testing.T) {
	t.Parallel()

	netlink := &NetLink{}

	initialLinks, err := netlink.LinkList()
	require.NoError(t, err)
	require.NotEmpty(t, initialLinks)

	loopbackFound := false
	for _, link := range initialLinks {
		if link.Name != "lo" {
			continue
		}
		loopbackFound = true
		assert.Equal(t, DeviceTypeLoopback, link.DeviceType)
		break
	}
	assert.True(t, loopbackFound, "loopback interface not found")

	testLink := Link{
		Name: makeLinkName(),
		// note if [Link.VirtualType] is set, [Link.DeviceType]
		// is ignored and gets set to [DeviceTypeNone] in LinkAdd.
		DeviceType:  DeviceTypeNone,
		VirtualType: "wireguard",
		MTU:         1420,
	}
	index, err := netlink.LinkAdd(testLink)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = netlink.LinkDel(index)
	})

	links, err := netlink.LinkList()
	require.NoError(t, err)

	testLink.Index = index
	for _, link := range links {
		if link.Name != testLink.Name {
			continue
		}
		assert.Equal(t, testLink, link)
		return
	}
	t.Errorf("created link %q not found", testLink.Name)
}

func Test_NetLink_LinkSetMTU(t *testing.T) {
	t.Parallel()

	netlink := &NetLink{}

	testLink := Link{
		Name:        makeLinkName(),
		DeviceType:  DeviceTypeNone,
		VirtualType: "wireguard",
		MTU:         1420,
	}
	index, err := netlink.LinkAdd(testLink)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = netlink.LinkDel(index)
	})
	testLink.Index = index

	err = netlink.LinkSetMTU(index, 1500)
	require.NoError(t, err)

	link, err := netlink.LinkByIndex(index)
	require.NoError(t, err)
	testLink.MTU = 1500
	assert.Equal(t, testLink, link)
}
