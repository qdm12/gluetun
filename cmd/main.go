package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	libhealthcheck "github.com/qdm12/golibs/healthcheck"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/alpine"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/dns"
	"github.com/qdm12/private-internet-access-docker/internal/firewall"
	"github.com/qdm12/private-internet-access-docker/internal/healthcheck"
	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/qdm12/private-internet-access-docker/internal/mullvad"
	"github.com/qdm12/private-internet-access-docker/internal/openvpn"
	"github.com/qdm12/private-internet-access-docker/internal/params"
	"github.com/qdm12/private-internet-access-docker/internal/pia"
	"github.com/qdm12/private-internet-access-docker/internal/routing"
	"github.com/qdm12/private-internet-access-docker/internal/server"
	"github.com/qdm12/private-internet-access-docker/internal/settings"
	"github.com/qdm12/private-internet-access-docker/internal/shadowsocks"
	"github.com/qdm12/private-internet-access-docker/internal/splash"
	"github.com/qdm12/private-internet-access-docker/internal/tinyproxy"
	"github.com/qdm12/private-internet-access-docker/internal/windscribe"
)

func main() {
	if libhealthcheck.Mode(os.Args) {
		if err := healthcheck.HealthCheck(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := createLogger()
	fatalOnError := makeFatalOnError(logger, cancel)
	paramsReader := params.NewReader(logger)
	fmt.Println(splash.Splash(
		paramsReader.GetVersion(),
		paramsReader.GetVcsRef(),
		paramsReader.GetBuildDate()))

	client := network.NewClient(15 * time.Second)
	// Create configurators
	fileManager := files.NewFileManager()
	alpineConf := alpine.NewConfigurator(fileManager)
	ovpnConf := openvpn.NewConfigurator(logger, fileManager)
	dnsConf := dns.NewConfigurator(logger, client, fileManager)
	firewallConf := firewall.NewConfigurator(logger)
	routingConf := routing.NewRouting(logger, fileManager)
	piaConf := pia.NewConfigurator(client, fileManager, firewallConf)
	mullvadConf := mullvad.NewConfigurator(fileManager, logger)
	windscribeConf := windscribe.NewConfigurator(fileManager)
	tinyProxyConf := tinyproxy.NewConfigurator(fileManager, logger)
	shadowsocksConf := shadowsocks.NewConfigurator(fileManager, logger)
	streamMerger := command.NewStreamMerger()

	printVersions(ctx, logger, map[string]func(ctx context.Context) (string, error){
		"OpenVPN":     ovpnConf.Version,
		"Unbound":     dnsConf.Version,
		"IPtables":    firewallConf.Version,
		"TinyProxy":   tinyProxyConf.Version,
		"ShadowSocks": shadowsocksConf.Version,
	})

	allSettings, err := settings.GetAllSettings(paramsReader)
	fatalOnError(err)
	logger.Info(allSettings.String())

	err = alpineConf.CreateUser("nonrootuser", allSettings.System.UID)
	fatalOnError(err)
	err = fileManager.SetOwnership("/etc/unbound", allSettings.System.UID, allSettings.System.GID)
	fatalOnError(err)
	err = fileManager.SetOwnership("/etc/tinyproxy", allSettings.System.UID, allSettings.System.GID)
	fatalOnError(err)

	if err := ovpnConf.CheckTUN(); err != nil {
		logger.Warn(err)
		err = ovpnConf.CreateTUN()
		fatalOnError(err)
	}

	var openVPNUser, openVPNPassword string
	switch allSettings.VPNSP {
	case constants.PrivateInternetAccess:
		openVPNUser = allSettings.PIA.User
		openVPNPassword = allSettings.PIA.Password
	case constants.Mullvad:
		openVPNUser = allSettings.Mullvad.User
		openVPNPassword = "m"
	case constants.Windscribe:
		openVPNUser = allSettings.Windscribe.User
		openVPNPassword = allSettings.Windscribe.Password
	}
	err = ovpnConf.WriteAuthFile(openVPNUser, openVPNPassword, allSettings.System.UID, allSettings.System.GID)
	fatalOnError(err)

	defaultInterface, defaultGateway, defaultSubnet, err := routingConf.DefaultRoute()
	fatalOnError(err)

	// Temporarily reset chain policies allowing Kubernetes sidecar to
	// successfully restart the container. Without this, the existing rules will
	// pre-exist, preventing the nslookup of the PIA region address. These will
	// simply be redundant at Docker runtime as they will already be set this way
	// Thanks to @npawelek https://github.com/npawelek
	err = firewallConf.AcceptAll(ctx)
	fatalOnError(err)

	connected, signalConnected := context.WithCancel(context.Background())
	go collectStreamLines(ctx, streamMerger, logger, signalConnected)

	waiter := command.NewWaiter()

	var connections []models.OpenVPNConnection
	switch allSettings.VPNSP {
	case constants.PrivateInternetAccess:
		connections, err = piaConf.GetOpenVPNConnections(
			allSettings.PIA.Region,
			allSettings.OpenVPN.NetworkProtocol,
			allSettings.PIA.Encryption,
			allSettings.OpenVPN.TargetIP)
		if err != nil {
			break
		}
		err = piaConf.BuildConf(
			connections,
			allSettings.PIA.Encryption,
			allSettings.OpenVPN.Verbosity,
			allSettings.System.UID,
			allSettings.System.GID,
			allSettings.OpenVPN.Root,
			allSettings.OpenVPN.Cipher,
			allSettings.OpenVPN.Auth)
	case constants.Mullvad:
		connections, err = mullvadConf.GetOpenVPNConnections(
			allSettings.Mullvad.Country,
			allSettings.Mullvad.City,
			allSettings.Mullvad.ISP,
			allSettings.OpenVPN.NetworkProtocol,
			allSettings.Mullvad.Port,
			allSettings.OpenVPN.TargetIP)
		if err != nil {
			break
		}
		err = mullvadConf.BuildConf(
			connections,
			allSettings.OpenVPN.Verbosity,
			allSettings.System.UID,
			allSettings.System.GID,
			allSettings.OpenVPN.Root,
			allSettings.OpenVPN.Cipher)
	case constants.Windscribe:
		connections, err = windscribeConf.GetOpenVPNConnections(
			allSettings.Windscribe.Region,
			allSettings.OpenVPN.NetworkProtocol,
			allSettings.Windscribe.Port,
			allSettings.OpenVPN.TargetIP)
		if err != nil {
			break
		}
		err = windscribeConf.BuildConf(
			connections,
			allSettings.OpenVPN.Verbosity,
			allSettings.System.UID,
			allSettings.System.GID,
			allSettings.OpenVPN.Root,
			allSettings.OpenVPN.Cipher,
			allSettings.OpenVPN.Auth)
	}
	fatalOnError(err)

	err = routingConf.AddRoutesVia(ctx, allSettings.Firewall.AllowedSubnets, defaultGateway, defaultInterface)
	fatalOnError(err)
	err = firewallConf.Clear(ctx)
	fatalOnError(err)
	err = firewallConf.BlockAll(ctx)
	fatalOnError(err)
	err = firewallConf.CreateGeneralRules(ctx)
	fatalOnError(err)
	err = firewallConf.CreateVPNRules(ctx, constants.TUN, defaultInterface, connections)
	fatalOnError(err)
	err = firewallConf.CreateLocalSubnetsRules(ctx, defaultSubnet, allSettings.Firewall.AllowedSubnets, defaultInterface)
	fatalOnError(err)

	if allSettings.TinyProxy.Enabled {
		err = tinyProxyConf.MakeConf(
			allSettings.TinyProxy.LogLevel,
			allSettings.TinyProxy.Port,
			allSettings.TinyProxy.User,
			allSettings.TinyProxy.Password,
			allSettings.System.UID,
			allSettings.System.GID)
		fatalOnError(err)
		err = firewallConf.AllowAnyIncomingOnPort(ctx, allSettings.TinyProxy.Port)
		fatalOnError(err)
		stream, waitFn, err := tinyProxyConf.Start(ctx)
		fatalOnError(err)
		waiter.Add(func() error {
			err := waitFn()
			logger.Error("tinyproxy: %s", err)
			return err
		})
		go streamMerger.Merge(ctx, stream, command.MergeName("tinyproxy"), command.MergeColor(constants.ColorTinyproxy()))
	}

	if allSettings.ShadowSocks.Enabled {
		err = shadowsocksConf.MakeConf(
			allSettings.ShadowSocks.Port,
			allSettings.ShadowSocks.Password,
			allSettings.ShadowSocks.Method,
			allSettings.System.UID,
			allSettings.System.GID)
		fatalOnError(err)
		err = firewallConf.AllowAnyIncomingOnPort(ctx, allSettings.ShadowSocks.Port)
		fatalOnError(err)
		stdout, stderr, waitFn, err := shadowsocksConf.Start(ctx, "0.0.0.0", allSettings.ShadowSocks.Port, allSettings.ShadowSocks.Password, allSettings.ShadowSocks.Log)
		fatalOnError(err)
		waiter.Add(func() error {
			err := waitFn()
			logger.Error("shadowsocks: %s", err)
			return err
		})
		go streamMerger.Merge(ctx, stdout, command.MergeName("shadowsocks"), command.MergeColor(constants.ColorShadowsocks()))
		go streamMerger.Merge(ctx, stderr, command.MergeName("shadowsocks error"), command.MergeColor(constants.ColorShadowsocksError()))
	}

	httpServer := server.New("0.0.0.0:8000", logger)

	go openvpnRunLoop(ctx, ovpnConf, streamMerger, logger, httpServer, waiter, fatalOnError)

	waiter.Add(func() error {
		err := httpServer.Run(ctx)
		logger.Error("http server: %s", err)
		return err
	})

	go func() {
		<-connected.Done() // blocks until openvpn is connected
		onConnected(ctx, allSettings, logger, dnsConf, fileManager, waiter,
			streamMerger, httpServer, routingConf, defaultInterface, piaConf)
	}()

	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)
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
	}
	if allSettings.PIA.PortForwarding.Enabled {
		logger.Info("Clearing forwarded port status file %s", allSettings.PIA.PortForwarding.Filepath)
		if err := fileManager.Remove(string(allSettings.PIA.PortForwarding.Filepath)); err != nil {
			logger.Error(err)
		}
	}
	errors := waiter.WaitForAll(ctx)
	for _, err := range errors {
		logger.Error(err)
	}
}

