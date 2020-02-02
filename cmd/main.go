package main

import (
	"context"
	"fmt"
	"io"
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
	"github.com/qdm12/private-internet-access-docker/internal/tinyproxy"
)

func main() {
	logger, err := logging.NewLogger(logging.ConsoleEncoding, logging.InfoLevel, -1)
	if err != nil {
		panic(err)
	}
	e := env.New(logger)
	fmt.Printf(`=========================================
=========================================
============= PIA CONTAINER =============
=========================================
=========================================
========== by github.com/qdm12 ==========
`)
	client := network.NewClient(3 * time.Second)
	// Create configurators
	fileManager := files.NewFileManager()
	ovpnConf := openvpn.NewConfigurator(logger, fileManager)
	dnsConf := dns.NewConfigurator(logger, client, fileManager)
	piaConf := pia.NewConfigurator(client)
	firewallConf := firewall.NewConfigurator(logger, fileManager)
	tinyProxyConf := tinyproxy.NewConfigurator(fileManager)
	shadowsocksConf := shadowsocks.NewConfigurator(fileManager)

	e.PrintVersion("OpenVPN", ovpnConf.Version)
	e.PrintVersion("Unbound", dnsConf.Version)
	e.PrintVersion("IPtables", firewallConf.Version)
	e.PrintVersion("TinyProxy", tinyProxyConf.Version)
	e.PrintVersion("ShadowSocks", shadowsocksConf.Version)

	allSettings, err := settings.GetAllSettings(params.NewParamsReader(logger))
	e.FatalOnError(err)
	logger.Info(allSettings.String())

	logger.Info("Checking /dev/tun device")
	err = ovpnConf.CheckTUN()
	e.FatalOnError(err)

	logger.Info("Writing auth file")
	err = ovpnConf.WriteAuthFile(allSettings.PIA.User, allSettings.PIA.Password)
	e.FatalOnError(err)

	stdouts := make(map[string]io.ReadCloser)

	if allSettings.DNS.Enabled {
		logger.Info("Setting up DNS over TLS")
		err = dnsConf.DownloadRootHints()
		e.FatalOnError(err)
		err = dnsConf.DownloadRootKey()
		e.FatalOnError(err)
		err = dnsConf.MakeUnboundConf(allSettings.DNS)
		e.FatalOnError(err)
		stdouts["Unbound"], err = dnsConf.Start()
		e.FatalOnError(err)
		err = dnsConf.SetLocalNameserver()
		e.FatalOnError(err)
	}

	logger.Info("Configuring PIA")
	lines, err := piaConf.DownloadOvpnConfig(allSettings.PIA.Encryption, allSettings.OpenVPN.NetworkProtocol, allSettings.PIA.Region)
	e.FatalOnError(err)
	VPNIPs, port, VPNDevice, err := piaConf.ParseConfig(lines)
	e.FatalOnError(err)
	lines, err = piaConf.ModifyLines(lines, VPNIPs, port)
	e.FatalOnError(err)
	fileManager.WriteLinesToFile(string(constants.OpenVPNConf), lines)
	e.FatalOnError(err)

	logger.Info("Configuring firewall")
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
	err = firewallConf.CreateVPNRules(VPNDevice, VPNIPs, defaultInterface, port, allSettings.OpenVPN.NetworkProtocol)
	e.FatalOnError(err)
	err = firewallConf.CreateLocalSubnetsRules(defaultSubnet, allSettings.Firewall.AllowedSubnets, defaultInterface)
	e.FatalOnError(err)

	if allSettings.TinyProxy.Enabled {
		logger.Info("Configuring Tinyproxy")
		err = tinyProxyConf.MakeConf(allSettings.TinyProxy.LogLevel, allSettings.ShadowSocks.Port, allSettings.TinyProxy.User, allSettings.TinyProxy.Password)
		e.FatalOnError(err)
		stdouts["TinyProxy"], err = tinyProxyConf.Start()
		e.FatalOnError(err)
	}

	if allSettings.ShadowSocks.Enabled {
		logger.Info("Configuring Shadowsocks")
		err = shadowsocksConf.MakeConf(allSettings.ShadowSocks.Port, allSettings.TinyProxy.Password)
		e.FatalOnError(err)
		stdouts["Shadowsocks"], err = shadowsocksConf.Start(allSettings.ShadowSocks.Log)
		e.FatalOnError(err)
	}

	if allSettings.PIA.PortForwarding.Enabled {
		time.AfterFunc(10*time.Second, func() {
			piaConf.PortForward(allSettings.PIA.PortForwarding.Filepath)
		})
	}

	stdouts["Shadowsocks"], err = ovpnConf.Start()
	e.FatalOnError(err)

	// Blocking line merging reader for all programs: openvpn, tinyproxy, unbound and shadowsocks
	command.NewCommander().MergeLineReaders(
		context.Background(),
		func(line string) {
			logger.Info(line)
		},
		stdouts,
	)
}
