package wevpn

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/updater/resolver"
	"github.com/qdm12/gluetun/internal/updater/resolver/mock_resolver"
	"github.com/stretchr/testify/assert"
)

func Test_resolveHosts(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	ctx := context.Background()
	presolver := mock_resolver.NewMockParallel(ctrl)
	hosts := []string{"host1", "host2"}
	const minServers = 10

	expectedHostToIPs := map[string][]net.IP{
		"host1": {{1, 2, 3, 4}},
		"host2": {{2, 3, 4, 5}},
	}
	expectedWarnings := []string{"warning1", "warning2"}
	expectedErr := errors.New("dummy")

	const (
		maxFailRatio    = 0.1
		maxDuration     = 20 * time.Second
		betweenDuration = time.Second
		maxNoNew        = 2
		maxFails        = 2
	)
	expectedSettings := resolver.ParallelSettings{
		MaxFailRatio: maxFailRatio,
		MinFound:     minServers,
		Repeat: resolver.RepeatSettings{
			MaxDuration:     maxDuration,
			BetweenDuration: betweenDuration,
			MaxNoNew:        maxNoNew,
			MaxFails:        maxFails,
			SortIPs:         true,
		},
	}
	presolver.EXPECT().Resolve(ctx, hosts, expectedSettings).
		Return(expectedHostToIPs, expectedWarnings, expectedErr)

	hostToIPs, warnings, err := resolveHosts(ctx, presolver, hosts, minServers)
	assert.Equal(t, expectedHostToIPs, hostToIPs)
	assert.Equal(t, expectedWarnings, warnings)
	assert.Equal(t, expectedErr, err)
}
