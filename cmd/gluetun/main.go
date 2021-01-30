package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	nativeos "os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/qdm12/dns/pkg/unbound"
	"github.com/qdm12/gluetun/internal/alpine"
	"github.com/qdm12/gluetun/internal/cli"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/dns"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/healthcheck"
	"github.com/qdm12/gluetun/internal/httpproxy"
	gluetunLogging "github.com/qdm12/gluetun/internal/logging"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/openvpn"
	"github.com/qdm12/gluetun/internal/params"
	"github.com/qdm12/gluetun/internal/publicip"
	"github.com/qdm12/gluetun/internal/routing"
	"github.com/qdm12/gluetun/internal/server"
	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/gluetun/internal/shadowsocks"
	"github.com/qdm12/gluetun/internal/storage"
	"github.com/qdm12/gluetun/internal/unix"
	"github.com/qdm12/gluetun/internal/updater"
	versionpkg "github.com/qdm12/gluetun/internal/version"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
	"github.com/qdm12/golibs/os/user"
	"github.com/qdm12/updated/pkg/dnscrypto"
)

//nolint:gochecknoglobals
var (
	version   = "unknown"
	commit    = "unknown"
	buildDate = "an unknown date"
)

func main() {
	buildInfo := models.BuildInformation{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	logger, err := logging.NewLogger(logging.ConsoleEncoding, logging.InfoLevel)
	if err != nil {
		fmt.Println(err)
		nativeos.Exit(1)
	}

	args := nativeos.Args
	os := os.New()
	osUser := user.New()
	unix := unix.New()
	cli := cli.New()

	errorCh := make(chan error)
	go func() {
		errorCh <- _main(ctx, buildInfo, args, logger, os, osUser, unix, cli)
	}()

	signalsCh := make(chan nativeos.Signal, 1)
	signal.Notify(signalsCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		nativeos.Interrupt,
	)

	select {
	case signal := <-signalsCh:
		logger.Warn("Caught OS signal %s, shutting down", signal)
	case err := <-errorCh:
		close(errorCh)
		if err == nil { // expected exit such as healthcheck
			nativeos.Exit(0)
		}
		logger.Error(err)
	}

	cancel()

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

//nolint:gocognit,gocyclo
func _main(ctx context.Context, buildInfo models.BuildInformation,
	args []string, logger logging.Logger, os os.OS, osUser user.OSUser, unix unix.Unix,
	cli cli.CLI) error {
	if len(args) > 1 { // cli operation
		switch args[1] {
		case "healthcheck":
			return cli.HealthCheck(ctx)
		case "clientkey":
			return cli.ClientKey(args[2:], os.OpenFile)
		case "openvpnconfig":
			return cli.OpenvpnConfig(os)
		case "update":
			return cli.Update(ctx, args[2:], os)
		default:
			return fmt.Errorf("command %q is unknown", args[1])
		}
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	const clientTimeout = 15 * time.Second
	httpClient := &http.Client{Timeout: clientTimeout}
	// Create configurators
	alpineConf := alpine.NewConfigurator(os.OpenFile, osUser)
	ovpnConf := openvpn.NewConfigurator(logger, os, unix)
	dnsCrypto := dnscrypto.New(httpClient, "", "")
	const cacertsPath = "/etc/ssl/certs/ca-certificates.crt"
	dnsConf := unbound.NewConfigurator(logger, os.OpenFile, dnsCrypto,
		"/etc/unbound", "/usr/sbin/unbound", cacertsPath)
	routingConf := routing.NewRouting(logger)
	firewallConf := firewall.NewConfigurator(logger, routingConf, os.OpenFile)

	paramsReader := params.NewReader(logger, os)
	fmt.Println(gluetunLogging.Splash(buildInfo))

	printVersions(ctx, logger, map[string]func(ctx context.Context) (string, error){
		"OpenVPN":  ovpnConf.Version,
		"Unbound":  dnsConf.Version,
		"IPtables": firewallConf.Version,
	})

	allSettings, warnings, err := settings.GetAllSettings(paramsReader)
	for _, warning := range warnings {
		logger.Warn(warning)
	}
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

	// TODO run this in a loop or in openvpn to reload from file without restarting
	storage := storage.New(logger, os, constants.ServersData)
	allServers, err := storage.SyncServers(constants.GetAllServers())
	if err != nil {
		return err
	}

	// Should never change
	puid, pgid := allSettings.System.PUID, allSettings.System.PGID

	const defaultUsername = "nonrootuser"
	nonRootUsername, err := alpineConf.CreateUser(defaultUsername, puid)
	if err != nil {
		return err
	}
	if nonRootUsername != defaultUsername {
		logger.Info("using existing username %s corresponding to user id %d", nonRootUsername, puid)
	}

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

	localSubnet, err := routingConf.LocalSubnet()
	if err != nil {
		return err
	}

	defaultIP, err := routingConf.DefaultIP()
	if err != nil {
		return err
	}

	firewallConf.SetNetworkInformation(defaultInterface, defaultGateway, localSubnet, defaultIP)

	if err := routingConf.Setup(); err != nil {
		return err
	}
	defer func() {
		routingConf.SetVerbose(false)
		if err := routingConf.TearDown(); err != nil {
			logger.Error(err)
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

	wg := &sync.WaitGroup{}

	openvpnLooper := openvpn.NewLooper(allSettings.OpenVPN, nonRootUsername, puid, pgid, allServers,
		ovpnConf, firewallConf, routingConf, logger, httpClient, os.OpenFile, tunnelReadyCh, cancel)
	wg.Add(1)
	// wait for restartOpenvpn
	go openvpnLooper.Run(ctx, wg)

	updaterLooper := updater.NewLooper(allSettings.Updater,
		allServers, storage, openvpnLooper.SetServers, httpClient, logger)
	wg.Add(1)
	// wait for updaterLooper.Restart() or its ticket launched with RunRestartTicker
	go updaterLooper.Run(ctx, wg)

	unboundLooper := dns.NewLooper(dnsConf, allSettings.DNS, httpClient,
		logger, nonRootUsername, puid, pgid)
	wg.Add(1)
	// wait for unboundLooper.Restart or its ticker launched with RunRestartTicker
	go unboundLooper.Run(ctx, wg)

	publicIPLooper := publicip.NewLooper(
		httpClient, logger, allSettings.PublicIP, puid, pgid, os)
	wg.Add(1)
	go publicIPLooper.Run(ctx, wg)
	wg.Add(1)
	go publicIPLooper.RunRestartTicker(ctx, wg)

	httpProxyLooper := httpproxy.NewLooper(logger, allSettings.HTTPProxy)
	wg.Add(1)
	go httpProxyLooper.Run(ctx, wg)

	shadowsocksLooper := shadowsocks.NewLooper(allSettings.ShadowSocks, logger)
	wg.Add(1)
	go shadowsocksLooper.Run(ctx, wg)

	wg.Add(1)
	go routeReadyEvents(ctx, wg, buildInfo, tunnelReadyCh,
		unboundLooper, updaterLooper, publicIPLooper, routingConf, logger, httpClient,
		allSettings.VersionInformation, allSettings.OpenVPN.Provider.PortForwarding.Enabled, openvpnLooper.PortForward,
	)
	controlServerAddress := fmt.Sprintf("0.0.0.0:%d", allSettings.ControlServer.Port)
	controlServerLogging := allSettings.ControlServer.Log
	httpServer := server.New(controlServerAddress, controlServerLogging,
		logger, buildInfo, openvpnLooper, unboundLooper, updaterLooper, publicIPLooper)
	wg.Add(1)
	go httpServer.Run(ctx, wg)

	healthcheckServer := healthcheck.NewServer(
		constants.HealthcheckAddress, logger)
	wg.Add(1)
	go healthcheckServer.Run(ctx, wg)

	// Start openvpn for the first time in a blocking call
	// until openvpn is launched
	_, _ = openvpnLooper.SetStatus(constants.Running) // TODO option to disable with variable

	<-ctx.Done()

	if allSettings.OpenVPN.Provider.PortForwarding.Enabled {
		logger.Info("Clearing forwarded port status file %s", allSettings.OpenVPN.Provider.PortForwarding.Filepath)
		if err := os.Remove(string(allSettings.OpenVPN.Provider.PortForwarding.Filepath)); err != nil {
			logger.Error(err)
		}
	}

	wg.Wait()

	return nil
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

func routeReadyEvents(ctx context.Context, wg *sync.WaitGroup, buildInfo models.BuildInformation,
	tunnelReadyCh <-chan struct{},
	unboundLooper dns.Looper, updaterLooper updater.Looper, publicIPLooper publicip.Looper,
	routing routing.Routing, logger logging.Logger, httpClient *http.Client,
	versionInformation, portForwardingEnabled bool, startPortForward func(vpnGateway net.IP)) {
	defer wg.Done()
	tickerWg := &sync.WaitGroup{}
	// for linters only
	var restartTickerContext context.Context
	var restartTickerCancel context.CancelFunc = func() {}
	first := true
	for {
		select {
		case <-ctx.Done():
			restartTickerCancel() // for linters only
			tickerWg.Wait()
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
			tickerWg.Wait()
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
					logger.Error(err)
				} else {
					logger.Info(message)
				}
			}

			//nolint:gomnd
			tickerWg.Add(2)
			go unboundLooper.RunRestartTicker(restartTickerContext, tickerWg)
			go updaterLooper.RunRestartTicker(restartTickerContext, tickerWg)
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