func makeFatalOnError(logger logging.Logger, cancel func()) func(err error) {
	return func(err error) {
		if err != nil {
			logger.Error(err)
			cancel()
			time.Sleep(100 * time.Millisecond) // wait for operations to terminate
			os.Exit(1)
		}
	}
}

func createLogger() logging.Logger {
	logger, err := logging.NewLogger(logging.ConsoleEncoding, logging.InfoLevel, -1)
	if err != nil {
		panic(err)
	}
	return logger
}

func printVersions(ctx context.Context, logger logging.Logger, versionFunctions map[string]func(ctx context.Context) (string, error)) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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

func collectStreamLines(ctx context.Context, streamMerger command.StreamMerger, logger logging.Logger, signalConnected func()) {
	// Blocking line merging paramsReader for all programs: openvpn, tinyproxy, unbound and shadowsocks
	logger.Info("Launching standard output merger")
	streamMerger.CollectLines(ctx, func(line string) {
		logger.Info(line)
		if strings.Contains(line, "Initialization Sequence Completed") {
			signalConnected()
		}
	}, func(err error) {
		logger.Error(err)
	})
}

func openvpnRunLoop(ctx context.Context, ovpnConf openvpn.Configurator, streamMerger command.StreamMerger,
	logger logging.Logger, httpServer server.Server, waiter command.Waiter, fatalOnError func(err error)) {
	waitErrors := make(chan error)
	for {
		if ctx.Err() == context.Canceled {
			break
		}
		logger.Info("openvpn: starting")
		openvpnCtx, openvpnCancel := context.WithCancel(ctx)
		stream, waitFn, err := ovpnConf.Start(openvpnCtx)
		fatalOnError(err)
		httpServer.SetOpenVPNRestart(openvpnCancel)
		go streamMerger.Merge(openvpnCtx, stream, command.MergeName("openvpn"), command.MergeColor(constants.ColorOpenvpn()))
		waiter.Add(func() error {
			return <-waitErrors
		})
		err = waitFn()
		waitErrors <- err
		if openvpnCtx.Err() == context.Canceled {
			logger.Info("openvpn: shutting down")
		}
		logger.Error("openvpn: %s", err)
		openvpnCancel()
	}
}

