package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	nativeos "os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/qdm12/dns/pkg/unbound"
	"github.com/qdm12/gluetun/internal/alpine"
	"github.com/qdm12/gluetun/internal/cli"
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/dns"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/healthcheck"
	"github.com/qdm12/gluetun/internal/httpproxy"
	gluetunLogging "github.com/qdm12/gluetun/internal/logging"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/openvpn"
	"github.com/qdm12/gluetun/internal/publicip"
	"github.com/qdm12/gluetun/internal/routing"
	"github.com/qdm12/gluetun/internal/server"
	"github.com/qdm12/gluetun/internal/shadowsocks"
	"github.com/qdm12/gluetun/internal/shutdown"
	"github.com/qdm12/gluetun/internal/storage"
	"github.com/qdm12/gluetun/internal/unix"
	"github.com/qdm12/gluetun/internal/updater"
	versionpkg "github.com/qdm12/gluetun/internal/version"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
	"github.com/qdm12/golibs/os/user"
	"github.com/qdm12/golibs/params"
	"github.com/qdm12/updated/pkg/dnscrypto"
)

//nolint:gochecknoglobals
var (
	version   = "unknown"
	commit    = "unknown"
	buildDate = "an unknown date"
)

var (
	errSetupRouting = errors.New("cannot setup routing")
	errCreateUser   = errors.New("cannot create user")
)

func main() {
	buildInfo := models.BuildInformation{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
	}

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, nativeos.Interrupt)
	ctx, cancel := context.WithCancel(ctx)

	logger := logging.NewParent(logging.Settings{})

	args := nativeos.Args
	os := os.New()
	osUser := user.New()
	unix := unix.New()
	cli := cli.New()

	errorCh := make(chan error)
	go func() {
		errorCh <- _main(ctx, buildInfo, args, logger, os, osUser, unix, cli)
	}()

	select {
	case <-ctx.Done():
		stop()
		logger.Warn("Caught OS signal, shutting down")
	case err := <-errorCh:
		stop()
		close(errorCh)
		if err == nil { // expected exit such as healthcheck
			nativeos.Exit(0)
		}
		logger.Error(err)
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

	nativeos.Exit(1)
}

var (
	errCommandUnknown = errors.New("command is unknown")
)

