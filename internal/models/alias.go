package models

type (
	// VPNDevice is the device name used to tunnel using Openvpn
	VPNDevice string
	// DNSProvider is a DNS over TLS server provider name
	DNSProvider string
	// DNSHost is the DNS host to use for TLS validation
	DNSHost string
	// PIAEncryption defines the level of encryption for communication with PIA servers
	PIAEncryption string
	// PIARegion is used to define the list of regions available for PIA
	PIARegion string
	// MullvadCountry is used as the country for a Mullvad server
	MullvadCountry string
	// MullvadCity is used as the city for a Mullvad server
	MullvadCity string
	// MullvadProvider is used as the Internet service provider for a Mullvad server
	MullvadProvider string
	// WindscribeCity is used as the region for a Windscribe server
	WindscribeRegion string
	// URL is an HTTP(s) URL address
	URL string
	// Filepath is a local filesytem file path
	Filepath string
	// TinyProxyLogLevel is the log level for TinyProxy
	TinyProxyLogLevel string
	// VPNProvider is the name of the VPN provider to be used
	VPNProvider string
	// NetworkProtocol contains the network protocol to be used to communicate with the VPN servers
	NetworkProtocol string
)
