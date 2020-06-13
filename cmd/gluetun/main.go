package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/alpine"
	"github.com/qdm12/private-internet-access-docker/internal/cli"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/dns"
	"github.com/qdm12/private-internet-access-docker/internal/firewall"
	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/qdm12/private-internet-access-docker/internal/openvpn"
	"github.com/qdm12/private-internet-access-docker/internal/params"
	"github.com/qdm12/private-internet-access-docker/internal/provider"
	"github.com/qdm12/private-internet-access-docker/internal/publicip"
	"github.com/qdm12/private-internet-access-docker/internal/routing"
	"github.com/qdm12/private-internet-access-docker/internal/server"
	"github.com/qdm12/private-internet-access-docker/internal/settings"
	"github.com/qdm12/private-internet-access-docker/internal/shadowsocks"
	"github.com/qdm12/private-internet-access-docker/internal/splash"
	"github.com/qdm12/private-internet-access-docker/internal/tinyproxy"
)

func main() {
	ctx := context.Background()
	os.Exit(_main(ctx, os.Args))
}

func _main(background context.Context, args []string) int {
	if len(args) > 1 { // cli operation
		var err error
		switch args[1] {
		case "healthcheck":
			err = cli.HealthCheck()
		case "clientkey":
			err = cli.ClientKey(args[2:])
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

	providerConf := provider.New(allSettings.VPNSP, logger, client, fileManager, firewallConf)

	if !allSettings.Firewall.Enabled {
		firewallConf.Disable()
	}

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

	err = ovpnConf.WriteAuthFile(
		allSettings.OpenVPN.User,
		allSettings.OpenVPN.Password,
		allSettings.System.UID,
		allSettings.System.GID)
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

	connectedCh := make(chan struct{})
	signalConnected := func() {
		connectedCh <- struct{}{}
	}
	defer close(connectedCh)
	go collectStreamLines(ctx, streamMerger, logger, signalConnected)

	waiter := command.NewWaiter()

	connections, err := providerConf.GetOpenVPNConnections(allSettings.Provider.ServerSelection)
	fatalOnError(err)
	err = providerConf.BuildConf(
		connections,
		allSettings.OpenVPN.Verbosity,
		allSettings.System.UID,
		allSettings.System.GID,
		allSettings.OpenVPN.Root,
		allSettings.OpenVPN.Cipher,
		allSettings.OpenVPN.Auth,
		allSettings.Provider.ExtraConfigOptions,
	)
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
	err = firewallConf.RunUserPostRules(ctx, fileManager, "/iptables/post-rules.txt")
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
		nameserver := ""
		if allSettings.DNS.Enabled {
			nameserver = "127.0.0.1"
		}
		err = shadowsocksConf.MakeConf(
			allSettings.ShadowSocks.Port,
			allSettings.ShadowSocks.Password,
			allSettings.ShadowSocks.Method,
			nameserver,
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
		firstRun := true
		for {
			select {
			case <-ctx.Done():
				return
			case <-connectedCh: // blocks until openvpn is connected
				onConnected(ctx, allSettings, logger, dnsConf, fileManager, waiter,
					streamMerger, httpServer, routingConf, defaultInterface, providerConf, firstRun)
				firstRun = false
			}
		}
	}()

	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)
	exitStatus := 0
	select {
	case signal := <-signalsCh:
		logger.Warn("Caught OS signal %s, shutting down", signal)
		exitStatus = 1
		cancel()
	case <-ctx.Done():
		logger.Warn("context canceled, shutting down")
	}
	logger.Info("Clearing ip status file %s", allSettings.System.IPStatusFilepath)
	if err := fileManager.Remove(string(allSettings.System.IPStatusFilepath)); err != nil {
		logger.Error(err)
		exitStatus = 1
	}
	if allSettings.Provider.PortForwarding.Enabled {
		logger.Info("Clearing forwarded port status file %s", allSettings.Provider.PortForwarding.Filepath)
		if err := fileManager.Remove(string(allSettings.Provider.PortForwarding.Filepath)); err != nil {
			logger.Error(err)
			exitStatus = 1
		}
	}
	timeoutCtx, timeoutCancel := context.WithTimeout(background, time.Second)
	defer timeoutCancel()
	for _, err := range waiter.WaitForAll(timeoutCtx) {
		logger.Error(err)
		exitStatus = 1
	}
	return exitStatus
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
		line = trimEventualProgramPrefix(line)
		logger.Info(line)
		if strings.Contains(line, "Initialization Sequence Completed") {
			signalConnected()
		}
	}, func(err error) {
		logger.Warn(err)
	})
}

