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
	"github.com/qdm12/private-internet-access-docker/internal/env"
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

func main() { //nolint:gocognit
	logger, err := logging.NewLogger(logging.ConsoleEncoding, logging.InfoLevel, -1)
	if err != nil {
		panic(err)
	}
	if libhealthcheck.Mode(os.Args) {
		if err := healthcheck.HealthCheck(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	paramsReader := params.NewReader(logger)
	fmt.Println(splash.Splash(paramsReader))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	e := env.New(logger, cancel)

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

	e.PrintVersion(ctx, "OpenVPN", ovpnConf.Version)
	e.PrintVersion(ctx, "Unbound", dnsConf.Version)
	e.PrintVersion(ctx, "IPtables", firewallConf.Version)
	e.PrintVersion(ctx, "TinyProxy", tinyProxyConf.Version)
	e.PrintVersion(ctx, "ShadowSocks", shadowsocksConf.Version)

	allSettings, err := settings.GetAllSettings(paramsReader)
	e.FatalOnError(err)
	logger.Info(allSettings.String())

	err = alpineConf.CreateUser("nonrootuser", allSettings.System.UID)
	e.FatalOnError(err)
	err = fileManager.SetOwnership("/etc/unbound", allSettings.System.UID, allSettings.System.GID)
	e.FatalOnError(err)
	err = fileManager.SetOwnership("/etc/tinyproxy", allSettings.System.UID, allSettings.System.GID)
	e.FatalOnError(err)

	if err := ovpnConf.CheckTUN(); err != nil {
		logger.Warn(err)
		err = ovpnConf.CreateTUN()
		e.FatalOnError(err)
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
	e.FatalOnError(err)

	defaultInterface, defaultGateway, defaultSubnet, err := routingConf.DefaultRoute()
	e.FatalOnError(err)

	// Temporarily reset chain policies allowing Kubernetes sidecar to
	// successfully restart the container. Without this, the existing rules will
	// pre-exist, preventing the nslookup of the PIA region address. These will
	// simply be redundant at Docker runtime as they will already be set this way
	// Thanks to @npawelek https://github.com/npawelek
	err = firewallConf.AcceptAll(ctx)
	e.FatalOnError(err)

	connected, signalConnected := context.WithCancel(context.Background())
	go func() {
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
	}()

	waiter := command.NewWaiter()

	var connections []models.OpenVPNConnection
	switch allSettings.VPNSP {
	case constants.PrivateInternetAccess:
		connections, err = piaConf.GetOpenVPNConnections(
			allSettings.PIA.Region,
			allSettings.OpenVPN.NetworkProtocol,
			allSettings.PIA.Encryption,
			allSettings.OpenVPN.TargetIP)
		e.FatalOnError(err)
		err = piaConf.BuildConf(
			connections,
			allSettings.PIA.Encryption,
			allSettings.OpenVPN.Verbosity,
			allSettings.System.UID,
			allSettings.System.GID,
			allSettings.OpenVPN.Root,
			allSettings.OpenVPN.Cipher,
			allSettings.OpenVPN.Auth)
		e.FatalOnError(err)
	case constants.Mullvad:
		connections, err = mullvadConf.GetOpenVPNConnections(
			allSettings.Mullvad.Country,
			allSettings.Mullvad.City,
			allSettings.Mullvad.ISP,
			allSettings.OpenVPN.NetworkProtocol,
			allSettings.Mullvad.Port,
			allSettings.OpenVPN.TargetIP)
		e.FatalOnError(err)
		err = mullvadConf.BuildConf(
			connections,
			allSettings.OpenVPN.Verbosity,
			allSettings.System.UID,
			allSettings.System.GID,
			allSettings.OpenVPN.Root,
			allSettings.OpenVPN.Cipher)
		e.FatalOnError(err)
	case constants.Windscribe:
		connections, err = windscribeConf.GetOpenVPNConnections(
			allSettings.Windscribe.Region,
			allSettings.OpenVPN.NetworkProtocol,
			allSettings.Windscribe.Port,
			allSettings.OpenVPN.TargetIP)
		e.FatalOnError(err)
		err = windscribeConf.BuildConf(
			connections,
			allSettings.OpenVPN.Verbosity,
			allSettings.System.UID,
			allSettings.System.GID,
			allSettings.OpenVPN.Root,
			allSettings.OpenVPN.Cipher,
			allSettings.OpenVPN.Auth)
		e.FatalOnError(err)
	}

	err = routingConf.AddRoutesVia(ctx, allSettings.Firewall.AllowedSubnets, defaultGateway, defaultInterface)
	e.FatalOnError(err)
	err = firewallConf.Clear(ctx)
	e.FatalOnError(err)
	err = firewallConf.BlockAll(ctx)
	e.FatalOnError(err)
	err = firewallConf.CreateGeneralRules(ctx)
	e.FatalOnError(err)
	err = firewallConf.CreateVPNRules(ctx, constants.TUN, defaultInterface, connections)
	e.FatalOnError(err)
	err = firewallConf.CreateLocalSubnetsRules(ctx, defaultSubnet, allSettings.Firewall.AllowedSubnets, defaultInterface)
	e.FatalOnError(err)

	if allSettings.TinyProxy.Enabled {
		err = tinyProxyConf.MakeConf(
			allSettings.TinyProxy.LogLevel,
			allSettings.TinyProxy.Port,
			allSettings.TinyProxy.User,
			allSettings.TinyProxy.Password,
			allSettings.System.UID,
			allSettings.System.GID)
		e.FatalOnError(err)
		err = firewallConf.AllowAnyIncomingOnPort(ctx, allSettings.TinyProxy.Port)
		e.FatalOnError(err)
		stream, waitFn, err := tinyProxyConf.Start(ctx)
		e.FatalOnError(err)
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
		e.FatalOnError(err)
		err = firewallConf.AllowAnyIncomingOnPort(ctx, allSettings.ShadowSocks.Port)
		e.FatalOnError(err)
		stdout, stderr, waitFn, err := shadowsocksConf.Start(ctx, "0.0.0.0", allSettings.ShadowSocks.Port, allSettings.ShadowSocks.Password, allSettings.ShadowSocks.Log)
		e.FatalOnError(err)
		waiter.Add(func() error {
			err := waitFn()
			logger.Error("shadowsocks: %s", err)
			return err
		})
		go streamMerger.Merge(ctx, stdout, command.MergeName("shadowsocks"), command.MergeColor(constants.ColorShadowsocks()))
		go streamMerger.Merge(ctx, stderr, command.MergeName("shadowsocks error"), command.MergeColor(constants.ColorShadowsocksError()))
	}

	httpServer := server.New("0.0.0.0:8000", logger)

	// Runs openvpn and restarts it if it does not exit cleanly
	openvpnCancelSet, signalOpenvpnCancelSet := context.WithCancel(context.Background())
	go func() {
		waitErrors := make(chan error)
		for {
			openvpnCtx, openvpnCancel := context.WithCancel(ctx)
			stream, waitFn, err := ovpnConf.Start(openvpnCtx)
			e.FatalOnError(err)
			httpServer.SetOpenVPNRestart(openvpnCancel)
			signalOpenvpnCancelSet()
			go streamMerger.Merge(openvpnCtx, stream, command.MergeName("openvpn"), command.MergeColor(constants.ColorOpenvpn()))
			waiter.Add(func() error {
				err := <-waitErrors
				logger.Error("openvpn: %s", err)
				return err
			})
			if err := waitFn(); err != nil {
				waitErrors <- err
			} else {
				break
			}
			openvpnCancel()
		}
	}()

	<-openvpnCancelSet.Done()

	waiter.Add(func() error {
		err := httpServer.Run(ctx)
		logger.Error("http server: %s", err)
		return err
	})

	go func() {
		<-connected.Done() // blocks until openvpn is connected

		if allSettings.DNS.Enabled {
			initialDNSToUse := constants.DNSProviderMapping()[allSettings.DNS.Providers[0]]
			dnsConf.UseDNSInternally(initialDNSToUse.IPs[0])
			err = dnsConf.DownloadRootHints(allSettings.System.UID, allSettings.System.GID)
			e.FatalOnError(err)
			err = dnsConf.DownloadRootKey(allSettings.System.UID, allSettings.System.GID)
			e.FatalOnError(err)
			err = dnsConf.MakeUnboundConf(allSettings.DNS, allSettings.System.UID, allSettings.System.GID)
			e.FatalOnError(err)
			stream, waitFn, err := dnsConf.Start(ctx, allSettings.DNS.VerbosityDetailsLevel)
			e.FatalOnError(err)
			waiter.Add(func() error {
				err := waitFn()
				logger.Error("unbound: %s", err)
				return err
			})
			go streamMerger.Merge(ctx, stream, command.MergeName("unbound"), command.MergeColor(constants.ColorUnbound()))
			dnsConf.UseDNSInternally(net.IP{127, 0, 0, 1})       // use Unbound
			err = dnsConf.UseDNSSystemWide(net.IP{127, 0, 0, 1}) // use Unbound
			e.FatalOnError(err)
			err = dnsConf.WaitForUnbound()
			e.FatalOnError(err)
			logger.Info("DNS over TLS with Unbound setup completed")
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

		if allSettings.PIA.PortForwarding.Enabled {
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
			pfLogger.Info("writing forwarded port to %s", allSettings.PIA.PortForwarding.Filepath)
			if err := piaConf.WritePortForward(allSettings.PIA.PortForwarding.Filepath, port, allSettings.System.UID, allSettings.System.GID); err != nil {
				pfLogger.Error(err)
			}
			pfLogger.Info("allowing forwarded port %d through firewall", port)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			if err := piaConf.AllowPortForwardFirewall(ctx, constants.TUN, port); err != nil {
				pfLogger.Error(err)
			}
		}
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
