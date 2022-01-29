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
	"github.com/qdm12/gluetun/internal/configuration/sources"
	"github.com/qdm12/gluetun/internal/configuration/sources/env"
	"github.com/qdm12/gluetun/internal/configuration/sources/files"
	"github.com/qdm12/gluetun/internal/configuration/sources/mux"
	"github.com/qdm12/gluetun/internal/configuration/sources/secrets"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/dns"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/healthcheck"
	"github.com/qdm12/gluetun/internal/httpproxy"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/qdm12/gluetun/internal/openvpn"
	"github.com/qdm12/gluetun/internal/portforward"
	"github.com/qdm12/gluetun/internal/pprof"
	"github.com/qdm12/gluetun/internal/publicip"
	"github.com/qdm12/gluetun/internal/routing"
	"github.com/qdm12/gluetun/internal/server"
	"github.com/qdm12/gluetun/internal/shadowsocks"
	"github.com/qdm12/gluetun/internal/storage"
	"github.com/qdm12/gluetun/internal/tun"
	"github.com/qdm12/gluetun/internal/updater"
	"github.com/qdm12/gluetun/internal/vpn"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/goshutdown"
	"github.com/qdm12/goshutdown/goroutine"
	"github.com/qdm12/goshutdown/group"
	"github.com/qdm12/goshutdown/order"
	"github.com/qdm12/gosplash"
	"github.com/qdm12/updated/pkg/dnscrypto"
)

//nolint:gochecknoglobals
var (
	version = "unknown"
	commit  = "unknown"
	created = "an unknown date"
)

var (
	errSetupRouting = errors.New("cannot setup routing")
	errCreateUser   = errors.New("cannot create user")
)