func trimEventualProgramPrefix(s string) string {
	switch {
	case strings.HasPrefix(s, "unbound: "):
		prefixRegex := regexp.MustCompile(`unbound: \[[0-9]{10}\] unbound\[[0-9]+:0\] `)
		prefix := prefixRegex.FindString(s)
		return fmt.Sprintf("unbound: %s", s[len(prefix):])
	case strings.HasPrefix(s, "shadowsocks: "):
		prefixRegex := regexp.MustCompile(`shadowsocks:[ ]+2[0-9]{3}\-[0-1][0-9]\-[0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9] `)
		prefix := prefixRegex.FindString(s)
		return fmt.Sprintf("shadowsocks: %s", s[len(prefix):])
	case strings.HasPrefix(s, "tinyproxy: "):
		logLevelRegex := regexp.MustCompile(`INFO|CONNECT|NOTICE|WARNING|ERROR|CRITICAL`)
		logLevel := logLevelRegex.FindString(s)
		prefixRegex := regexp.MustCompile(`tinyproxy: (INFO|CONNECT|NOTICE|WARNING|ERROR|CRITICAL)[ ]+(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) [0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9] \[[0-9]+\]: `)
		prefix := prefixRegex.FindString(s)
		return fmt.Sprintf("tinyproxy: %s %s", logLevel, s[len(prefix):])
	default:
		return s
	}
}

func openvpnRunLoop(ctx context.Context, ovpnConf openvpn.Configurator, streamMerger command.StreamMerger,
	logger logging.Logger, httpServer server.Server, waiter command.Waiter, fatalOnError func(err error)) {
	logger = logger.WithPrefix("openvpn: ")
	waitErrors := make(chan error)
	for ctx.Err() == nil {
		logger.Info("starting")
		openvpnCtx, openvpnCancel := context.WithCancel(ctx)
		stream, waitFn, err := ovpnConf.Start(openvpnCtx)
		fatalOnError(err)
		httpServer.SetOpenVPNRestart(openvpnCancel)
		go streamMerger.Merge(openvpnCtx, stream, command.MergeName("openvpn"), command.MergeColor(constants.ColorOpenvpn()))
		waiter.Add(func() error {
			return <-waitErrors
		})
		err = waitFn()
		waitErrors <- fmt.Errorf("openvpn: %w", err)
		switch {
		case ctx.Err() != nil:
			logger.Warn("context canceled: exiting openvpn run loop")
		case openvpnCtx.Err() == context.Canceled:
			logger.Info("triggered openvpn restart")
		default:
			logger.Warn(err)
			openvpnCancel()
		}
	}
}

func onConnected(ctx context.Context, allSettings settings.Settings,
	logger logging.Logger, dnsConf dns.Configurator, fileManager files.FileManager,
	waiter command.Waiter, streamMerger command.StreamMerger, httpServer server.Server,
	routingConf routing.Routing, defaultInterface string,
	providerConf provider.Provider, firstRun bool,
) {
	if allSettings.Provider.PortForwarding.Enabled {
		time.AfterFunc(5*time.Second, func() {
			setupPortForwarding(logger, providerConf, allSettings.Provider.PortForwarding.Filepath, allSettings.System.UID, allSettings.System.GID)
		})
	}
	if allSettings.DNS.Enabled && firstRun {
		go unboundRunLoop(ctx, logger, dnsConf, allSettings.DNS, allSettings.System.UID, allSettings.System.GID, waiter, streamMerger, httpServer)
	}

	vpnGatewayIP, err := routingConf.VPNGatewayIP(defaultInterface)
	if err != nil {
		logger.Warn(err)
	} else {
		logger.Info("Gateway VPN IP address: %s", vpnGatewayIP)
	}
	publicIP, err := publicip.NewIPGetter(network.NewClient(3 * time.Second)).Get()
	if err != nil {
		logger.Error(err)
	} else {
		logger.Info("Public IP address is %s", publicIP)
		err = fileManager.WriteLinesToFile(
			string(allSettings.System.IPStatusFilepath),
			[]string{publicIP.String()},
			files.Ownership(allSettings.System.UID, allSettings.System.GID),
			files.Permissions(0400))
		if err != nil {
			logger.Error(err)
		}
	}
}

func fallbackToUnencryptedIPv4DNS(dnsConf dns.Configurator, providers []models.DNSProvider) error {
	var targetIP net.IP
	for _, provider := range providers {
		data := constants.DNSProviderMapping()[provider]
		for _, targetIP = range data.IPs {
			if targetIP.To4() != nil {
				dnsConf.UseDNSInternally(targetIP)
				return dnsConf.UseDNSSystemWide(targetIP)
			}
		}
	}
	// No IPv4 address found
	return fmt.Errorf("no ipv4 DNS address found for providers %s", providers)
}

