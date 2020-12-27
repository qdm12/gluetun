package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

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
	"github.com/qdm12/gluetun/internal/updater"
	versionpkg "github.com/qdm12/gluetun/internal/version"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
)

//nolint:gochecknoglobals
var (
	buildInfo models.BuildInformation
	version   = "unknown"
	commit    = "unknown"
	buildDate = "an unknown date"
)

func main() {
	buildInfo.Version = version
	buildInfo.Commit = commit
	buildInfo.BuildDate = buildDate
	ctx := context.Background()
	os.Exit(_main(ctx, os.Args))
}

func _main(background context.Context, args []string) int { //nolint:gocognit,gocyclo
	if len(args) > 1 { // cli operation
		var err error
		switch args[1] {
		case "healthcheck":
			err = cli.HealthCheck(background)
		case "clientkey":
			err = cli.ClientKey(args[2:])
		case "openvpnconfig":
			err = cli.OpenvpnConfig()
		case "update":
			err = cli.Update(args[2:])
		default:
			err = fmt.Errorf("command %q is unknown", args[1])
		}
		if err != nil {
			fmt.Println(err)
			return 1
		}
		return 0
	}
	ctx, cancel := context.WithCancel(background)
	defer cancel()
	logger := createLogger()

	const clientTimeout = 15 * time.Second
	httpClient := &http.Client{Timeout: clientTimeout}
	client := network.NewClient(clientTimeout)
	// Create configurators
	fileManager := files.NewFileManager()
	alpineConf := alpine.NewConfigurator(fileManager)
	ovpnConf := openvpn.NewConfigurator(logger, fileManager)
	dnsConf := dns.NewConfigurator(logger, client, fileManager)
	routingConf := routing.NewRouting(logger)
	firewallConf := firewall.NewConfigurator(logger, routingConf, fileManager)
	streamMerger := command.NewStreamMerger()

	paramsReader := params.NewReader(logger, fileManager)
	fmt.Println(gluetunLogging.Splash(buildInfo))

	printVersions(ctx, logger, map[string]func(ctx context.Context) (string, error){
		"OpenVPN":  ovpnConf.Version,
		"Unbound":  dnsConf.Version,
		"IPtables": firewallConf.Version,
	})

	allSettings, err := settings.GetAllSettings(paramsReader)
	if err != nil {
		logger.Error(err)
		return 1
	}
	logger.Info(allSettings.String())

	// TODO run this in a loop or in openvpn to reload from file without restarting
	storage := storage.New(logger)
	const updateServerFile = true
	allServers, err := storage.SyncServers(constants.GetAllServers(), updateServerFile)
	if err != nil {
		logger.Error(err)
		return 1
	}

	// Should never change
	uid, gid := allSettings.System.UID, allSettings.System.GID

	const defaultUsername = "nonrootuser"
	nonRootUsername, err := alpineConf.CreateUser(defaultUsername, uid)
	if err != nil {
		logger.Error(err)
		return 1
	}
	err = fileManager.SetOwnership("/etc/unbound", uid, gid)
	if err != nil {
		logger.Error(err)
		return 1
	}

	if allSettings.Firewall.Debug {
		firewallConf.SetDebug()
		routingConf.SetDebug()
	}

	defaultInterface, defaultGateway, err := routingConf.DefaultRoute()
	if err != nil {
		logger.Error(err)
		return 1
	}

	localSubnet, err := routingConf.LocalSubnet()
	if err != nil {
		logger.Error(err)
		return 1
	}

	defaultIP, err := routingConf.DefaultIP()
	if err != nil {
		logger.Error(err)
		return 1
	}

	firewallConf.SetNetworkInformation(defaultInterface, defaultGateway, localSubnet, defaultIP)

	if err := routingConf.Setup(); err != nil {
		logger.Error(err)
		return 1
	}
	defer func() {
		routingConf.SetVerbose(false)
		if err := routingConf.TearDown(); err != nil {
			logger.Error(err)
		}
	}()

	if err := firewallConf.SetOutboundSubnets(ctx, allSettings.Firewall.OutboundSubnets); err != nil {
		logger.Error(err)
		return 1
	}
	if err := routingConf.SetOutboundRoutes(allSettings.Firewall.OutboundSubnets); err != nil {
		logger.Error(err)
		return 1
	}

	if err := ovpnConf.CheckTUN(); err != nil {
		logger.Warn(err)
		err = ovpnConf.CreateTUN()
		if err != nil {
			logger.Error(err)
			return 1
		}
	}

	tunnelReadyCh, dnsReadyCh := make(chan struct{}), make(chan struct{})
	signalTunnelReady := func() { tunnelReadyCh <- struct{}{} }
	signalDNSReady := func() { dnsReadyCh <- struct{}{} }
	defer close(tunnelReadyCh)
	defer close(dnsReadyCh)

	if allSettings.Firewall.Enabled {
		err := firewallConf.SetEnabled(ctx, true) // disabled by default
		if err != nil {
			logger.Error(err)
			return 1
		}
	}

	for _, vpnPort := range allSettings.Firewall.VPNInputPorts {
		err = firewallConf.SetAllowedPort(ctx, vpnPort, string(constants.TUN))
		if err != nil {
			logger.Error(err)
			return 1
		}
	}

	for _, port := range allSettings.Firewall.InputPorts {
		err = firewallConf.SetAllowedPort(ctx, port, defaultInterface)
		if err != nil {
			logger.Error(err)
			return 1
		}
	} // TODO move inside firewall?

	wg := &sync.WaitGroup{}

	go collectStreamLines(ctx, streamMerger, logger, signalTunnelReady)

	openvpnLooper := openvpn.NewLooper(allSettings.OpenVPN, nonRootUsername, uid, gid, allServers,
		ovpnConf, firewallConf, routingConf, logger, httpClient, fileManager, streamMerger, cancel)
	wg.Add(1)
	// wait for restartOpenvpn
	go openvpnLooper.Run(ctx, wg)

	updaterLooper := updater.NewLooper(allSettings.Updater,
		allServers, storage, openvpnLooper.SetServers, httpClient, logger)
	wg.Add(1)
	// wait for updaterLooper.Restart() or its ticket launched with RunRestartTicker
	go updaterLooper.Run(ctx, wg)

	unboundLooper := dns.NewLooper(dnsConf, allSettings.DNS, logger, streamMerger, nonRootUsername, uid, gid)
	wg.Add(1)
	// wait for unboundLooper.Restart or its ticker launched with RunRestartTicker
	go unboundLooper.Run(ctx, wg, signalDNSReady)

	publicIPLooper := publicip.NewLooper(client, logger, fileManager,
		allSettings.System.IPStatusFilepath, allSettings.PublicIPPeriod, uid, gid)
	wg.Add(1)
	go publicIPLooper.Run(ctx, wg)
	wg.Add(1)
	go publicIPLooper.RunRestartTicker(ctx, wg)
	publicIPLooper.SetPeriod(allSettings.PublicIPPeriod) // call after RunRestartTicker

	httpProxyLooper := httpproxy.NewLooper(logger, allSettings.HTTPProxy)
	wg.Add(1)
	go httpProxyLooper.Run(ctx, wg)

	shadowsocksLooper := shadowsocks.NewLooper(allSettings.ShadowSocks, logger, defaultInterface)
	restartShadowsocks := shadowsocksLooper.Restart
	wg.Add(1)
	go shadowsocksLooper.Run(ctx, wg)

	if allSettings.HTTPProxy.Enabled {
		httpProxyLooper.Restart()
	}
	if allSettings.ShadowSocks.Enabled {
		restartShadowsocks()
	}

	wg.Add(1)
	go routeReadyEvents(ctx, wg, tunnelReadyCh, dnsReadyCh,
		unboundLooper, updaterLooper, publicIPLooper, routingConf, logger, httpClient,
		allSettings.VersionInformation, allSettings.OpenVPN.Provider.PortForwarding.Enabled, openvpnLooper.PortForward,
	)
	controlServerAddress := fmt.Sprintf("0.0.0.0:%d", allSettings.ControlServer.Port)
	controlServerLogging := allSettings.ControlServer.Log
	httpServer := server.New(controlServerAddress, controlServerLogging,
		logger, buildInfo, openvpnLooper, unboundLooper, updaterLooper)
	wg.Add(1)
	go httpServer.Run(ctx, wg)

	healthcheckServer := healthcheck.NewServer(
		constants.HealthcheckAddress, logger)
	wg.Add(1)
	go healthcheckServer.Run(ctx, wg)

	// Start openvpn for the first time in a blocking call
	// until openvpn is launched
	_, _ = openvpnLooper.SetStatus(constants.Running) // TODO option to disable with variable

	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)
	shutdownErrorsCount := 0
	select {
	case signal := <-signalsCh:
		logger.Warn("Caught OS signal %s, shutting down", signal)
		cancel()
	case <-ctx.Done():
		logger.Warn("context canceled, shutting down")
	}
	logger.Info("Clearing ip status file %s", allSettings.System.IPStatusFilepath)
	if err := fileManager.Remove(string(allSettings.System.IPStatusFilepath)); err != nil {
		logger.Error(err)
		shutdownErrorsCount++
	}
	if allSettings.OpenVPN.Provider.PortForwarding.Enabled {
		logger.Info("Clearing forwarded port status file %s", allSettings.OpenVPN.Provider.PortForwarding.Filepath)
		if err := fileManager.Remove(string(allSettings.OpenVPN.Provider.PortForwarding.Filepath)); err != nil {
			logger.Error(err)
			shutdownErrorsCount++
		}
	}
	const shutdownGracePeriod = 5 * time.Second
	waiting, waited := context.WithTimeout(context.Background(), shutdownGracePeriod)
	go func() {
		defer waited()
		wg.Wait()
	}()
	<-waiting.Done()
	if waiting.Err() == context.DeadlineExceeded {
		if shutdownErrorsCount > 0 {
			logger.Warn("Shutdown had %d errors", shutdownErrorsCount)
		}
		logger.Warn("Shutdown timed out")
		return 1
	}
	if shutdownErrorsCount > 0 {
		logger.Warn("Shutdown had %d errors")
		return 1
	}
	logger.Info("Shutdown successful")
	return 0
}

