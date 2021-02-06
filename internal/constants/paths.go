package constants

const (
	// UnboundConf is the file path to the Unbound configuration file.
	UnboundConf string = "/etc/unbound/unbound.conf"
	// ResolvConf is the file path to the system resolv.conf file.
	ResolvConf string = "/etc/resolv.conf"
	// CACertificates is the file path to the CA certificates file.
	CACertificates string = "/etc/ssl/certs/ca-certificates.crt"
	// OpenVPNAuthConf is the file path to the OpenVPN auth file.
	OpenVPNAuthConf string = "/etc/openvpn/auth.conf"
	// OpenVPNConf is the file path to the OpenVPN client configuration file.
	OpenVPNConf string = "/etc/openvpn/target.ovpn"
	// PIAPortForward is the file path to the port forwarding JSON information for PIA servers.
	PIAPortForward string = "/gluetun/piaportforward.json"
	// TunnelDevice is the file path to tun device.
	TunnelDevice string = "/dev/net/tun"
	// NetRoute is the path to the file containing information on the network route.
	NetRoute string = "/proc/net/route"
	// RootHints is the filepath to the root.hints file used by Unbound.
	RootHints string = "/etc/unbound/root.hints"
	// RootKey is the filepath to the root.key file used by Unbound.
	RootKey string = "/etc/unbound/root.key"
	// Client key filepath, used by Cyberghost.
	ClientKey string = "/gluetun/client.key"
	// Client certificate filepath, used by Cyberghost.
	ClientCertificate string = "/gluetun/client.crt"
	// Servers information filepath.
	ServersData = "/gluetun/servers.json"
)
