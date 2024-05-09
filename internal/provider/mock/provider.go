package mock

import (
	"context"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
	"github.com/qdm12/log"
)

type Provider struct {
	logger *log.Logger
}

func New() *Provider {
	logger := log.New(log.SetComponent("mock provider"),
		log.SetCallerFile(true), log.SetCallerLine(true))
	return &Provider{
		logger: logger,
	}
}

func (p *Provider) Name() string {
	return providers.Mock
}

func (p *Provider) GetConnection(_ settings.ServerSelection, _ bool) (
	connection models.Connection, err error) {
	p.logger.Info("getting connection")
	return models.Connection{}, nil
}

func (p *Provider) OpenVPNConfig(_ models.Connection,
	_ settings.OpenVPN, _ bool) (lines []string) {
	p.logger.Info("generating openvpn config")
	return nil
}

func (p *Provider) FetchServers(_ context.Context, _ int) (
	servers []models.Server, err error) {
	p.logger.Info("fetching servers")
	return servers, nil
}

func (p *Provider) PortForward(_ context.Context,
	_ utils.PortForwardObjects) (port uint16, err error) {
	const mockPort = 12345
	p.logger.Info("port forward")
	return mockPort, nil
}

func (p *Provider) KeepPortForward(ctx context.Context,
	_ utils.PortForwardObjects) (err error) {
	p.logger.Info("keeping port forward start")
	defer p.logger.Info("keeping port forward exited")
	const keepAlivePeriod = 10 * time.Second
	keepAliveTimer := time.NewTimer(keepAlivePeriod)
	for {
		select {
		case <-ctx.Done():
			p.logger.Info("keeping port forward context canceled")
			if !keepAliveTimer.Stop() {
				<-keepAliveTimer.C
			}
			return ctx.Err()
		case <-keepAliveTimer.C:
			p.logger.Info("keeping port forward ticked")
			keepAliveTimer.Reset(keepAlivePeriod)
		}
	}
}
