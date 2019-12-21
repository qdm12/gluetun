package main

import (
	"fmt"
	"os"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/private-internet-access-docker/internal/command"
	"github.com/qdm12/private-internet-access-docker/internal/settings"
)

func main() {
	// TODO use colors, emojis, maybe move to Golibs
	fmt.Printf(`
	=========================================
	=========================================
	============= PIA CONTAINER =============
	=========================================
	=========================================
	== by github.com/qdm12 - Quentin McGaw ==
	`)
	printVersion("OpenVPN", command.VersionOpenVPN)
	printVersion("Unbound", command.VersionUnbound)
	printVersion("IPtables", command.VersionIptables)
	printVersion("TinyProxy", command.VersionTinyProxy)
	printVersion("ShadowSocks", command.VersionShadowSocks)
	openVPNSettings, err := settings.GetOpenVPNSettings()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(openVPNSettings)
	PIASettings, err := settings.GetPIASettings()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(PIASettings)
	DNSSettings, err := settings.GetDNSSettings()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(DNSSettings)
	firewallSettings, err := settings.GetFirewallSettings()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(firewallSettings)
	tinyProxySettings, err := settings.GetTinyProxySettings()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(tinyProxySettings)
	shadowSocksSettings, err := settings.GetShadowSocksSettings()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(shadowSocksSettings)
}

func printVersion(program string, commandFn func() (string, error)) {
	version, err := commandFn()
	if err != nil {
		logging.Err(err)
	} else {
		fmt.Printf("%s version: %s\n", program, version)
	}
}
