//go:build integration

package healthcheck

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_healthCheck_ping(t *testing.T) {
	t.Parallel()

	const timeout = time.Second

	testCases := map[string]struct {
		address string
		err     error
	}{
		"github.com": {
			address: "github.com",
		},
		"99.99.99.99": {
			address: "99.99.99.99",
			err:     context.DeadlineExceeded,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			pinger := newPinger(testCase.address)

			err := healthCheck(ctx, pinger)

			assert.ErrorIs(t, testCase.err, err)
		})
	}
}
