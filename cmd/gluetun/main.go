package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
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
	"github.com/qdm12/private-internet-access-docker/internal/dns"
	"github.com/qdm12/private-internet-access-docker/internal/firewall"
	gluetunLogging "github.com/qdm12/private-internet-access-docker/internal/logging"
	"github.com/qdm12/private-internet-access-docker/internal/openvpn"
	"github.com/qdm12/private-internet-access-docker/internal/params"
	"github.com/qdm12/private-internet-access-docker/internal/publicip"
	"github.com/qdm12/private-internet-access-docker/internal/routing"
	"github.com/qdm12/private-internet-access-docker/internal/server"
	"github.com/qdm12/private-internet-access-docker/internal/settings"
	"github.com/qdm12/private-internet-access-docker/internal/shadowsocks"
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
		case "openvpnconfig":
			err = cli.OpenvpnConfig()
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

	client := network.NewClient(15 * time.Second)
	// Create configurators
	fileManager := files.NewFileManager()
	alpineConf := alpine.NewConfigurator(fileManager)
	ovpnConf := openvpn.NewConfigurator(logger, fileManager)
	dnsConf := dns.NewConfigurator(logger, client, fileManager)
	routingConf := routing.NewRouting(logger, fileManager)
	firewallConf := firewall.NewConfigurator(logger, routingConf, fileManager)
	tinyProxyConf := tinyproxy.NewConfigurator(fileManager, logger)
	shadowsocksConf := shadowsocks.NewConfigurator(fileManager, logger)
	streamMerger := command.NewStreamMerger()

	paramsReader := params.NewReader(logger, fileManager)
	fmt.Println(gluetunLogging.Splash(
		paramsReader.GetVersion(),
		paramsReader.GetVcsRef(),
		paramsReader.GetBuildDate()))

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

	err = alpineConf.CreateUser("nonrootuser", uid)
	fatalOnError(err)
	err = fileManager.SetOwnership("/etc/unbound", uid, gid)
	fatalOnError(err)
	err = fileManager.SetOwnership("/etc/tinyproxy", uid, gid)
	fatalOnError(err)

	if allSettings.Firewall.Debug {
		firewallConf.SetDebug()
		routingConf.SetDebug()
	}

	if err := ovpnConf.CheckTUN(); err != nil {
		logger.Warn(err)
		err = ovpnConf.CreateTUN()
		fatalOnError(err)
	}

	connectedCh := make(chan struct{})
	signalConnected := func() {
		connectedCh <- struct{}{}
	}
	defer close(connectedCh)
	go collectStreamLines(ctx, streamMerger, logger, signalConnected)

	// TODO replace these with methods on loopers and pass loopers around
	restartOpenvpn := make(chan struct{})
	portForward := make(chan struct{})
	restartUnbound := make(chan struct{})
	restartPublicIP := make(chan struct{})
	restartTinyproxy := make(chan struct{})
	restartShadowsocks := make(chan struct{})

	if allSettings.Firewall.Enabled {
		err := firewallConf.SetEnabled(ctx, true) // disabled by default
		fatalOnError(err)
	}

	err = firewallConf.SetAllowedSubnets(ctx, allSettings.Firewall.AllowedSubnets)
	fatalOnError(err)

	openvpnLooper := openvpn.NewLooper(allSettings.VPNSP, allSettings.OpenVPN, uid, gid,
		ovpnConf, firewallConf, logger, client, fileManager, streamMerger, fatalOnError)
	// wait for restartOpenvpn
	go openvpnLooper.Run(ctx, restartOpenvpn, portForward, wg)

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
		var restartTickerContext context.Context
		var restartTickerCancel context.CancelFunc = func() {}
		for {
			select {
			case <-ctx.Done():
				restartTickerCancel()
				return
			case <-connectedCh: // blocks until openvpn is connected
				restartTickerCancel()
				restartTickerContext, restartTickerCancel = context.WithCancel(ctx)
				go unboundLooper.RunRestartTicker(restartTickerContext, restartUnbound)
				onConnected(allSettings, logger, routingConf, portForward, restartUnbound, restartPublicIP)
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
	waiting, waited := context.WithTimeout(context.Background(), time.Second)
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

func makeFatalOnError(logger logging.Logger, cancel context.CancelFunc, wg *sync.WaitGroup) func(err error) {
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
		line, level := gluetunLogging.PostProcessLine(line)
		if line == "" {
			return
		}
		switch level {
		case logging.InfoLevel:
			logger.Info(line)
		case logging.WarnLevel:
			logger.Warn(line)
		case logging.ErrorLevel:
			logger.Error(line)
		}
		if strings.Contains(line, "Initialization Sequence Completed") {
			signalConnected()
		}
	}, func(err error) {
		logger.Warn(err)
	})
}

func onConnected(allSettings settings.Settings, logger logging.Logger, routingConf routing.Routing,
	portForward, restartUnbound, restartPublicIP chan<- struct{},
) {
	restartUnbound <- struct{}{}
	restartPublicIP <- struct{}{}
	if allSettings.OpenVPN.Provider.PortForwarding.Enabled {
		time.AfterFunc(5*time.Second, func() {
			portForward <- struct{}{}
		})
	}
	defaultInterface, _, err := routingConf.DefaultRoute()
	if err != nil {
		logger.Warn(err)
	} else {
		vpnGatewayIP, err := routingConf.VPNGatewayIP(defaultInterface)
		if err != nil {
			logger.Warn(err)
		} else {
			logger.Info("Gateway VPN IP address: %s", vpnGatewayIP)
		}
	}
}
