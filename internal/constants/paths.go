package constants

import (
	"github.com/qdm12/gluetun/internal/models"
)

const (
	// UnboundConf is the file path to the Unbound configuration file
	UnboundConf models.Filepath = "/etc/unbound/unbound.conf"
	// ResolvConf is the file path to the system resolv.conf file
	ResolvConf models.Filepath = "/etc/resolv.conf"
	// CACertificates is the file path to the CA certificates file
	CACertificates models.Filepath = "/etc/ssl/certs/ca-certificates.crt"
	// OpenVPNAuthConf is the file path to the OpenVPN auth file
	OpenVPNAuthConf models.Filepath = "/etc/openvpn/auth.conf"
	// OpenVPNConf is the file path to the OpenVPN client configuration file
	OpenVPNConf models.Filepath = "/etc/openvpn/target.ovpn"
	// TunnelDevice is the file path to tun device
	TunnelDevice models.Filepath = "/dev/net/tun"
	// NetRoute is the path to the file containing information on the network route
	NetRoute models.Filepath = "/proc/net/route"
	// TinyProxyConf is the filepath to the tinyproxy configuration file
	TinyProxyConf models.Filepath = "/etc/tinyproxy/tinyproxy.conf"
	// ShadowsocksConf is the filepath to the shadowsocks configuration file
	ShadowsocksConf models.Filepath = "/etc/shadowsocks.json"
	// RootHints is the filepath to the root.hints file used by Unbound
	RootHints models.Filepath = "/etc/unbound/root.hints"
	// RootKey is the filepath to the root.key file used by Unbound
	RootKey models.Filepath = "/etc/unbound/root.key"
)
