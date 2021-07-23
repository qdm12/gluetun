package cli

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/healthcheck"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/params"
)

func (c *cli) HealthCheck(ctx context.Context, env params.Env,
	logger logging.Logger) error {
	// Extract the health server port from the configuration.
	config := configuration.Health{}
	err := config.Read(env, logger)
	if err != nil {
		return err
	}
	_, port, err := net.SplitHostPort(config.ServerAddress)
	if err != nil {
		return err
	}

	const timeout = 10 * time.Second
	httpClient := &http.Client{Timeout: timeout}
	healthchecker := healthcheck.NewChecker(httpClient)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	url := "http://127.0.0.1:" + port
	return healthchecker.Check(ctx, url)
}