func main() {
	buildInfo := models.BuildInformation{
		Version: version,
		Commit:  commit,
		Created: created,
	}

	background := context.Background()
	signalCtx, stop := signal.NotifyContext(background, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	ctx, cancel := context.WithCancel(background)

	logger := logging.New(logging.Settings{
		Level: logging.LevelInfo,
	})

	args := os.Args
	tun := tun.New()
	netLinker := netlink.New()
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

	select {
	case <-signalCtx.Done():
		stop()
		fmt.Println("")
		logger.Warn("Caught OS signal, shutting down")
		cancel()
	case err := <-errorCh:
		stop()
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
	case <-errorCh:
		if !timer.Stop() {
			<-timer.C
		}
		logger.Info("Shutdown successful")
	case <-timer.C:
		logger.Warn("Shutdown timed out")
	}

	os.Exit(1)
}

var (
	errCommandUnknown = errors.New("command is unknown")
)

//nolint:gocognit,gocyclo
func _main(ctx context.Context, buildInfo models.BuildInformation,
	args []string, logger logging.ParentLogger, source sources.Source,
	tun tun.Interface, netLinker netlink.NetLinker, cmder command.RunStarter,
	cli cli.CLIer) error {
	if len(args) > 1 { // cli operation
		switch args[1] {
		case "healthcheck":
			return cli.HealthCheck(ctx, source, logger)
		case "clientkey":
			return cli.ClientKey(args[2:])
		case "openvpnconfig":
			return cli.OpenvpnConfig(logger, source)
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

	// TODO run this in a loop or in openvpn to reload from file without restarting
	storageLogger := logger.NewChild(logging.Settings{Prefix: "storage: "})
	storage, err := storage.New(storageLogger, constants.ServersData)
	if err != nil {
		return err
	}

	allServers := storage.GetServers()

	err = allSettings.Validate(allServers)
	if err != nil {
		return err
	}

	logger.PatchLevel(*allSettings.Log.Level)

	allSettings.Pprof.HTTPServer.Logger = logger
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
		logger.NewChild(logging.Settings{Prefix: "openvpn configurator: "}),
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
		return fmt.Errorf("%w: %s", errCreateUser, err)
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

	firewallLogLevel := *allSettings.Log.Level
	if *allSettings.Firewall.Debug {
		firewallLogLevel = logging.LevelDebug
	}
	routingLogger := logger.NewChild(logging.Settings{
		Prefix: "routing: ",
		Level:  firewallLogLevel,
	})
	routingConf := routing.New(netLinker, routingLogger)

	defaultInterface, defaultGateway, err := routingConf.DefaultRoute()
	if err != nil {
		return err
	}

	localNetworks, err := routingConf.LocalNetworks()
	if err != nil {
		return err
	}

	defaultIP, err := routingConf.DefaultIP()
	if err != nil {
		return err
	}

	firewallLogger := logger.NewChild(logging.Settings{
		Prefix: "firewall: ",
		Level:  firewallLogLevel,
	})
	firewallConf := firewall.NewConfig(firewallLogger, cmder,
		defaultInterface, defaultGateway, localNetworks, defaultIP)

	if err := routingConf.Setup(); err != nil {
		if strings.Contains(err.Error(), "operation not permitted") {
			logger.Warn("💡 Tip: Are you passing NET_ADMIN capability to gluetun?")
		}
		return fmt.Errorf("%w: %s", errSetupRouting, err)
	}
	defer func() {
		logger.Info("routing cleanup...")
		if err := routingConf.TearDown(); err != nil {
			logger.Error("cannot teardown routing: " + err.Error())
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

	if *allSettings.Firewall.Enabled {
		err := firewallConf.SetEnabled(ctx, true) // disabled by default
		if err != nil {
			return err
		}
	}

	for _, port := range allSettings.Firewall.InputPorts {
		err = firewallConf.SetAllowedPort(ctx, port, defaultInterface)
		if err != nil {
			return err
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

	pprofReady := make(chan struct{})
	pprofHandler, pprofCtx, pprofDone := goshutdown.NewGoRoutineHandler("pprof server")
	go pprofServer.Run(pprofCtx, pprofReady, pprofDone)
	otherGroupHandler.Add(pprofHandler)
	<-pprofReady

	portForwardLogger := logger.NewChild(logging.Settings{Prefix: "port forwarding: "})
	portForwardLooper := portforward.NewLoop(allSettings.VPN.Provider.PortForwarding,
		httpClient, firewallConf, portForwardLogger)
	portForwardHandler, portForwardCtx, portForwardDone := goshutdown.NewGoRoutineHandler(
		"port forwarding", goroutine.OptionTimeout(time.Second))
	go portForwardLooper.Run(portForwardCtx, portForwardDone)

	unboundLogger := logger.NewChild(logging.Settings{Prefix: "dns over tls: "})
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

	publicIPLooper := publicip.NewLoop(httpClient,
		logger.NewChild(logging.Settings{Prefix: "ip getter: "}),
		allSettings.PublicIP, puid, pgid)
	pubIPHandler, pubIPCtx, pubIPDone := goshutdown.NewGoRoutineHandler(
		"public IP", goroutine.OptionTimeout(defaultShutdownTimeout))
	go publicIPLooper.Run(pubIPCtx, pubIPDone)
	otherGroupHandler.Add(pubIPHandler)

	pubIPTickerHandler, pubIPTickerCtx, pubIPTickerDone := goshutdown.NewGoRoutineHandler(
		"public IP", goroutine.OptionTimeout(defaultShutdownTimeout))
	go publicIPLooper.RunRestartTicker(pubIPTickerCtx, pubIPTickerDone)
	tickersGroupHandler.Add(pubIPTickerHandler)

	vpnLogger := logger.NewChild(logging.Settings{Prefix: "vpn: "})
	vpnLooper := vpn.NewLoop(allSettings.VPN, allSettings.Firewall.VPNInputPorts,
		allServers, ovpnConf, netLinker, firewallConf, routingConf, portForwardLooper,
		cmder, publicIPLooper, unboundLooper, vpnLogger, httpClient,
		buildInfo, *allSettings.Version.Enabled)
	vpnHandler, vpnCtx, vpnDone := goshutdown.NewGoRoutineHandler(
		"vpn", goroutine.OptionTimeout(time.Second))
	go vpnLooper.Run(vpnCtx, vpnDone)

	updaterLooper := updater.NewLooper(allSettings.Updater,
		allServers, storage, vpnLooper.SetServers, httpClient,
		logger.NewChild(logging.Settings{Prefix: "updater: "}))
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
		logger.NewChild(logging.Settings{Prefix: "http proxy: "}),
		allSettings.HTTPProxy)
	httpProxyHandler, httpProxyCtx, httpProxyDone := goshutdown.NewGoRoutineHandler(
		"http proxy", goroutine.OptionTimeout(defaultShutdownTimeout))
	go httpProxyLooper.Run(httpProxyCtx, httpProxyDone)
	otherGroupHandler.Add(httpProxyHandler)

	shadowsocksLooper := shadowsocks.NewLooper(allSettings.Shadowsocks,
		logger.NewChild(logging.Settings{Prefix: "shadowsocks: "}))
	shadowsocksHandler, shadowsocksCtx, shadowsocksDone := goshutdown.NewGoRoutineHandler(
		"shadowsocks proxy", goroutine.OptionTimeout(defaultShutdownTimeout))
	go shadowsocksLooper.Run(shadowsocksCtx, shadowsocksDone)
	otherGroupHandler.Add(shadowsocksHandler)

	controlServerAddress := *allSettings.ControlServer.Address
	controlServerLogging := *allSettings.ControlServer.Log
	httpServerHandler, httpServerCtx, httpServerDone := goshutdown.NewGoRoutineHandler(
		"http server", goroutine.OptionTimeout(defaultShutdownTimeout))
	httpServer := server.New(httpServerCtx, controlServerAddress, controlServerLogging,
		logger.NewChild(logging.Settings{Prefix: "http server: "}),
		buildInfo, vpnLooper, portForwardLooper, unboundLooper, updaterLooper, publicIPLooper)
	go httpServer.Run(httpServerCtx, httpServerDone)
	controlGroupHandler.Add(httpServerHandler)

	healthLogger := logger.NewChild(logging.Settings{Prefix: "healthcheck: "})
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
