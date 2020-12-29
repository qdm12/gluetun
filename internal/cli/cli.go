package cli

import (
	"context"

	"github.com/qdm12/gluetun/internal/os"
)

type CLI interface {
	ClientKey(args []string, openFile os.OpenFileFunc) error
	HealthCheck(ctx context.Context) error
	OpenvpnConfig(os os.OS) error
	Update(args []string, os os.OS) error
}

type cli struct{}

func New() CLI {
	return &cli{}
}