//nolint:gocognit,gocyclo
func _main(ctx context.Context, buildInfo models.BuildInformation,
	args []string, logger logging.ParentLogger, os os.OS,
	osUser user.OSUser, unix unix.Unix, cli cli.CLI) error {
	if len(args) > 1 { // cli operation
		switch args[1] {
		case "healthcheck":
			return cli.HealthCheck(ctx)
		case "clientkey":
			return cli.ClientKey(args[2:], os.OpenFile)
		case "openvpnconfig":
			return cli.OpenvpnConfig(os, logger)
		case "update":
			return cli.Update(ctx, args[2:], os, logger)
		default:
			return fmt.Errorf("%w: %s", errCommandUnknown, args[1])
		}
	}

	const clientTimeout = 15 * time.Second
	httpClient := &http.Client{Timeout: clientTimeout}
	// Create configurators
	alpineConf := alpine.NewConfigurator(os.OpenFile, osUser)
	ovpnConf := openvpn.NewConfigurator(
		logger.NewChild(logging.Settings{Prefix: "openvpn configurator: "}),
		os, unix)
	dnsCrypto := dnscrypto.New(httpClient, "", "")
	const cacertsPath = "/etc/ssl/certs/ca-certificates.crt"
	dnsConf := unbound.NewConfigurator(nil, os.OpenFile, dnsCrypto,
		"/etc/unbound", "/usr/sbin/unbound", cacertsPath)
	routingConf := routing.NewRouting(
		logger.NewChild(logging.Settings{Prefix: "routing: "}))
	firewallConf := firewall.NewConfigurator(
		logger.NewChild(logging.Settings{Prefix: "firewall: "}),
		routingConf, os.OpenFile)

	fmt.Println(gluetunLogging.Splash(buildInfo))

	printVersions(ctx, logger, map[string]func(ctx context.Context) (string, error){
		"Alpine":   alpineConf.Version,
		"OpenVPN":  ovpnConf.Version,
		"Unbound":  dnsConf.Version,
		"IPtables": firewallConf.Version,
	})

	var allSettings configuration.Settings
	err := allSettings.Read(params.NewEnv(), os,
		logger.NewChild(logging.Settings{Prefix: "configuration: "}))
	if err != nil {
		return err
	}
	logger.Info(allSettings.String())

	allSettings.OpenVPN.Version, _ = ovpnConf.Version(ctx)

	if err := os.MkdirAll("/tmp/gluetun", 0644); err != nil {
		return err
	}
	if err := os.MkdirAll("/gluetun", 0644); err != nil {
		return err
	}

	// TODO run this in a loop or in openvpn to reload from file without restarting
	storage := storage.New(
		logger.NewChild(logging.Settings{Prefix: "storage: "}),
		os, constants.ServersData)
	allServers, err := storage.SyncServers(constants.GetAllServers())
	if err != nil {
		return err
	}

	// Should never change
	puid, pgid := allSettings.System.PUID, allSettings.System.PGID

	const defaultUsername = "nonrootuser"
	nonRootUsername, err := alpineConf.CreateUser(defaultUsername, puid)
	if err != nil {
		return fmt.Errorf("%w: %s", errCreateUser, err)
	}
	if nonRootUsername != defaultUsername {
		logger.Info("using existing username %s corresponding to user id %d", nonRootUsername, puid)
	}
	// set it for Unbound
	// TODO remove this when migrating to qdm12/dns v2
	allSettings.DNS.Unbound.Username = nonRootUsername

	if err := os.Chown("/etc/unbound", puid, pgid); err != nil {
		return err
	}

	if allSettings.Firewall.Debug {
		firewallConf.SetDebug()
		routingConf.SetDebug()
	}

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

	firewallConf.SetNetworkInformation(defaultInterface, defaultGateway, localNetworks, defaultIP)

	if err := routingConf.Setup(); err != nil {
		if strings.Contains(err.Error(), "operation not permitted") {
			logger.Warn("💡 Tip: Are you passing NET_ADMIN capability to gluetun?")
		}
		return fmt.Errorf("%w: %s", errSetupRouting, err)
	}
	defer func() {
		routingConf.SetVerbose(false)
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

	if err := ovpnConf.CheckTUN(); err != nil {
		logger.Warn(err)
		err = ovpnConf.CreateTUN()
		if err != nil {
			return err
		}
	}

	tunnelReadyCh := make(chan struct{})
	defer close(tunnelReadyCh)

	if allSettings.Firewall.Enabled {
		err := firewallConf.SetEnabled(ctx, true) // disabled by default
		if err != nil {
			return err
		}
	}

	for _, vpnPort := range allSettings.Firewall.VPNInputPorts {
		err = firewallConf.SetAllowedPort(ctx, vpnPort, string(constants.TUN))
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

	const (
		shutdownMaxTimeout     = 3 * time.Second
		shutdownRoutineTimeout = 400 * time.Millisecond
		shutdownOpenvpnTimeout = time.Second
	)

	healthy := make(chan bool)
	controlWave := shutdown.NewWave("control")
	tickerWave := shutdown.NewWave("tickers")
	healthWave := shutdown.NewWave("health")
	dnsWave := shutdown.NewWave("DNS")
	vpnWave := shutdown.NewWave("VPN")
	serverWave := shutdown.NewWave("servers")

	openvpnLooper := openvpn.NewLooper(allSettings.OpenVPN, nonRootUsername, puid, pgid, allServers,
		ovpnConf, firewallConf, routingConf, logger, httpClient, os.OpenFile, tunnelReadyCh, healthy)
	openvpnCtx, openvpnDone := vpnWave.Add("openvpn", shutdownOpenvpnTimeout)
	// wait for restartOpenvpn
	go openvpnLooper.Run(openvpnCtx, openvpnDone)

	updaterLooper := updater.NewLooper(allSettings.Updater,
		allServers, storage, openvpnLooper.SetServers, httpClient,
		logger.NewChild(logging.Settings{Prefix: "updater: "}))
	updaterCtx, updaterDone := tickerWave.Add("updater", shutdownRoutineTimeout)
	// wait for updaterLooper.Restart() or its ticket launched with RunRestartTicker
	go updaterLooper.Run(updaterCtx, updaterDone)

	unboundLogger := logger.NewChild(logging.Settings{Prefix: "dns over tls: "})
	unboundLooper := dns.NewLooper(dnsConf, allSettings.DNS, httpClient,
		unboundLogger, os.OpenFile)
	dnsCtx, dnsDone := dnsWave.Add("unbound", shutdownRoutineTimeout)
	// wait for unboundLooper.Restart or its ticker launched with RunRestartTicker
	go unboundLooper.Run(dnsCtx, dnsDone)

	publicIPLooper := publicip.NewLooper(httpClient,
		logger.NewChild(logging.Settings{Prefix: "ip getter: "}),
		allSettings.PublicIP, puid, pgid, os)
	pubIPCtx, pubIPDone := serverWave.Add("public IP", shutdownRoutineTimeout)
	go publicIPLooper.Run(pubIPCtx, pubIPDone)

	pubIPTickerCtx, pubIPTickerDone := tickerWave.Add("public IP", shutdownRoutineTimeout)
	go publicIPLooper.RunRestartTicker(pubIPTickerCtx, pubIPTickerDone)

	httpProxyLooper := httpproxy.NewLooper(
		logger.NewChild(logging.Settings{Prefix: "http proxy: "}),
		allSettings.HTTPProxy)
	httpProxyCtx, httpProxyDone := serverWave.Add("http proxy", shutdownRoutineTimeout)
	go httpProxyLooper.Run(httpProxyCtx, httpProxyDone)

	shadowsocksLooper := shadowsocks.NewLooper(allSettings.ShadowSocks,
		logger.NewChild(logging.Settings{Prefix: "shadowsocks: "}))
	shadowsocksCtx, shadowsocksDone := serverWave.Add("shadowsocks proxy", shutdownRoutineTimeout)
	go shadowsocksLooper.Run(shadowsocksCtx, shadowsocksDone)

	eventsRoutingCtx, eventsRoutingDone := controlWave.Add("events routing", shutdownRoutineTimeout)
	go routeReadyEvents(eventsRoutingCtx, eventsRoutingDone, buildInfo, tunnelReadyCh,
		unboundLooper, updaterLooper, publicIPLooper, routingConf, logger, httpClient,
		allSettings.VersionInformation, allSettings.OpenVPN.Provider.PortForwarding.Enabled, openvpnLooper.PortForward,
	)
	controlServerAddress := ":" + strconv.Itoa(int(allSettings.ControlServer.Port))
	controlServerLogging := allSettings.ControlServer.Log
	httpServer := server.New(controlServerAddress, controlServerLogging,
		logger.NewChild(logging.Settings{Prefix: "http server: "}),
		buildInfo, openvpnLooper, unboundLooper, updaterLooper, publicIPLooper)
	httpServerCtx, httpServerDone := controlWave.Add("http server", shutdownRoutineTimeout)
	go httpServer.Run(httpServerCtx, httpServerDone)

	healthcheckServer := healthcheck.NewServer(constants.HealthcheckAddress,
		logger.NewChild(logging.Settings{Prefix: "healthcheck: "}))
	healthServerCtx, healthServerDone := healthWave.Add("HTTP health server", shutdownRoutineTimeout)
	go healthcheckServer.Run(healthServerCtx, healthy, healthServerDone)

	shutdownOrder := shutdown.NewOrder()
	shutdownOrder.Append(controlWave, tickerWave, healthWave,
		dnsWave, vpnWave, serverWave,
	)

	// Start openvpn for the first time in a blocking call
	// until openvpn is launched
	_, _ = openvpnLooper.SetStatus(constants.Running) // TODO option to disable with variable

	<-ctx.Done()

	if allSettings.OpenVPN.Provider.PortForwarding.Enabled {
		logger.Info("Clearing forwarded port status file %s", allSettings.OpenVPN.Provider.PortForwarding.Filepath)
		if err := os.Remove(allSettings.OpenVPN.Provider.PortForwarding.Filepath); err != nil {
			logger.Error(err)
		}
	}

	return shutdownOrder.Shutdown(shutdownMaxTimeout, logger)
}

func printVersions(ctx context.Context, logger logging.Logger,
	versionFunctions map[string]func(ctx context.Context) (string, error)) {
	const timeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for name, f := range versionFunctions {
		version, err := f(ctx)
		if err != nil {
			logger.Error(err)
		} else {
			logger.Info("%s version: %s", name, version)
		}
	}
}

func routeReadyEvents(ctx context.Context, done chan<- struct{}, buildInfo models.BuildInformation,
	tunnelReadyCh <-chan struct{},
	unboundLooper dns.Looper, updaterLooper updater.Looper, publicIPLooper publicip.Looper,
	routing routing.Routing, logger logging.Logger, httpClient *http.Client,
	versionInformation, portForwardingEnabled bool, startPortForward func(vpnGateway net.IP)) {
	defer close(done)

	// for linters only
	var restartTickerContext context.Context
	var restartTickerCancel context.CancelFunc = func() {}

	unboundTickerDone := make(chan struct{})
	close(unboundTickerDone)
	updaterTickerDone := make(chan struct{})
	close(updaterTickerDone)

	first := true
	for {
		select {
		case <-ctx.Done():
			restartTickerCancel() // for linters only
			<-unboundTickerDone
			<-updaterTickerDone
			return
		case <-tunnelReadyCh: // blocks until openvpn is connected
			vpnDestination, err := routing.VPNDestinationIP()
			if err != nil {
				logger.Warn(err)
			} else {
				logger.Info("VPN routing IP address: %s", vpnDestination)
			}

			if unboundLooper.GetSettings().Enabled {
				_, _ = unboundLooper.SetStatus(constants.Running)
			}

			restartTickerCancel() // stop previous restart tickers
			<-unboundTickerDone
			<-updaterTickerDone
			restartTickerContext, restartTickerCancel = context.WithCancel(ctx)

			// Runs the Public IP getter job once
			_, _ = publicIPLooper.SetStatus(constants.Running)
			if !versionInformation {
				break
			}

			if first {
				first = false
				message, err := versionpkg.GetMessage(ctx, buildInfo, httpClient)
				if err != nil {
					logger.Error("cannot get version information: " + err.Error())
				} else {
					logger.Info(message)
				}
			}

			unboundTickerDone = make(chan struct{})
			updaterTickerDone = make(chan struct{})
			go unboundLooper.RunRestartTicker(restartTickerContext, unboundTickerDone)
			go updaterLooper.RunRestartTicker(restartTickerContext, updaterTickerDone)
			if portForwardingEnabled {
				// vpnGateway required only for PIA
				vpnGateway, err := routing.VPNLocalGatewayIP()
				if err != nil {
					logger.Error(err)
				}
				logger.Info("VPN gateway IP address: %s", vpnGateway)
				startPortForward(vpnGateway)
			}
		}
	}
}
