package healthcheck

import (
	"context"
	"errors"
	"net/http"
	"time"
)

func (s *Server) Run(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	loopDone := make(chan struct{})
	go s.runHealthcheckLoop(ctx, loopDone)

	server := http.Server{
		Addr:              s.config.ServerAddress,
		Handler:           s.handler,
		ReadHeaderTimeout: s.config.ReadHeaderTimeout,
		ReadTimeout:       s.config.ReadTimeout,
	}
	serverDone := make(chan struct{})
	go func() {
		defer close(serverDone)
		<-ctx.Done()
		const shutdownGraceDuration = 2 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGraceDuration)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("failed shutting down: " + err.Error())
		}
	}()

	s.logger.Info("listening on " + s.config.ServerAddress)
	err := server.ListenAndServe()
	if err != nil && !errors.Is(ctx.Err(), context.Canceled) {
		s.logger.Error(err.Error())
	}

	<-loopDone
	<-serverDone
}
