package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
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
	wg := &sync.WaitGroup{}
	fatalOnError := makeFatalOnError(logger, cancel, wg)
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

	// Should never change
	uid, gid := allSettings.System.UID, allSettings.System.GID

	providerConf := provider.New(allSettings.VPNSP, logger, client, fileManager, firewallConf)

	if !allSettings.Firewall.Enabled {
		firewallConf.Disable()
	}

	err = alpineConf.CreateUser("nonrootuser", uid)
	fatalOnError(err)
	err = fileManager.SetOwnership("/etc/unbound", uid, gid)
	fatalOnError(err)
	err = fileManager.SetOwnership("/etc/tinyproxy", uid, gid)
	fatalOnError(err)

	if err := ovpnConf.CheckTUN(); err != nil {
		logger.Warn(err)
		err = ovpnConf.CreateTUN()
		fatalOnError(err)
	}

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

	connections, err := providerConf.GetOpenVPNConnections(allSettings.OpenVPN.Provider.ServerSelection)
	fatalOnError(err)
	err = providerConf.BuildConf(
		connections,
		allSettings.OpenVPN.Verbosity,
		uid,
		gid,
		allSettings.OpenVPN.Root,
		allSettings.OpenVPN.Cipher,
		allSettings.OpenVPN.Auth,
		allSettings.OpenVPN.Provider.ExtraConfigOptions,
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

	restartOpenvpn := make(chan struct{})
	restartUnbound := make(chan struct{})
	restartPublicIP := make(chan struct{})
	restartTinyproxy := make(chan struct{})
	restartShadowsocks := make(chan struct{})

	openvpnLooper := openvpn.NewLooper(ovpnConf, allSettings.OpenVPN, logger, streamMerger, fatalOnError, uid, gid)
	// wait for restartOpenvpn
	go openvpnLooper.Run(ctx, restartOpenvpn, wg)

	unboundLooper := dns.NewLooper(dnsConf, allSettings.DNS, logger, streamMerger, uid, gid)
	// wait for restartUnbound
	go unboundLooper.Run(ctx, restartUnbound, wg)

	publicIPLooper := publicip.NewLooper(client, logger, fileManager, allSettings.System.IPStatusFilepath, uid, gid)
	go publicIPLooper.Run(ctx, restartPublicIP)
	go publicIPLooper.RunRestartTicker(ctx, restartPublicIP)

	tinyproxyLooper := tinyproxy.NewLooper(tinyProxyConf, firewallConf, allSettings.TinyProxy, logger, streamMerger, uid, gid)
	go tinyproxyLooper.Run(ctx, restartTinyproxy, wg)

	shadowsocksLooper := shadowsocks.NewLooper(shadowsocksConf, firewallConf, allSettings.ShadowSocks, allSettings.DNS, logger, streamMerger, uid, gid)
	go shadowsocksLooper.Run(ctx, restartShadowsocks, wg)

	if allSettings.TinyProxy.Enabled {
		restartTinyproxy <- struct{}{}
	}
	if allSettings.ShadowSocks.Enabled {
		restartShadowsocks <- struct{}{}
	}

	go func() {
		first := true
		var restartTickerContext context.Context
		var restartTickerCancel context.CancelFunc = func() {}
		for {
			select {
			case <-ctx.Done():
				restartTickerCancel()
				return
			case <-connectedCh: // blocks until openvpn is connected
				if first {
					first = false
					restartUnbound <- struct{}{}
				}
				restartTickerCancel()
				restartTickerContext, restartTickerCancel = context.WithCancel(ctx)
				go unboundLooper.RunRestartTicker(restartTickerContext, restartUnbound)
				onConnected(allSettings, logger, routingConf, defaultInterface, providerConf, restartPublicIP)
			}
		}
	}()

	httpServer := server.New("0.0.0.0:8000", logger, restartOpenvpn, restartUnbound)
	go httpServer.Run(ctx, wg)

	// Start openvpn for the first time
	restartOpenvpn <- struct{}{}

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
	if allSettings.OpenVPN.Provider.PortForwarding.Enabled {
		logger.Info("Clearing forwarded port status file %s", allSettings.OpenVPN.Provider.PortForwarding.Filepath)
		if err := fileManager.Remove(string(allSettings.OpenVPN.Provider.PortForwarding.Filepath)); err != nil {
			logger.Error(err)
			exitStatus = 1
		}
	}
	wg.Wait()
	return exitStatus
}

func makeFatalOnError(logger logging.Logger, cancel func(), wg *sync.WaitGroup) func(err error) {
	return func(err error) {
		if err != nil {
			logger.Error(err)
			cancel()
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			go func() {
				wg.Wait()
				cancel()
			}()
			<-ctx.Done() // either timeout or wait group completed
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

func onConnected(allSettings settings.Settings,
	logger logging.Logger, routingConf routing.Routing, defaultInterface string,
	providerConf provider.Provider, restartPublicIP chan<- struct{},
) {
	restartPublicIP <- struct{}{}
	uid, gid := allSettings.System.UID, allSettings.System.GID
	if allSettings.OpenVPN.Provider.PortForwarding.Enabled {
		time.AfterFunc(5*time.Second, func() {
			setupPortForwarding(logger, providerConf, allSettings.OpenVPN.Provider.PortForwarding.Filepath, uid, gid)
		})
	}

	vpnGatewayIP, err := routingConf.VPNGatewayIP(defaultInterface)
	if err != nil {
		logger.Warn(err)
	} else {
		logger.Info("Gateway VPN IP address: %s", vpnGatewayIP)
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
