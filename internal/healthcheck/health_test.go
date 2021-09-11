package healthcheck

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_healthCheck(t *testing.T) {
	t.Parallel()

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	someErr := errors.New("error")

	testCases := map[string]struct {
		ctx      context.Context
		runErr   error
		stopCall bool
		err      error
	}{
		"success": {
			ctx: context.Background(),
		},
		"error": {
			ctx:    context.Background(),
			runErr: someErr,
			err:    someErr,
		},
		"context canceled": {
			ctx:      canceledCtx,
			stopCall: true,
			err:      context.Canceled,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			stopped := make(chan struct{})

			pinger := NewMockPinger(ctrl)
			pinger.EXPECT().Run().DoAndReturn(func() error {
				if testCase.stopCall {
					<-stopped
				}
				return testCase.runErr
			})

			if testCase.stopCall {
				pinger.EXPECT().Stop().DoAndReturn(func() {
					close(stopped)
				})
			}

			err := healthCheck(testCase.ctx, pinger)

			assert.ErrorIs(t, testCase.err, err)
		})
	}

	t.Run("canceled real pinger", func(t *testing.T) {
		t.Parallel()

		pinger := newPinger()

		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		err := healthCheck(canceledCtx, pinger)

		assert.ErrorIs(t, context.Canceled, err)
	})
}
