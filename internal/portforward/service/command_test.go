//go:build linux

package service

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/command"
	"github.com/stretchr/testify/require"
)

func Test_Service_runCommand(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	ctx := context.Background()
	cmder := command.New()
	const commandTemplate = `/bin/sh -c "echo {{PORTS}}"`
	ports := []uint16{1234, 5678}
	logger := NewMockLogger(ctrl)
	logger.EXPECT().Info("1234,5678")

	err := runCommand(ctx, cmder, logger, commandTemplate, ports)

	require.NoError(t, err)
}