func onConnected(ctx context.Context, allSettings settings.Settings,
	logger logging.Logger, dnsConf dns.Configurator, fileManager files.FileManager,
	waiter command.Waiter, streamMerger command.StreamMerger, httpServer server.Server,
	routingConf routing.Routing, defaultInterface string,
	piaConf pia.Configurator,
) {
	if allSettings.PIA.PortForwarding.Enabled {
		time.AfterFunc(5*time.Second, func() {
			setupPortForwarding(logger, piaConf, allSettings.PIA, allSettings.System.UID, allSettings.System.GID)
		})
	}

	if allSettings.DNS.Enabled {
		go unboundRunLoop(ctx, logger, dnsConf, allSettings.DNS, allSettings.System.UID, allSettings.System.GID, waiter, streamMerger, httpServer)
	}

	ip, err := routingConf.CurrentPublicIP(defaultInterface)
	if err != nil {
		logger.Error(err)
	} else {
		logger.Info("Tunnel IP is %s, see more information at https://ipinfo.io/%s", ip, ip)
		err = fileManager.WriteLinesToFile(
			string(allSettings.System.IPStatusFilepath),
			[]string{ip.String()},
			files.Ownership(allSettings.System.UID, allSettings.System.GID),
			files.Permissions(0400))
		if err != nil {
			logger.Error(err)
		}
	}
}