func unboundRun(ctx, oldCtx context.Context, oldCancel context.CancelFunc, timer *time.Timer,
	dnsConf dns.Configurator, settings settings.DNS, uid, gid int,
	streamMerger command.StreamMerger, waiter command.Waiter, httpServer server.Server) (
	newCtx context.Context, newCancel context.CancelFunc, setupErr, startErr, waitErr error) {
	if timer != nil {
		timer.Stop()
		timer.Reset(settings.UpdatePeriod)
	}
	if err := dnsConf.DownloadRootHints(uid, gid); err != nil {
		return oldCtx, oldCancel, err, nil, nil
	}
	if err := dnsConf.DownloadRootKey(uid, gid); err != nil {
		return oldCtx, oldCancel, err, nil, nil
	}
	if err := dnsConf.MakeUnboundConf(settings, uid, gid); err != nil {
		return oldCtx, oldCancel, err, nil, nil
	}
	newCtx, newCancel = context.WithCancel(ctx)
	oldCancel()
	stream, waitFn, err := dnsConf.Start(newCtx, settings.VerbosityDetailsLevel)
	if err != nil {
		return newCtx, newCancel, nil, err, nil
	}
	go streamMerger.Merge(newCtx, stream, command.MergeName("unbound"), command.MergeColor(constants.ColorUnbound()))
	dnsConf.UseDNSInternally(net.IP{127, 0, 0, 1})                         // use Unbound
	if err := dnsConf.UseDNSSystemWide(net.IP{127, 0, 0, 1}); err != nil { // use Unbound
		return newCtx, newCancel, nil, err, nil
	}
	if err := dnsConf.WaitForUnbound(); err != nil {
		return newCtx, newCancel, nil, err, nil
	}
	// Unbound is up and running at this point
	httpServer.SetUnboundRestart(newCancel)
	waitError := make(chan error)
	waiterError := make(chan error)
	waiter.Add(func() error { //nolint:scopelint
		return <-waiterError
	})
	go func() {
		err := fmt.Errorf("unbound: %w", waitFn())
		waitError <- err
		waiterError <- err
	}()
	if timer == nil {
		waitErr := <-waitError
		return newCtx, newCancel, nil, nil, waitErr
	}
	select {
	case <-timer.C:
		return newCtx, newCancel, nil, nil, nil
	case waitErr := <-waitError:
		return newCtx, newCancel, nil, nil, waitErr
	}
}

func unboundRunLoop(ctx context.Context, logger logging.Logger, dnsConf dns.Configurator,
	settings settings.DNS, uid, gid int,
	waiter command.Waiter, streamMerger command.StreamMerger, httpServer server.Server,
) {
	logger = logger.WithPrefix("unbound dns over tls setup: ")
	if err := fallbackToUnencryptedIPv4DNS(dnsConf, settings.Providers); err != nil {
		logger.Error(err)
	}
	var timer *time.Timer
	if settings.UpdatePeriod > 0 {
		timer = time.NewTimer(settings.UpdatePeriod)
	}
	unboundCtx, unboundCancel := context.WithCancel(ctx)
	defer unboundCancel()
	for ctx.Err() == nil {
		var setupErr, startErr, waitErr error
		unboundCtx, unboundCancel, setupErr, startErr, waitErr = unboundRun(
			ctx, unboundCtx, unboundCancel, timer, dnsConf, settings,
			uid, gid, streamMerger, waiter, httpServer)
		switch {
		case ctx.Err() != nil:
			logger.Warn("context canceled: exiting unbound run loop")
		case timer != nil && !timer.Stop():
			logger.Info("planned restart of unbound")
		case unboundCtx.Err() == context.Canceled:
			logger.Info("triggered restart of unbound")
		case setupErr != nil:
			logger.Warn(setupErr)
		case startErr != nil:
			logger.Error(startErr)
			unboundCancel()
			if err := fallbackToUnencryptedIPv4DNS(dnsConf, settings.Providers); err != nil {
				logger.Error(err)
			}
		case waitErr != nil:
			logger.Warn(waitErr)
			if err := fallbackToUnencryptedIPv4DNS(dnsConf, settings.Providers); err != nil {
				logger.Error(err)
			}
			logger.Warn("restarting unbound because of unexpected exit")
		}
	}
}

func setupPortForwarding(logger logging.Logger, providerConf provider.Provider, filepath models.Filepath, uid, gid int) {
	pfLogger := logger.WithPrefix("port forwarding: ")
	var port uint16
	var err error
	for {
		port, err = providerConf.GetPortForward()
		if err != nil {
			pfLogger.Error(err)
			pfLogger.Info("retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		} else {
			pfLogger.Info("port forwarded is %d", port)
			break
		}
	}
	pfLogger.Info("writing forwarded port to %s", filepath)
	if err := providerConf.WritePortForward(filepath, port, uid, gid); err != nil {
		pfLogger.Error(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := providerConf.AllowPortForwardFirewall(ctx, constants.TUN, port); err != nil {
		pfLogger.Error(err)
	}
}
