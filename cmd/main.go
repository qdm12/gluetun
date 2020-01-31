package main

import (
	"fmt"
	"time"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/command"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/dns"
	"github.com/qdm12/private-internet-access-docker/internal/env"
	"github.com/qdm12/private-internet-access-docker/internal/firewall"
	"github.com/qdm12/private-internet-access-docker/internal/openvpn"
	"github.com/qdm12/private-internet-access-docker/internal/params"
	"github.com/qdm12/private-internet-access-docker/internal/pia"
	"github.com/qdm12/private-internet-access-docker/internal/settings"
)

func main() {
	// TODO use colors, emojis, maybe move to Golibs
	logger, err := logging.NewLogger(logging.ConsoleEncoding, logging.InfoLevel, -1)
	if err != nil {
		panic(err)
	}
	e := env.New(logger)
	fmt.Printf(`
	=========================================
	=========================================
	============= PIA CONTAINER =============
	=========================================
	=========================================
	== by github.com/qdm12 - Quentin McGaw ==
	`)
	cmd := command.NewCommand()
	e.PrintVersion("OpenVPN", cmd.VersionOpenVPN)
	e.PrintVersion("Unbound", cmd.VersionUnbound)
	e.PrintVersion("IPtables", cmd.VersionIptables)
	e.PrintVersion("TinyProxy", cmd.VersionTinyProxy)
	e.PrintVersion("ShadowSocks", cmd.VersionShadowSocks)
	paramsReader := params.NewParamsReader(logger)
	allSettings, err := settings.GetAllSettings(paramsReader)
	e.FatalOnError(err)
	logger.Info(allSettings)
	fileManager := files.NewFileManager()
	ovpnConf := openvpn.NewConfigurator(logger, fileManager)
	logger.Info("Writing auth file")
	err = ovpnConf.WriteAuthFile(allSettings.PIA.User, allSettings.PIA.Password)
	e.FatalOnError(err)
	logger.Info("Checking /dev/tun device")
	err = ovpnConf.CheckTUN()
	e.FatalOnError(err)
	client := network.NewClient(3 * time.Second)
	if allSettings.DNS.Enabled {
		logger.Info("Setting up DNS over TLS")
		dnsConf := dns.NewConfigurator(logger, client, fileManager)
		err = dnsConf.MakeUnboundConf(allSettings.DNS)
		e.FatalOnError(err)
		err = cmd.Unbound()
		e.FatalOnError(err)
		err = dnsConf.SetLocalNameserver()
		e.FatalOnError(err)
	}
	piaConf := pia.NewConfigurator(client)
	logger.Info("Configuring PIA")
	lines, err := piaConf.DownloadOvpnConfig(allSettings.PIA.Encryption, allSettings.OpenVPN.NetworkProtocol, allSettings.PIA.Region)
	e.FatalOnError(err)
	VPNIPs, port, VPNDevice, err := piaConf.ParseConfig(lines)
	e.FatalOnError(err)
	lines, err = piaConf.ModifyLines(lines, VPNIPs, port)
	e.FatalOnError(err)
	fileManager.WriteLinesToFile(constants.OpenVPNConf, lines)
	e.FatalOnError(err)
	firewallConf := firewall.NewConfigurator(logger, fileManager)
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
}
