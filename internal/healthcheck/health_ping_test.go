//go:build integration

package healthcheck

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_healthCheck_ping(t *testing.T) {
	t.Parallel()

	pinger := newPinger()

	err := healthCheck(context.Background(), pinger)

	assert.NoError(t, err)
}
