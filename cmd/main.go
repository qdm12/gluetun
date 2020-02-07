package main

import (
	"context"
	"fmt"
	"time"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/dns"
	"github.com/qdm12/private-internet-access-docker/internal/env"
	"github.com/qdm12/private-internet-access-docker/internal/firewall"
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
	fmt.Println(splash.Splash())
	logger, err := logging.NewLogger(logging.ConsoleEncoding, logging.InfoLevel, -1)
	if err != nil {
		panic(err)
	}
	e := env.New(logger)
	client := network.NewClient(3 * time.Second)
	// Create configurators
	fileManager := files.NewFileManager()
	ovpnConf := openvpn.NewConfigurator(logger, fileManager)
	dnsConf := dns.NewConfigurator(logger, client, fileManager)
	firewallConf := firewall.NewConfigurator(logger, fileManager)
	piaConf := pia.NewConfigurator(client, fileManager, firewallConf, logger)
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

	allSettings, err := settings.GetAllSettings(params.NewParamsReader(logger))
	e.FatalOnError(err)
	logger.Info(allSettings.String())

	err = ovpnConf.CheckTUN()
	e.FatalOnError(err)

	err = ovpnConf.WriteAuthFile(allSettings.PIA.User, allSettings.PIA.Password, uid, gid)
	e.FatalOnError(err)

	// Temporarily reset chain policies allowing Kubernetes sidecar to
	// successfully restart the container. Without this, the existing rules will
	// pre-exist, preventing the nslookup of the PIA region address. These will
	// simply be redundant at Docker runtime as they will already be set this way
	// Thanks to @npawelek https://github.com/npawelek
	err = firewallConf.Clear()
	e.FatalOnError(err)
	err = firewallConf.AcceptAll()
	e.FatalOnError(err)

	if allSettings.DNS.Enabled {
		err = dnsConf.DownloadRootHints(uid, gid)
		e.FatalOnError(err)
		err = dnsConf.DownloadRootKey(uid, gid)
		e.FatalOnError(err)
		err = dnsConf.MakeUnboundConf(allSettings.DNS, uid, gid)
		e.FatalOnError(err)
		stream, err := dnsConf.Start(allSettings.DNS.VerbosityDetailsLevel)
		e.FatalOnError(err)
		go streamMerger.Merge("unbound", stream)
		err = dnsConf.SetLocalNameserver()
		e.FatalOnError(err)
	}

	lines, err := piaConf.DownloadOvpnConfig(allSettings.PIA.Encryption, allSettings.OpenVPN.NetworkProtocol, allSettings.PIA.Region)
	e.FatalOnError(err)
	VPNIPs, port, VPNDevice, err := piaConf.ParseConfig(lines)
	e.FatalOnError(err)
	lines = piaConf.ModifyLines(lines, VPNIPs, port)
	fileManager.WriteLinesToFile(string(constants.OpenVPNConf), lines)
	e.FatalOnError(err)

	defaultInterface, defaultGateway, defaultSubnet, err := firewallConf.GetDefaultRoute()
	e.FatalOnError(err)
	err = firewallConf.AddRoutesVia(allSettings.Firewall.AllowedSubnets, defaultGateway, defaultInterface, VPNDevice)
	e.FatalOnError(err)
	err = firewallConf.Clear()
	e.FatalOnError(err)
	err = firewallConf.BlockAll()
	e.FatalOnError(err)
	err = firewallConf.CreateGeneralRules()
	e.FatalOnError(err)
	err = firewallConf.CreateVPNRules(VPNDevice, VPNIPs, defaultInterface, port, allSettings.OpenVPN.NetworkProtocol)
	e.FatalOnError(err)
	err = firewallConf.CreateLocalSubnetsRules(defaultSubnet, allSettings.Firewall.AllowedSubnets, defaultInterface)
	e.FatalOnError(err)

	if allSettings.TinyProxy.Enabled {
		err = tinyProxyConf.MakeConf(allSettings.TinyProxy.LogLevel, allSettings.ShadowSocks.Port, allSettings.TinyProxy.User, allSettings.TinyProxy.Password, uid, gid)
		e.FatalOnError(err)
		stream, err := tinyProxyConf.Start()
		e.FatalOnError(err)
		go streamMerger.Merge("tinyproxy", stream)
	}

	if allSettings.ShadowSocks.Enabled {
		err = shadowsocksConf.MakeConf(allSettings.ShadowSocks.Port, allSettings.TinyProxy.Password, uid, gid)
		e.FatalOnError(err)
		stream, err := shadowsocksConf.Start("0.0.0.0", allSettings.ShadowSocks.Port, allSettings.ShadowSocks.Password, allSettings.ShadowSocks.Log)
		e.FatalOnError(err)
		go streamMerger.Merge("shadowsocks", stream)
	}

	if allSettings.PIA.PortForwarding.Enabled {
		time.AfterFunc(10*time.Second, func() {
			port, err := piaConf.GetPortForward()
			if err != nil {
				logger.Error("port forwarding:", err)
			}
			if err := piaConf.WritePortForward(allSettings.PIA.PortForwarding.Filepath, port); err != nil {
				logger.Error("port forwarding:", err)
			}
			if err := piaConf.AllowPortForwardFirewall(VPNDevice, port); err != nil {
				logger.Error("port forwarding:", err)
			}
		})
	}

	stream, err := ovpnConf.Start()
	e.FatalOnError(err)
	go streamMerger.Merge("openvpn", stream)

	// Blocking line merging reader for all programs: openvpn, tinyproxy, unbound and shadowsocks
	logger.Info("Launching standard output merger")
	err = streamMerger.CollectLines(func(line string) { logger.Info(line) })
	e.FatalOnError(err)
}