func fallbackToUnencryptedDNS(dnsConf dns.Configurator, provider models.DNSProvider, ipv6 bool) error {
	targetDNS := constants.DNSProviderMapping()[provider]
	var targetIP net.IP
	for _, targetIP = range targetDNS.IPs {
		if ipv6 && targetIP.To4() == nil {
			break
		} else if !ipv6 && targetIP.To4() != nil {
			break
		}
	}
	dnsConf.UseDNSInternally(targetIP)
	return dnsConf.UseDNSSystemWide(targetIP)
}

func unboundRun(ctx, unboundCtx context.Context, unboundCancel context.CancelFunc, dnsConf dns.Configurator, settings settings.DNS, uid, gid int,
	streamMerger command.StreamMerger, waiter command.Waiter, httpServer server.Server) (newCtx context.Context, newCancel context.CancelFunc, err error) {
	if err := dnsConf.DownloadRootHints(uid, gid); err != nil {
		return unboundCtx, unboundCancel, err
	}
	if err := dnsConf.DownloadRootKey(uid, gid); err != nil {
		return unboundCtx, unboundCancel, err
	}
	if err := dnsConf.MakeUnboundConf(settings, uid, gid); err != nil {
		return unboundCtx, unboundCancel, err
	}
	unboundCancel()
	if settings.UpdatePeriod > 0 {
		newCtx, newCancel = context.WithTimeout(ctx, settings.UpdatePeriod)
	} else {
		newCtx, newCancel = context.WithCancel(ctx)
	}
	stream, waitFn, err := dnsConf.Start(newCtx, settings.VerbosityDetailsLevel)
	if err != nil {
		newCancel()
		if fallbackErr := fallbackToUnencryptedDNS(dnsConf, settings.Providers[0], settings.IPv6); err != nil {
			return newCtx, newCancel, fmt.Errorf("%s: %w", err, fallbackErr)
		}
		return newCtx, newCancel, err
	}
	go streamMerger.Merge(newCtx, stream, command.MergeName("unbound"), command.MergeColor(constants.ColorUnbound()))
	dnsConf.UseDNSInternally(net.IP{127, 0, 0, 1})                         // use Unbound
	if err := dnsConf.UseDNSSystemWide(net.IP{127, 0, 0, 1}); err != nil { // use Unbound
		newCancel()
		if fallbackErr := fallbackToUnencryptedDNS(dnsConf, settings.Providers[0], settings.IPv6); err != nil {
			return newCtx, newCancel, fmt.Errorf("%s: %w", err, fallbackErr)
		}
		return newCtx, newCancel, err
	}
	if err := dnsConf.WaitForUnbound(); err != nil {
		newCancel()
		if fallbackErr := fallbackToUnencryptedDNS(dnsConf, settings.Providers[0], settings.IPv6); err != nil {
			return newCtx, newCancel, fmt.Errorf("%s: %w", err, fallbackErr)
		}
		return newCtx, newCancel, err
	}
	// Unbound is up and running at this point
	httpServer.SetUnboundRestart(newCancel)
	waitErrors := make(chan error)
	waiter.Add(func() error { //nolint:scopelint
		return <-waitErrors
	})
	err = waitFn()
	waitErrors <- err
	if newCtx.Err() == context.Canceled || newCtx.Err() == context.DeadlineExceeded {
		return newCtx, newCancel, nil
	}
	return newCtx, newCancel, err
}

