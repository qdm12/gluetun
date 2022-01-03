package cli

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/sources"
	"github.com/qdm12/gluetun/internal/healthcheck"
)

type HealthChecker interface {
	HealthCheck(ctx context.Context, source sources.Source, warner Warner) error
}

func (c *CLI) HealthCheck(ctx context.Context, source sources.Source, warner Warner) error {
	// Extract the health server port from the configuration.
	config, err := source.ReadHealth()
	if err != nil {
		return err
	}

	err = config.Validate()
	if err != nil {
		return err
	}

	_, port, err := net.SplitHostPort(config.ServerAddress)
	if err != nil {
		return err
	}

	const timeout = 10 * time.Second
	httpClient := &http.Client{Timeout: timeout}
	client := healthcheck.NewClient(httpClient)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	url := "http://127.0.0.1:" + port
	return client.Check(ctx, url)
}
