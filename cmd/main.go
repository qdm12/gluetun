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
	e.PrintVersion("OpenVPN", command.VersionOpenVPN)
	e.PrintVersion("Unbound", command.VersionUnbound)
	e.PrintVersion("IPtables", command.VersionIptables)
	e.PrintVersion("TinyProxy", command.VersionTinyProxy)
	e.PrintVersion("ShadowSocks", command.VersionShadowSocks)
	params := params.NewParams(logger)
	allSettings, err := settings.GetAllSettings(params)
	e.FatalOnError(err)
	logger.Info(allSettings)
	fileManager := files.NewFileManager()
	ovpnConf := openvpn.NewConfigurator(logger, fileManager)
	err = ovpnConf.WriteAuthFile(allSettings.PIA.User, allSettings.PIA.Password)
	e.FatalOnError(err)
	err = ovpnConf.CheckTUN()
	e.FatalOnError(err)
	client := network.NewClient(3 * time.Second)
	if allSettings.DNS.Enabled {
		dnsConf := dns.NewConfigurator(logger, client)
		lines, warnings := dnsConf.MakeUnboundConf(allSettings.DNS)
		for _, warning := range warnings {
			logger.Warn(warning)
		}
		err = fileManager.WriteLinesToFile(constants.UnboundConf, lines)
		e.FatalOnError(err)
	}
	piaConf := pia.NewConfigurator()
	ovpnLines, err := piaConf.Get(client, allSettings.PIA.Encryption, allSettings.OpenVPN.NetworkProtocol, allSettings.PIA.Region)
	e.FatalOnError(err)
	IPs, port, device, err := piaConf.Read(ovpnLines)
	e.FatalOnError(err)
	ovpnLines, err = piaConf.Modify(ovpnLines, IPs, port)
	e.FatalOnError(err)
	err = fileManager.WriteLinesToFile(constants.OpenVPNConf, ovpnLines)
	e.FatalOnError(err)
	fmt.Println(device)
}