func unboundRunLoop(ctx context.Context, logger logging.Logger, dnsConf dns.Configurator,
	settings settings.DNS, uid, gid int,
	waiter command.Waiter, streamMerger command.StreamMerger, httpServer server.Server,
) {
	logger = logger.WithPrefix("unbound dns over tls setup: ")
	if err := fallbackToUnencryptedDNS(dnsConf, settings.Providers[0], settings.IPv6); err != nil {
		logger.Error(err)
	}
	unboundCtx, unboundCancel := context.WithCancel(ctx)
	defer unboundCancel()
	for ctx.Err() == nil {
		var err error
		unboundCtx, unboundCancel, err = unboundRun(ctx, unboundCtx, unboundCancel, dnsConf, settings, uid, gid, streamMerger, waiter, httpServer)
		if err != nil {
			logger.Error(err)
			time.Sleep(10 * time.Second)
			continue
		}
		logger.Info("attempting restart")
	}
	logger.Info("shutting down")
}

func setupPortForwarding(logger logging.Logger, piaConf pia.Configurator, settings settings.PIA, uid, gid int) {
	pfLogger := logger.WithPrefix("port forwarding: ")
	var port uint16
	var err error
	for {
		port, err = piaConf.GetPortForward()
		if err != nil {
			pfLogger.Error(err)
			pfLogger.Info("retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		} else {
			pfLogger.Info("port forwarded is %d", port)
			break
		}
	}
	pfLogger.Info("writing forwarded port to %s", settings.PortForwarding.Filepath)
	if err := piaConf.WritePortForward(settings.PortForwarding.Filepath, port, uid, gid); err != nil {
		pfLogger.Error(err)
	}
	pfLogger.Info("allowing forwarded port %d through firewall", port)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := piaConf.AllowPortForwardFirewall(ctx, constants.TUN, port); err != nil {
		pfLogger.Error(err)
	}
}
