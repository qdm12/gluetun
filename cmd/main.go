package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	libhealthcheck "github.com/qdm12/golibs/healthcheck"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/golibs/signals"
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
	"github.com/qdm12/private-internet-access-docker/internal/settings"
	"github.com/qdm12/private-internet-access-docker/internal/shadowsocks"
	"github.com/qdm12/private-internet-access-docker/internal/splash"
	"github.com/qdm12/private-internet-access-docker/internal/tinyproxy"
	"github.com/qdm12/private-internet-access-docker/internal/windscribe"
)

func main() {
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
	e := env.New(logger)
	client := network.NewClient(15 * time.Second)
	// Create configurators
	fileManager := files.NewFileManager()
	alpineConf := alpine.NewConfigurator(fileManager)
	ovpnConf := openvpn.NewConfigurator(logger, fileManager)
	dnsConf := dns.NewConfigurator(logger, client, fileManager)
	firewallConf := firewall.NewConfigurator(logger)
	routingConf := routing.NewRouting(logger, fileManager)
	piaConf := pia.NewConfigurator(client, fileManager, firewallConf, logger)
	mullvadConf := mullvad.NewConfigurator(fileManager, logger)
	windscribeConf := windscribe.NewConfigurator(fileManager)
	tinyProxyConf := tinyproxy.NewConfigurator(fileManager, logger)
	shadowsocksConf := shadowsocks.NewConfigurator(fileManager, logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	streamMerger := command.NewStreamMerger(ctx)

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

	go func() {
		// Blocking line merging paramsReader for all programs: openvpn, tinyproxy, unbound and shadowsocks
		logger.Info("Launching standard output merger")
		err = streamMerger.CollectLines(func(line string) {
			logger.Info(line)
			if strings.Contains(line, "Initialization Sequence Completed") {
				onConnected(ctx, logger, routingConf, fileManager, piaConf,
					defaultInterface,
					allSettings.VPNSP,
					allSettings.PIA.PortForwarding.Enabled,
					allSettings.PIA.PortForwarding.Filepath,
					allSettings.System.IPStatusFilepath,
					allSettings.System.UID,
					allSettings.System.GID)
			}
		})
		e.FatalOnError(err)
	}()

	waiter := command.NewWaiter()
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
		waiter.Add(waitFn)
		go streamMerger.Merge(stream, command.MergeName("unbound"))
		dnsConf.UseDNSInternally(net.IP{127, 0, 0, 1})       // use Unbound
		err = dnsConf.UseDNSSystemWide(net.IP{127, 0, 0, 1}) // use Unbound
		e.FatalOnError(err)
		err = dnsConf.WaitForUnbound()
		e.FatalOnError(err)
	}

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
		waiter.Add(waitFn)
		go streamMerger.Merge(stream, command.MergeName("tinyproxy"))
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
		waiter.Add(waitFn)
		go streamMerger.Merge(stdout, command.MergeName("shadowsocks"))
		go streamMerger.Merge(stderr, command.MergeName("shadowsocks error"))
	}

	stream, waitFn, err := ovpnConf.Start(ctx)
	e.FatalOnError(err)
	waiter.Add(waitFn)
	go streamMerger.Merge(stream, command.MergeName("openvpn"))
	signals.WaitForExit(func(signal string) int {
		logger.Warn("Caught OS signal %s, shutting down", signal)
		if allSettings.VPNSP == "pia" && allSettings.PIA.PortForwarding.Enabled {
			if err := piaConf.ClearPortForward(allSettings.PIA.PortForwarding.Filepath, allSettings.System.UID, allSettings.System.GID); err != nil {
				logger.Error(err)
			}
		}
		logger.Info("Waiting for processes to exit...")
		errors := waiter.WaitForAll()
		for _, err := range errors {
			logger.Error(err)
		}
		return 0
	})
}

func onConnected(
	ctx context.Context,
	logger logging.Logger,
	routingConf routing.Routing,
	fileManager files.FileManager,
	piaConf pia.Configurator,
	defaultInterface string,
	vpnsp models.VPNProvider,
	portForwarding bool,
	portForwardingFilepath models.Filepath,
	ipStatusFilepath models.Filepath,
	uid, gid int,
) {
	ip, err := routingConf.CurrentPublicIP(defaultInterface)
	if err != nil {
		logger.Error(err)
	} else {
		logger.Info("Tunnel IP is %s, see more information at https://ipinfo.io/%s", ip, ip)
		err := fileManager.WriteLinesToFile(
			string(ipStatusFilepath),
			[]string{ip.String()},
			files.Ownership(uid, gid),
			files.Permissions(0400))
		if err != nil {
			logger.Error(err)
		}
	}
	if vpnsp != constants.PrivateInternetAccess || !portForwarding {
		return
	}
	port, err := piaConf.GetPortForward()
	if err != nil {
		logger.Error("port forwarding:", err)
		return
	}
	logger.Info("port forwarding: Port %d", port)
	if err := piaConf.WritePortForward(portForwardingFilepath, port, uid, gid); err != nil {
		logger.Error("port forwarding:", err)
		return
	}
	if err := piaConf.AllowPortForwardFirewall(ctx, constants.TUN, port); err != nil {
		logger.Error("port forwarding:", err)
		return
	}
}
