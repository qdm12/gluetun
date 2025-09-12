package cli

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/healthcheck"
	"github.com/qdm12/gosettings/reader"
)

type HealthCheckCommand struct {
	reader *reader.Reader
}

func NewHealthCheckCommand(reader *reader.Reader) *HealthCheckCommand {
	return &HealthCheckCommand{
		reader: reader,
	}
}

func (c *HealthCheckCommand) Name() string {
	return "healthcheck"
}

func (c *HealthCheckCommand) Description() string {
	return "Check the health of the VPN connection of another Gluetun instance"
}

func (c *HealthCheckCommand) Run(ctx context.Context) (err error) {
	// Extract the health server port from the configuration.
	var config settings.Health
	err = config.Read(c.reader)
	if err != nil {
		return err
	}

	config.SetDefaults()

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
