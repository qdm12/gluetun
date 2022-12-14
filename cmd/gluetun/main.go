package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	_ "time/tzdata"

	_ "github.com/breml/rootcerts"
	"github.com/qdm12/dns/pkg/unbound"
	"github.com/qdm12/gluetun/internal/alpine"
	"github.com/qdm12/gluetun/internal/cli"
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/configuration/sources/env"
	"github.com/qdm12/gluetun/internal/configuration/sources/files"
	mux "github.com/qdm12/gluetun/internal/configuration/sources/merge"
	"github.com/qdm12/gluetun/internal/configuration/sources/secrets"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/dns"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/healthcheck"
	"github.com/qdm12/gluetun/internal/httpproxy"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/qdm12/gluetun/internal/openvpn"
	"github.com/qdm12/gluetun/internal/openvpn/extract"
	"github.com/qdm12/gluetun/internal/portforward"
	"github.com/qdm12/gluetun/internal/pprof"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/gluetun/internal/publicip"
	"github.com/qdm12/gluetun/internal/publicip/ipinfo"
	"github.com/qdm12/gluetun/internal/routing"
	"github.com/qdm12/gluetun/internal/server"
	"github.com/qdm12/gluetun/internal/shadowsocks"
	"github.com/qdm12/gluetun/internal/storage"
	"github.com/qdm12/gluetun/internal/tun"
	updater "github.com/qdm12/gluetun/internal/updater/loop"
	"github.com/qdm12/gluetun/internal/updater/resolver"
	"github.com/qdm12/gluetun/internal/updater/unzip"
	"github.com/qdm12/gluetun/internal/vpn"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/goshutdown"
	"github.com/qdm12/goshutdown/goroutine"
	"github.com/qdm12/goshutdown/group"
	"github.com/qdm12/goshutdown/order"
	"github.com/qdm12/gosplash"
	"github.com/qdm12/log"
	"github.com/qdm12/updated/pkg/dnscrypto"
)

//nolint:gochecknoglobals
var (
	version = "unknown"
	commit  = "unknown"
	created = "an unknown date"
)

func main() {
	buildInfo := models.BuildInformation{
		Version: version,
		Commit:  commit,
		Created: created,
	}

	background := context.Background()
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(background)

	logger := log.New(log.SetLevel(log.LevelInfo))

	args := os.Args
	tun := tun.New()
	netLinkDebugLogger := logger.New(log.SetComponent("netlink"))
	netLinker := netlink.New(netLinkDebugLogger)
	cli := cli.New()
	cmder := command.NewCmder()

	envReader := env.New(logger)
	filesReader := files.New()
	secretsReader := secrets.New()
	muxReader := mux.New(envReader, filesReader, secretsReader)

	errorCh := make(chan error)
	go func() {
		errorCh <- _main(ctx, buildInfo, args, logger, muxReader, tun, netLinker, cmder, cli)
	}()

	var err error
	select {
	case signal := <-signalCh:
		fmt.Println("")
		logger.Warn("Caught OS signal " + signal.String() + ", shutting down")
		cancel()
	case err = <-errorCh:
		close(errorCh)
		if err == nil { // expected exit such as healthcheck
			os.Exit(0)
		}
		logger.Error(err.Error())
		cancel()
	}

	const shutdownGracePeriod = 5 * time.Second
	timer := time.NewTimer(shutdownGracePeriod)
	select {
	case shutdownErr := <-errorCh:
		if !timer.Stop() {
			<-timer.C
		}
		if shutdownErr != nil {
			logger.Warnf("Shutdown not completed gracefully: %s", shutdownErr)
			os.Exit(1)
		}

		logger.Info("Shutdown successful")
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	case <-timer.C:
		logger.Warn("Shutdown timed out")
		os.Exit(1)
	case signal := <-signalCh:
		logger.Warn("Caught OS signal " + signal.String() + ", forcing shut down")
		os.Exit(1)
	}
}

var (
	errCommandUnknown = errors.New("command is unknown")
)

