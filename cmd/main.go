package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	libhealthcheck "github.com/qdm12/golibs/healthcheck"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/golibs/signals"
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
	"github.com/qdm12/private-internet-access-docker/internal/settings"
	"github.com/qdm12/private-internet-access-docker/internal/shadowsocks"
	"github.com/qdm12/private-internet-access-docker/internal/splash"
	"github.com/qdm12/private-internet-access-docker/internal/tinyproxy"
)

const (
	uid, gid = 1000, 1000
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
	paramsReader := params.NewParamsReader(logger)
	fmt.Println(splash.Splash(paramsReader))
	e := env.New(logger)
	client := network.NewClient(15 * time.Second)
	// Create configurators
	fileManager := files.NewFileManager()
	ovpnConf := openvpn.NewConfigurator(logger, fileManager)
	dnsConf := dns.NewConfigurator(logger, client, fileManager)
	firewallConf := firewall.NewConfigurator(logger, fileManager)
	piaConf := pia.NewConfigurator(client, fileManager, firewallConf, logger)
	mullvadConf := mullvad.NewConfigurator(fileManager, logger)
	tinyProxyConf := tinyproxy.NewConfigurator(fileManager, logger)
	shadowsocksConf := shadowsocks.NewConfigurator(fileManager, logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	streamMerger := command.NewStreamMerger(ctx)

	e.PrintVersion("OpenVPN", ovpnConf.Version)
	e.PrintVersion("Unbound", dnsConf.Version)
	e.PrintVersion("IPtables", firewallConf.Version)
	e.PrintVersion("TinyProxy", tinyProxyConf.Version)
	e.PrintVersion("ShadowSocks", shadowsocksConf.Version)

	allSettings, err := settings.GetAllSettings(paramsReader)
	e.FatalOnError(err)
	logger.Info(allSettings.String())

	if err := ovpnConf.CheckTUN(); err != nil {
		logger.Warn(err)
		err = ovpnConf.CreateTUN()
		e.FatalOnError(err)
	}

	var openVPNUser, openVPNPassword string
	switch allSettings.VPNSP {
	case "pia":
		openVPNUser = allSettings.PIA.User
		openVPNPassword = allSettings.PIA.Password
	case "mullvad":
		openVPNUser = allSettings.Mullvad.User
		openVPNPassword = "m"
	}
	err = ovpnConf.WriteAuthFile(openVPNUser, openVPNPassword, uid, gid)
	e.FatalOnError(err)

	// Temporarily reset chain policies allowing Kubernetes sidecar to
	// successfully restart the container. Without this, the existing rules will
	// pre-exist, preventing the nslookup of the PIA region address. These will
	// simply be redundant at Docker runtime as they will already be set this way
	// Thanks to @npawelek https://github.com/npawelek
	err = firewallConf.AcceptAll()
	e.FatalOnError(err)

	go func() {
		// Blocking line merging reader for all programs: openvpn, tinyproxy, unbound and shadowsocks
		logger.Info("Launching standard output merger")
		err = streamMerger.CollectLines(func(line string) { logger.Info(line) })
		e.FatalOnError(err)
	}()

	if allSettings.DNS.Enabled {
		initialDNSToUse := constants.DNSProviderMapping()[allSettings.DNS.Providers[0]]
		dnsConf.UseDNSInternally(initialDNSToUse.IPs[0])
		err = dnsConf.DownloadRootHints(uid, gid)
		e.FatalOnError(err)
		err = dnsConf.DownloadRootKey(uid, gid)
		e.FatalOnError(err)
		err = dnsConf.MakeUnboundConf(allSettings.DNS, uid, gid)
		e.FatalOnError(err)
		stream, waitFn, err := dnsConf.Start(allSettings.DNS.VerbosityDetailsLevel)
		e.FatalOnError(err)
		go func() {
			e.FatalOnError(waitFn())
		}()
		go streamMerger.Merge("unbound", stream)
		dnsConf.UseDNSInternally(net.IP{127, 0, 0, 1})       // use Unbound
		err = dnsConf.UseDNSSystemWide(net.IP{127, 0, 0, 1}) // use Unbound
		e.FatalOnError(err)
		err = dnsConf.WaitForUnbound()
		e.FatalOnError(err)
	}

	var connections []models.OpenVPNConnection
	switch allSettings.VPNSP {
	case "pia":
		connections, err = piaConf.GetOpenVPNConnections(allSettings.PIA.Region, allSettings.OpenVPN.NetworkProtocol, allSettings.PIA.Encryption, allSettings.OpenVPN.TargetIP)
		e.FatalOnError(err)
		err = piaConf.BuildConf(connections, allSettings.PIA.Encryption, allSettings.OpenVPN.Verbosity, uid, gid, allSettings.OpenVPN.Root)
		e.FatalOnError(err)
	case "mullvad":
		connections, err = mullvadConf.GetOpenVPNConnections(allSettings.Mullvad.Country, allSettings.Mullvad.City, allSettings.Mullvad.ISP, allSettings.OpenVPN.NetworkProtocol, allSettings.Mullvad.Port, allSettings.OpenVPN.TargetIP)
		e.FatalOnError(err)
		err = mullvadConf.BuildConf(connections, allSettings.OpenVPN.Verbosity, uid, gid, allSettings.OpenVPN.Root)
		e.FatalOnError(err)
	}

	defaultInterface, defaultGateway, defaultSubnet, err := firewallConf.GetDefaultRoute()
	e.FatalOnError(err)
	err = firewallConf.AddRoutesVia(allSettings.Firewall.AllowedSubnets, defaultGateway, defaultInterface)
	e.FatalOnError(err)
	err = firewallConf.Clear()
	e.FatalOnError(err)
	err = firewallConf.BlockAll()
	e.FatalOnError(err)
	err = firewallConf.CreateGeneralRules()
	e.FatalOnError(err)
	err = firewallConf.CreateVPNRules(constants.TUN, defaultInterface, connections)
	e.FatalOnError(err)
	err = firewallConf.CreateLocalSubnetsRules(defaultSubnet, allSettings.Firewall.AllowedSubnets, defaultInterface)
	e.FatalOnError(err)

	if allSettings.TinyProxy.Enabled {
		err = tinyProxyConf.MakeConf(allSettings.TinyProxy.LogLevel, allSettings.TinyProxy.Port, allSettings.TinyProxy.User, allSettings.TinyProxy.Password, uid, gid)
		e.FatalOnError(err)
		err = firewallConf.AllowAnyIncomingOnPort(allSettings.TinyProxy.Port)
		e.FatalOnError(err)
		stream, waitFn, err := tinyProxyConf.Start()
		e.FatalOnError(err)
		go func() {
			if err := waitFn(); err != nil {
				logger.Error(err)
			}
		}()
		go streamMerger.Merge("tinyproxy", stream)
	}

	if allSettings.ShadowSocks.Enabled {
		err = shadowsocksConf.MakeConf(allSettings.ShadowSocks.Port, allSettings.ShadowSocks.Password, uid, gid)
		e.FatalOnError(err)
		err = firewallConf.AllowAnyIncomingOnPort(allSettings.ShadowSocks.Port)
		e.FatalOnError(err)
		stream, waitFn, err := shadowsocksConf.Start("0.0.0.0", allSettings.ShadowSocks.Port, allSettings.ShadowSocks.Password, allSettings.ShadowSocks.Log)
		e.FatalOnError(err)
		go func() {
			if err := waitFn(); err != nil {
				logger.Error(err)
			}
		}()
		go streamMerger.Merge("shadowsocks", stream)
	}

	if allSettings.VPNSP == "pia" && allSettings.PIA.PortForwarding.Enabled {
		time.AfterFunc(10*time.Second, func() {
			port, err := piaConf.GetPortForward()
			if err != nil {
				logger.Error("port forwarding:", err)
			}
			if err := piaConf.WritePortForward(allSettings.PIA.PortForwarding.Filepath, port); err != nil {
				logger.Error("port forwarding:", err)
			}
			if err := piaConf.AllowPortForwardFirewall(constants.TUN, port); err != nil {
				logger.Error("port forwarding:", err)
			}
		})
	}

	stream, waitFn, err := ovpnConf.Start()
	e.FatalOnError(err)
	go streamMerger.Merge("openvpn", stream)
	go signals.WaitForExit(func(signal string) int {
		logger.Warn("Caught OS signal %s, shutting down", signal)
		time.Sleep(100 * time.Millisecond) // wait for other processes to exit
		return 0
	})
	e.FatalOnError(waitFn())
}
