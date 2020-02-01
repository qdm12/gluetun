package constants

import (
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

const (
	// UnboundConf is the file path to the Unbound configuration file
	UnboundConf models.Filepath = "/etc/unbound/unbound.conf"
	// ResolvConf is the file path to the system resolv.conf file
	ResolvConf models.Filepath = "/etc/resolv.conf"
	// OpenVPNAuthConf is the file path to the OpenVPN auth file
	OpenVPNAuthConf models.Filepath = "/etc/openvpn/auth.conf"
	// OpenVPNConf is the file path to the OpenVPN client configuration file
	OpenVPNConf models.Filepath = "/etc/openvpn/target.ovpn"
	// TunnelDevice is the file path to tun device
	TunnelDevice models.Filepath = "/dev/net/tun"
	// NetRoute is the path to the file containing information on the network route
	NetRoute models.Filepath = "/proc/net/route"
)