//nolint:gocognit,gocyclo,maintidx
func _main(ctx context.Context, buildInfo models.BuildInformation,
	args []string, logger log.LoggerInterface, source Source,
	tun Tun, netLinker netLinker, cmder command.RunStarter,
	cli clier) error {
	if len(args) > 1 { // cli operation
		switch args[1] {
		case "healthcheck":
			return cli.HealthCheck(ctx, source, logger)
		case "clientkey":
			return cli.ClientKey(args[2:])
		case "openvpnconfig":
			return cli.OpenvpnConfig(logger, source, netLinker)
		case "update":
			return cli.Update(ctx, args[2:], logger)
		case "format-servers":
			return cli.FormatServers(args[2:])
		default:
			return fmt.Errorf("%w: %s", errCommandUnknown, args[1])
		}
	}

	announcementExp, err := time.Parse(time.RFC3339, "2021-02-15T00:00:00Z")
	if err != nil {
		return err
	}
	splashSettings := gosplash.Settings{
		User:         "qdm12",
		Repository:   "gluetun",
		Emails:       []string{"quentin.mcgaw@gmail.com"},
		Version:      buildInfo.Version,
		Commit:       buildInfo.Commit,
		BuildDate:    buildInfo.Created,
		Announcement: "Large settings parsing refactoring merged on 2022-01-06, please report any issue!",
		AnnounceExp:  announcementExp,
		// Sponsor information
		PaypalUser:    "qmcgaw",
		GithubSponsor: "qdm12",
	}
	for _, line := range gosplash.MakeLines(splashSettings) {
		fmt.Println(line)
	}

	allSettings, err := source.Read()
	if err != nil {
		return err
	}

	// Note: no need to validate minimal settings for the firewall:
	// - global log level is parsed from source
	// - firewall Debug and Enabled are booleans parsed from source

	logger.Patch(log.SetLevel(*allSettings.Log.Level))
	netLinker.PatchLoggerLevel(*allSettings.Log.Level)

	routingLogger := logger.New(log.SetComponent("routing"))
	if *allSettings.Firewall.Debug { // To remove in v4
		routingLogger.Patch(log.SetLevel(log.LevelDebug))
	}
	routingConf := routing.New(netLinker, routingLogger)

	defaultRoutes, err := routingConf.DefaultRoutes()
	if err != nil {
		return err
	}

	localNetworks, err := routingConf.LocalNetworks()
	if err != nil {
		return err
	}

	firewallLogger := logger.New(log.SetComponent("firewall"))
	if *allSettings.Firewall.Debug { // To remove in v4
		firewallLogger.Patch(log.SetLevel(log.LevelDebug))
	}
	firewallConf, err := firewall.NewConfig(ctx, firewallLogger, cmder,
		defaultRoutes, localNetworks)
	if err != nil {
		return err
	}

	if *allSettings.Firewall.Enabled {
		err = firewallConf.SetEnabled(ctx, true)
		if err != nil {
			return err
		}
	}

	// TODO run this in a loop or in openvpn to reload from file without restarting
	storageLogger := logger.New(log.SetComponent("storage"))
	storage, err := storage.New(storageLogger, constants.ServersData)
	if err != nil {
		return err
	}

	ipv6Supported, err := netLinker.IsIPv6Supported()
	if err != nil {
		return fmt.Errorf("checking for IPv6 support: %w", err)
	}

	err = allSettings.Validate(storage, ipv6Supported)
	if err != nil {
		return err
	}

	allSettings.Pprof.HTTPServer.Logger = logger.New(log.SetComponent("pprof"))
	pprofServer, err := pprof.New(allSettings.Pprof)
	if err != nil {
		return fmt.Errorf("cannot create Pprof server: %w", err)
	}

	puid, pgid := int(*allSettings.System.PUID), int(*allSettings.System.PGID)

	const clientTimeout = 15 * time.Second
	httpClient := &http.Client{Timeout: clientTimeout}
	// Create configurators
	alpineConf := alpine.New()
	ovpnConf := openvpn.New(
		logger.New(log.SetComponent("openvpn configurator")),
		cmder, puid, pgid)
	dnsCrypto := dnscrypto.New(httpClient, "", "")
	const cacertsPath = "/etc/ssl/certs/ca-certificates.crt"
	dnsConf := unbound.NewConfigurator(nil, cmder, dnsCrypto,
		"/etc/unbound", "/usr/sbin/unbound", cacertsPath)

	err = printVersions(ctx, logger, []printVersionElement{
		{name: "Alpine", getVersion: alpineConf.Version},
		{name: "OpenVPN 2.4", getVersion: ovpnConf.Version24},
		{name: "OpenVPN 2.5", getVersion: ovpnConf.Version25},
		{name: "Unbound", getVersion: dnsConf.Version},
		{name: "IPtables", getVersion: func(ctx context.Context) (version string, err error) {
			return firewall.Version(ctx, cmder)
		}},
	})
	if err != nil {
		return err
	}

	logger.Info(allSettings.String())

	if err := os.MkdirAll("/tmp/gluetun", 0644); err != nil {
		return err
	}
	if err := os.MkdirAll("/gluetun", 0644); err != nil {
		return err
	}

	const defaultUsername = "nonrootuser"
	nonRootUsername, err := alpineConf.CreateUser(defaultUsername, puid)
	if err != nil {
		return fmt.Errorf("cannot create user: %w", err)
	}
	if nonRootUsername != defaultUsername {
		logger.Info("using existing username " + nonRootUsername + " corresponding to user id " + fmt.Sprint(puid))
	}
	// set it for Unbound
	// TODO remove this when migrating to qdm12/dns v2
	allSettings.DNS.DoT.Unbound.Username = nonRootUsername
	allSettings.VPN.OpenVPN.ProcessUser = nonRootUsername

	if err := os.Chown("/etc/unbound", puid, pgid); err != nil {
		return err
	}

	if err := routingConf.Setup(); err != nil {
		if strings.Contains(err.Error(), "operation not permitted") {
			logger.Warn("💡 Tip: Are you passing NET_ADMIN capability to gluetun?")
		}
		return fmt.Errorf("cannot setup routing: %w", err)
	}
	defer func() {
		routingLogger.Info("routing cleanup...")
		if err := routingConf.TearDown(); err != nil {
			routingLogger.Error("cannot teardown routing: " + err.Error())
		}
	}()

	if err := firewallConf.SetOutboundSubnets(ctx, allSettings.Firewall.OutboundSubnets); err != nil {
		return err
	}
	if err := routingConf.SetOutboundRoutes(allSettings.Firewall.OutboundSubnets); err != nil {
		return err
	}

	const tunDevice = "/dev/net/tun"
	if err := tun.Check(tunDevice); err != nil {
		logger.Info(err.Error() + "; creating it...")
		err = tun.Create(tunDevice)
		if err != nil {
			return err
		}
	}

	for _, port := range allSettings.Firewall.InputPorts {
		for _, defaultRoute := range defaultRoutes {
			err = firewallConf.SetAllowedPort(ctx, port, defaultRoute.NetInterface)
			if err != nil {
				return err
			}
		}
	} // TODO move inside firewall?

	// Shutdown settings
	const totalShutdownTimeout = 3 * time.Second
	const defaultShutdownTimeout = 400 * time.Millisecond
	defaultShutdownOnSuccess := func(goRoutineName string) {
		logger.Info(goRoutineName + ": terminated ✔️")
	}
	defaultShutdownOnFailure := func(goRoutineName string, err error) {
		logger.Warn(goRoutineName + ": " + err.Error() + " ⚠️")
	}
	defaultGroupOptions := []group.Option{
		group.OptionTimeout(defaultShutdownTimeout),
		group.OptionOnSuccess(defaultShutdownOnSuccess)}

	controlGroupHandler := goshutdown.NewGroupHandler("control", defaultGroupOptions...)
	tickersGroupHandler := goshutdown.NewGroupHandler("tickers", defaultGroupOptions...)
	otherGroupHandler := goshutdown.NewGroupHandler("other", defaultGroupOptions...)

	if *allSettings.Pprof.Enabled {
		// TODO run in run loop so this can be patched at runtime
		pprofReady := make(chan struct{})
		pprofHandler, pprofCtx, pprofDone := goshutdown.NewGoRoutineHandler("pprof server")
		go pprofServer.Run(pprofCtx, pprofReady, pprofDone)
		otherGroupHandler.Add(pprofHandler)
		<-pprofReady
	}

	portForwardLogger := logger.New(log.SetComponent("port forwarding"))
	portForwardLooper := portforward.NewLoop(allSettings.VPN.Provider.PortForwarding,
		httpClient, firewallConf, portForwardLogger, puid, pgid)
	portForwardHandler, portForwardCtx, portForwardDone := goshutdown.NewGoRoutineHandler(
		"port forwarding", goroutine.OptionTimeout(time.Second))
	go portForwardLooper.Run(portForwardCtx, portForwardDone)

	unboundLogger := logger.New(log.SetComponent("dns over tls"))
	unboundLooper := dns.NewLoop(dnsConf, allSettings.DNS, httpClient,
		unboundLogger)
	dnsHandler, dnsCtx, dnsDone := goshutdown.NewGoRoutineHandler(
		"unbound", goroutine.OptionTimeout(defaultShutdownTimeout))
	// wait for unboundLooper.Restart or its ticker launched with RunRestartTicker
	go unboundLooper.Run(dnsCtx, dnsDone)
	otherGroupHandler.Add(dnsHandler)

	dnsTickerHandler, dnsTickerCtx, dnsTickerDone := goshutdown.NewGoRoutineHandler(
		"dns ticker", goroutine.OptionTimeout(defaultShutdownTimeout))
	go unboundLooper.RunRestartTicker(dnsTickerCtx, dnsTickerDone)
	controlGroupHandler.Add(dnsTickerHandler)

	ipFetcher := ipinfo.New(httpClient)
	publicIPLooper := publicip.NewLoop(ipFetcher,
		logger.New(log.SetComponent("ip getter")),
		allSettings.PublicIP, puid, pgid)
	pubIPHandler, pubIPCtx, pubIPDone := goshutdown.NewGoRoutineHandler(
		"public IP", goroutine.OptionTimeout(defaultShutdownTimeout))
	go publicIPLooper.Run(pubIPCtx, pubIPDone)
	otherGroupHandler.Add(pubIPHandler)

	pubIPTickerHandler, pubIPTickerCtx, pubIPTickerDone := goshutdown.NewGoRoutineHandler(
		"public IP", goroutine.OptionTimeout(defaultShutdownTimeout))
	go publicIPLooper.RunRestartTicker(pubIPTickerCtx, pubIPTickerDone)
	tickersGroupHandler.Add(pubIPTickerHandler)

	updaterLogger := logger.New(log.SetComponent("updater"))

	unzipper := unzip.New(httpClient)
	parallelResolver := resolver.NewParallelResolver(allSettings.Updater.DNSAddress)
	openvpnFileExtractor := extract.New()
	providers := provider.NewProviders(storage, time.Now, updaterLogger,
		httpClient, unzipper, parallelResolver, ipFetcher, openvpnFileExtractor)

	vpnLogger := logger.New(log.SetComponent("vpn"))
	vpnLooper := vpn.NewLoop(allSettings.VPN, ipv6Supported, allSettings.Firewall.VPNInputPorts,
		providers, storage, ovpnConf, netLinker, firewallConf, routingConf, portForwardLooper,
		cmder, publicIPLooper, unboundLooper, vpnLogger, httpClient,
		buildInfo, *allSettings.Version.Enabled)
	vpnHandler, vpnCtx, vpnDone := goshutdown.NewGoRoutineHandler(
		"vpn", goroutine.OptionTimeout(time.Second))
	go vpnLooper.Run(vpnCtx, vpnDone)

	updaterLooper := updater.NewLoop(allSettings.Updater,
		providers, storage, httpClient, updaterLogger)
	updaterHandler, updaterCtx, updaterDone := goshutdown.NewGoRoutineHandler(
		"updater", goroutine.OptionTimeout(defaultShutdownTimeout))
	// wait for updaterLooper.Restart() or its ticket launched with RunRestartTicker
	go updaterLooper.Run(updaterCtx, updaterDone)
	tickersGroupHandler.Add(updaterHandler)

	updaterTickerHandler, updaterTickerCtx, updaterTickerDone := goshutdown.NewGoRoutineHandler(
		"updater ticker", goroutine.OptionTimeout(defaultShutdownTimeout))
	go updaterLooper.RunRestartTicker(updaterTickerCtx, updaterTickerDone)
	controlGroupHandler.Add(updaterTickerHandler)

	httpProxyLooper := httpproxy.NewLoop(
		logger.New(log.SetComponent("http proxy")),
		allSettings.HTTPProxy)
	httpProxyHandler, httpProxyCtx, httpProxyDone := goshutdown.NewGoRoutineHandler(
		"http proxy", goroutine.OptionTimeout(defaultShutdownTimeout))
	go httpProxyLooper.Run(httpProxyCtx, httpProxyDone)
	otherGroupHandler.Add(httpProxyHandler)

	shadowsocksLooper := shadowsocks.NewLoop(allSettings.Shadowsocks,
		logger.New(log.SetComponent("shadowsocks")))
	shadowsocksHandler, shadowsocksCtx, shadowsocksDone := goshutdown.NewGoRoutineHandler(
		"shadowsocks proxy", goroutine.OptionTimeout(defaultShutdownTimeout))
	go shadowsocksLooper.Run(shadowsocksCtx, shadowsocksDone)
	otherGroupHandler.Add(shadowsocksHandler)

	controlServerAddress := *allSettings.ControlServer.Address
	controlServerLogging := *allSettings.ControlServer.Log
	httpServerHandler, httpServerCtx, httpServerDone := goshutdown.NewGoRoutineHandler(
		"http server", goroutine.OptionTimeout(defaultShutdownTimeout))
	httpServer, err := server.New(httpServerCtx, controlServerAddress, controlServerLogging,
		logger.New(log.SetComponent("http server")),
		buildInfo, vpnLooper, portForwardLooper, unboundLooper, updaterLooper, publicIPLooper,
		storage, ipv6Supported)
	if err != nil {
		return fmt.Errorf("cannot setup control server: %w", err)
	}
	httpServerReady := make(chan struct{})
	go httpServer.Run(httpServerCtx, httpServerReady, httpServerDone)
	<-httpServerReady
	controlGroupHandler.Add(httpServerHandler)

	healthLogger := logger.New(log.SetComponent("healthcheck"))
	healthcheckServer := healthcheck.NewServer(allSettings.Health, healthLogger, vpnLooper)
	healthServerHandler, healthServerCtx, healthServerDone := goshutdown.NewGoRoutineHandler(
		"HTTP health server", goroutine.OptionTimeout(defaultShutdownTimeout))
	go healthcheckServer.Run(healthServerCtx, healthServerDone)

	orderHandler := goshutdown.NewOrderHandler("gluetun",
		order.OptionTimeout(totalShutdownTimeout),
		order.OptionOnSuccess(defaultShutdownOnSuccess),
		order.OptionOnFailure(defaultShutdownOnFailure))
	orderHandler.Append(controlGroupHandler, tickersGroupHandler, healthServerHandler,
		vpnHandler, portForwardHandler, otherGroupHandler)

	// Start VPN for the first time in a blocking call
	// until the VPN is launched
	_, _ = vpnLooper.ApplyStatus(ctx, constants.Running) // TODO option to disable with variable

	<-ctx.Done()

	return orderHandler.Shutdown(context.Background())
}

