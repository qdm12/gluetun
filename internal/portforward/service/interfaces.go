package service

import (
	"context"
)

type PortAllower interface {
	SetAllowedPort(ctx context.Context, port uint16, intf string) (err error)
	RemoveAllowedPort(ctx context.Context, port uint16) (err error)
}

type Logger interface {
	Info(s string)
	Warn(s string)
	Error(s string)
}
