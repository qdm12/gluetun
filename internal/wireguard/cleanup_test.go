package wireguard

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_closers(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	var ACloseCalled, BCloseCalled, CCloseCalled bool
	var (
		AErr error
		BErr = errors.New("B failed")
		CErr = errors.New("C failed")
	)

	var closers closers
	closers.add("closing A", stepFive, func() error {
		ACloseCalled = true
		return AErr
	})

	closers.add("closing B", stepThree, func() error {
		BCloseCalled = true
		return BErr
	})

	closers.add("closing C", stepTwo, func() error {
		CCloseCalled = true
		return CErr
	})

	logger := NewMockLogger(ctrl)
	prevCall := logger.EXPECT().Debug("closing C...")
	prevCall = logger.EXPECT().Error("failed closing C: C failed").After(prevCall)
	prevCall = logger.EXPECT().Debug("closing B...").After(prevCall)
	prevCall = logger.EXPECT().Error("failed closing B: B failed").After(prevCall)
	logger.EXPECT().Debug("closing A...").After(prevCall)

	closers.cleanup(logger)

	closers.cleanup(logger) // run twice should not close already closed

	for _, closer := range closers {
		assert.True(t, closer.closed)
	}

	assert.True(t, ACloseCalled)
	assert.True(t, BCloseCalled)
	assert.True(t, CCloseCalled)
}
