package cli

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/openvpn/extract"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/gluetun/internal/publicip/ipinfo"
	"github.com/qdm12/gluetun/internal/storage"
	"github.com/qdm12/gluetun/internal/updater/resolver"
)

type OpenvpnConfigLogger interface {
	Info(s string)
	Warn(s string)
}

type Unzipper interface {
	FetchAndExtract(ctx context.Context, url string) (
		contents map[string][]byte, err error)
}

type ParallelResolver interface {
	Resolve(ctx context.Context, settings resolver.ParallelSettings) (
		hostToIPs map[string][]net.IP, warnings []string, err error)
}

type IPFetcher interface {
	FetchMultiInfo(ctx context.Context, ips []net.IP) (data []ipinfo.Response, err error)
}

type IPv6Checker interface {
	IsIPv6Supported() (supported bool, err error)
}

func (c *CLI) OpenvpnConfig(logger OpenvpnConfigLogger, source Source,
	ipv6Checker IPv6Checker) error {
	storage, err := storage.New(logger, constants.ServersData)
	if err != nil {
		return err
	}

	allSettings, err := source.Read()
	if err != nil {
		return err
	}

	ipv6Supported, err := ipv6Checker.IsIPv6Supported()
	if err != nil {
		return fmt.Errorf("checking for IPv6 support: %w", err)
	}

	if err = allSettings.Validate(storage, ipv6Supported); err != nil {
		return fmt.Errorf("validating settings: %w", err)
	}

	// Unused by this CLI command
	unzipper := (Unzipper)(nil)
	client := (*http.Client)(nil)
	warner := (Warner)(nil)
	parallelResolver := (ParallelResolver)(nil)
	ipFetcher := (IPFetcher)(nil)
	openvpnFileExtractor := extract.New()

	providers := provider.NewProviders(storage, time.Now, warner, client,
		unzipper, parallelResolver, ipFetcher, openvpnFileExtractor)
	providerConf := providers.Get(*allSettings.VPN.Provider.Name)
	connection, err := providerConf.GetConnection(
		allSettings.VPN.Provider.ServerSelection, ipv6Supported)
	if err != nil {
		return err
	}

	lines := providerConf.OpenVPNConfig(connection,
		allSettings.VPN.OpenVPN, ipv6Supported)

	fmt.Println(strings.Join(lines, "\n"))
	return nil
}