type printVersionElement struct {
	name       string
	getVersion func(ctx context.Context) (version string, err error)
}

type infoer interface {
	Info(s string)
}

func printVersions(ctx context.Context, logger infoer,
	elements []printVersionElement) (err error) {
	const timeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for _, element := range elements {
		version, err := element.getVersion(ctx)
		if err != nil {
			return err
		}
		logger.Info(element.name + " version: " + version)
	}

	return nil
}

type netLinker interface {
	Addresser
	Router
	Ruler
	Linker
	IsWireguardSupported() (ok bool, err error)
	IsIPv6Supported() (ok bool, err error)
	PatchLoggerLevel(level log.Level)
}

type Addresser interface {
	AddrList(link netlink.Link, family int) (
		addresses []netlink.Addr, err error)
	AddrAdd(link netlink.Link, addr *netlink.Addr) error
}

type Router interface {
	RouteList(link netlink.Link, family int) (
		routes []netlink.Route, err error)
	RouteAdd(route *netlink.Route) error
	RouteDel(route *netlink.Route) error
	RouteReplace(route *netlink.Route) error
}

type Ruler interface {
	RuleList(family int) (rules []netlink.Rule, err error)
	RuleAdd(rule *netlink.Rule) error
	RuleDel(rule *netlink.Rule) error
}

type Linker interface {
	LinkList() (links []netlink.Link, err error)
	LinkByName(name string) (link netlink.Link, err error)
	LinkByIndex(index int) (link netlink.Link, err error)
	LinkAdd(link netlink.Link) (err error)
	LinkDel(link netlink.Link) (err error)
	LinkSetUp(link netlink.Link) (err error)
	LinkSetDown(link netlink.Link) (err error)
}

type clier interface {
	ClientKey(args []string) error
	FormatServers(args []string) error
	OpenvpnConfig(logger cli.OpenvpnConfigLogger, source cli.Source, ipv6Checker cli.IPv6Checker) error
	HealthCheck(ctx context.Context, source cli.Source, warner cli.Warner) error
	Update(ctx context.Context, args []string, logger cli.UpdaterLogger) error
}

type Tun interface {
	Check(tunDevice string) error
	Create(tunDevice string) error
}

type Source interface {
	Read() (settings settings.Settings, err error)
	ReadHealth() (health settings.Health, err error)
	String() string
}
