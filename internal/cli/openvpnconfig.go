package cli

import (
	"context"
	"fmt"
	"net/http"
	"net/netip"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/qdm12/gluetun/internal/openvpn/extract"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/gluetun/internal/storage"
	"github.com/qdm12/gluetun/internal/updater/resolver"
	"github.com/qdm12/gosettings/reader"
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
		hostToIPs map[string][]netip.Addr, warnings []string, err error)
}

type IPFetcher interface {
	String() string
	CanFetchAnyIP() bool
	FetchInfo(ctx context.Context, ip netip.Addr) (data models.PublicIP, err error)
}

type IPv6Checker interface {
	FindIPv6SupportLevel(ctx context.Context,
		checkAddresses []netip.AddrPort, firewall netlink.Firewall,
	) (level netlink.IPv6SupportLevel, err error)
}

type OpenVPNConfigCommand struct {
	logger      OpenvpnConfigLogger
	reader      *reader.Reader
	ipv6Checker IPv6Checker
}

func NewOpenVPNConfigCommand(logger OpenvpnConfigLogger, reader *reader.Reader,
	ipv6Checker IPv6Checker,
) *OpenVPNConfigCommand {
	return &OpenVPNConfigCommand{
		logger:      logger,
		reader:      reader,
		ipv6Checker: ipv6Checker,
	}
}

func (c *OpenVPNConfigCommand) Name() string {
	return "openvpnconfig"
}

func (c *OpenVPNConfigCommand) Description() string {
	return "OPrint the OpenVPN configuration (for debugging)"
}

func (c *OpenVPNConfigCommand) Run(ctx context.Context) error {
	storage, err := storage.New(c.logger, constants.ServersData)
	if err != nil {
		return err
	}

	var allSettings settings.Settings
	err = allSettings.Read(c.reader, c.logger)
	if err != nil {
		return err
	}
	allSettings.SetDefaults()

	ipv6SupportLevel, err := c.ipv6Checker.FindIPv6SupportLevel(ctx,
		allSettings.IPv6.CheckAddresses, &noopFirewall{})
	if err != nil {
		return fmt.Errorf("checking for IPv6 support: %w", err)
	}

	if err = allSettings.Validate(storage, ipv6SupportLevel.IsSupported(), c.logger); err != nil {
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
		unzipper, parallelResolver, ipFetcher, openvpnFileExtractor, allSettings.Updater)
	providerConf := providers.Get(allSettings.VPN.Provider.Name)
	connection, err := providerConf.GetConnection(
		allSettings.VPN.Provider.ServerSelection, ipv6SupportLevel == netlink.IPv6Internet)
	if err != nil {
		return err
	}

	lines := providerConf.OpenVPNConfig(connection,
		allSettings.VPN.OpenVPN, ipv6SupportLevel.IsSupported())

	fmt.Println(strings.Join(lines, "\n"))
	return nil
}
