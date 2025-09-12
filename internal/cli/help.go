package cli

import (
	"fmt"
)

func (c *CLI) Help() {
	//nolint:lll
	fmt.Printf("Usage: gluetun [COMMAND [OPTIONS]]\n\nLightweight swiss-army-knife-like VPN client to multiple VPN service providers.\n\n")
	fmt.Println("Commands:")
	fmt.Printf("  update         \tUpdate the VPN servers data for some or all providers\n")
	fmt.Printf("  healthcheck    \tCheck the health of the VPN connection of another Gluetun instance\n")
	fmt.Printf("  openvpnconfig  \tPrint the OpenVPN configuration (for debugging)\n")
	fmt.Printf("  format-servers \tFormat the servers data into a Markdown table\n")
	fmt.Printf("  genkey         \tGenerate a new 32 bytes Wireguard key (base58 encoded)\n")
}
