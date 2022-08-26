package files

import (
	"net/netip"
	"testing"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func uint16Ptr(n uint16) *uint16 { return &n }

func Test_Source_readWireguardSelection_integration(t *testing.T) {
	t.Parallel()

	source := &Source{
		wireguardConfigPath: "./testdata/wg.conf",
	}

	wireguardSelection, err := source.readWireguardSelection()
	require.NoError(t, err)

	expectedWireguardSelection := settings.WireguardSelection{
		PublicKey:    "EIORcksHjtrjwP6uVveZLyR/00GSp3xlIkwB3JXikEM=",
		EndpointIP:   netip.AddrFrom4([4]byte{193, 32, 249, 66}),
		EndpointPort: uint16Ptr(51820),
	}
	assert.Equal(t, expectedWireguardSelection, wireguardSelection)
}