func createLogger() logging.Logger {
	logger, err := logging.NewLogger(logging.ConsoleEncoding, logging.InfoLevel)
	if err != nil {
		panic(err)
	}
	return logger
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

//nolint:lll
func collectStreamLines(ctx context.Context, streamMerger command.StreamMerger,
	logger logging.Logger, signalTunnelReady func()) {
	// Blocking line merging paramsReader for openvpn and unbound
	logger.Info("Launching standard output merger")
	streamMerger.CollectLines(ctx, func(line string) {
		line, level := gluetunLogging.PostProcessLine(line)
		if line == "" {
			return
		}
		switch level {
		case logging.DebugLevel:
			logger.Debug(line)
		case logging.InfoLevel:
			logger.Info(line)
		case logging.WarnLevel:
			logger.Warn(line)
		case logging.ErrorLevel:
			logger.Error(line)
		}
		switch {
		case strings.Contains(line, "Initialization Sequence Completed"):
			signalTunnelReady()
		case strings.Contains(line, "TLS Error: TLS key negotiation failed to occur within 60 seconds (check your network connectivity)"):
			logger.Warn("This means that either...")
			logger.Warn("1. The VPN server IP address you are trying to connect to is no longer valid, see https://github.com/qdm12/gluetun/wiki/Update-servers-information")
			logger.Warn("2. The VPN server crashed, try changing region")
			logger.Warn("3. Your Internet connection is not working, ensure it works")
			logger.Warn("Feel free to create an issue at https://github.com/qdm12/gluetun/issues/new/choose")
		}
	}, func(err error) {
		logger.Warn(err)
	})
}

func routeReadyEvents(ctx context.Context, wg *sync.WaitGroup, tunnelReadyCh, dnsReadyCh <-chan struct{},
	unboundLooper dns.Looper, updaterLooper updater.Looper, publicIPLooper publicip.Looper,
	routing routing.Routing, logger logging.Logger, httpClient *http.Client,
	versionInformation, portForwardingEnabled bool, startPortForward func(vpnGateway net.IP)) {
	defer wg.Done()
	tickerWg := &sync.WaitGroup{}
	// for linters only
	var restartTickerContext context.Context
	var restartTickerCancel context.CancelFunc = func() {}
	for {
		select {
		case <-ctx.Done():
			restartTickerCancel() // for linters only
			tickerWg.Wait()
			return
		case <-tunnelReadyCh: // blocks until openvpn is connected
			_, _ = unboundLooper.SetStatus(constants.Running)
			restartTickerCancel() // stop previous restart tickers
			tickerWg.Wait()
			restartTickerContext, restartTickerCancel = context.WithCancel(ctx)
			tickerWg.Add(2) //nolint:gomnd
			go unboundLooper.RunRestartTicker(restartTickerContext, tickerWg)
			go updaterLooper.RunRestartTicker(restartTickerContext, tickerWg)
			vpnDestination, err := routing.VPNDestinationIP()
			if err != nil {
				logger.Warn(err)
			} else {
				logger.Info("VPN routing IP address: %s", vpnDestination)
			}
			if portForwardingEnabled {
				// vpnGateway required only for PIA
				vpnGateway, err := routing.VPNLocalGatewayIP()
				if err != nil {
					logger.Error(err)
				}
				logger.Info("VPN gateway IP address: %s", vpnGateway)
				startPortForward(vpnGateway)
			}
		case <-dnsReadyCh:
			publicIPLooper.Restart() // TODO do not restart if disabled
			if !versionInformation {
				break
			}
			message, err := versionpkg.GetMessage(ctx, buildInfo, httpClient)
			if err != nil {
				logger.Error(err)
				break
			}
			logger.Info(message)
		}
	}
}
