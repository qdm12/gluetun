package models

type (
	// VPNDevice is the device name used to tunnel using Openvpn
	VPNDevice string
	// DNSProvider is a DNS over TLS server provider name
	DNSProvider string
	// DNSHost is the DNS host to use for TLS validation
	DNSHost string
	// URL is an HTTP(s) URL address
	URL string
	// Filepath is a local filesytem file path
	Filepath string
	// TinyProxyLogLevel is the log level for TinyProxy
	TinyProxyLogLevel string
	// VPNProvider is the name of the VPN provider to be used
	VPNProvider string // TODO
	// NetworkProtocol contains the network protocol to be used to communicate with the VPN servers
	NetworkProtocol string
)
