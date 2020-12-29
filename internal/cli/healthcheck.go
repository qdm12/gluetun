package cli

import (
	"context"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/healthcheck"
)

func (c *cli) HealthCheck(ctx context.Context) error {
	const timeout = 3 * time.Second
	httpClient := &http.Client{Timeout: timeout}
	healthchecker := healthcheck.NewChecker(httpClient)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	const url = "http://" + constants.HealthcheckAddress
	return healthchecker.Check(ctx, url)
}
