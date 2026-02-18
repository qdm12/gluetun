//go:build integration

package tcp

import (
	"errors"
	"net/netip"
	"testing"
	"time"

	"github.com/qdm12/gluetun/internal/command"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/pmtud/constants"
	"github.com/qdm12/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PathMTUDiscover(t *testing.T) {
	t.Parallel()
	noopLogger := log.New(log.SetLevel(log.LevelDebug))

	cmder := command.New()
	fw, err := firewall.NewConfig(t.Context(), noopLogger, cmder, nil, nil)
	if errors.Is(err, firewall.ErrIPTablesNotSupported) {
		t.Skip("iptables not installed, skipping TCP PMTUD tests")
	}
	require.NoError(t, err, "creating firewall config")

	dst := netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 1, 1, 1}), 80)
	const minMTU = constants.MinIPv6MTU
	const maxMTU = constants.MaxEthernetFrameSize
	const tryTimeout = time.Second
	mtu, err := PathMTUDiscover(t.Context(), dst, minMTU, maxMTU, tryTimeout, fw, noopLogger)
	require.NoError(t, err, "discovering path MTU")
	assert.Greater(t, mtu, uint32(0), "MTU should be greater than 0")
	t.Logf("discovered path MTU to %s is %d", dst, mtu)
}
