package cli

import (
	"fmt"
)

func (c *CLI) Help() {
	//nolint:lll
	fmt.Printf("Usage: gluetun [COMMAND [OPTIONS]]\n\nLightweight swiss-army-knife-like VPN client to multiple VPN service providers.\n\n")
	fmt.Println("Commands:")
	fmt.Printf("  version        \tPrint the version of gluetun\n")
	fmt.Printf("  update         \tUpdate the VPN providers and servers\n")
	fmt.Printf("  healthcheck    \tCheck the health of the VPN connection\n")
	fmt.Printf("  openvpnconfig  \tPrint the OpenVPN configuration (for debugging)\n")
	fmt.Printf("  format-servers \tFormat the servers into a format that can be used by OpenVPN\n")
	fmt.Printf("  genkey         \tGenerate a new OpenVPN key\n")
}
