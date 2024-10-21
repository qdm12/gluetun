package shadowsocks

import (
	"context"
	"fmt"
	"time"

	"github.com/qdm12/ss-server/pkg/tcpudp"
)

type service struct {
	// Injected settings
	settings tcpudp.Settings
	logger   Logger
	// Internal fields
	cancel context.CancelFunc
	done   <-chan struct{}
}

func newService(settings tcpudp.Settings,
	logger Logger) *service {
	return &service{
		settings: settings,
		logger:   logger,
	}
}

func (s *service) Start(ctx context.Context) (runError <-chan error, err error) {
	server, err := tcpudp.NewServer(s.settings, s.logger)
	if err != nil {
		return nil, fmt.Errorf("creating server: %w", err)
	}

	shadowsocksCtx, shadowsocksCancel := context.WithCancel(context.Background())
	s.cancel = shadowsocksCancel
	runErrorCh := make(chan error)
	done := make(chan struct{})
	s.done = done
	go func() {
		defer close(done)
		err = server.Listen(shadowsocksCtx)
		if shadowsocksCtx.Err() == nil {
			runErrorCh <- fmt.Errorf("listening: %w", err)
		}
	}()

	const minStabilityTime = 100 * time.Millisecond
	isStableTimer := time.NewTimer(minStabilityTime)
	select {
	case <-isStableTimer.C:
	case err = <-runErrorCh:
		return nil, fmt.Errorf("server became unstable within %s: %w",
			minStabilityTime, err)
	case <-ctx.Done():
		shadowsocksCancel()
		<-done
		return nil, ctx.Err()
	}

	return runErrorCh, nil
}

func (s *service) Stop() (err error) {
	s.cancel()
	<-s.done
	return nil
}
