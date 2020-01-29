package constants

const (
	// UnboundConf is the file path to the Unbound configuration file
	UnboundConf = "/etc/unbound/unbound.conf"
	// ResolvConf is the file path to the system resolv.conf file
	ResolvConf = "/etc/resolv.conf"
	// OpenVPNAuthConf is the file path to the OpenVPN auth file
	OpenVPNAuthConf = "/etc/openvpn/auth.conf"
	// OpenVPNConf is the file path to the OpenVPN client configuration file
	OpenVPNConf = "/etc/openvpn/target.ovpn"
	// TunnelDevice is the file path to tun device
	TunnelDevice = "/dev/net/tun"
	// NetRoute is the path to the file containing information on the network route
	NetRoute = "/proc/net/route"
)
