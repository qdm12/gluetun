package models

import (
	"fmt"
	"strings"
)

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
	VPNProvider string
	// NetworkProtocol contains the network protocol to be used to communicate with the VPN servers
	NetworkProtocol string
)

func marshalJSONString(s string) (data []byte, err error) {
	return []byte(fmt.Sprintf("%q", s)), nil
}

func unmarshalJSONString(data []byte) (s string) {
	s = string(data)
	s = strings.TrimPrefix(s, "\"")
	s = strings.TrimSuffix(s, "\"")
	return s
}

func (v *VPNProvider) MarshalJSON() ([]byte, error) {
	return marshalJSONString(string(*v))
}

func (v *VPNProvider) UnmarshalJSON(data []byte) error {
	*v = VPNProvider(unmarshalJSONString(data))
	return nil
}

func (n *NetworkProtocol) MarshalJSON() ([]byte, error) {
	return marshalJSONString(string(*n))
}

func (n *NetworkProtocol) UnmarshalJSON(data []byte) error {
	*n = NetworkProtocol(unmarshalJSONString(data))
	return nil
}

func (f *Filepath) MarshalJSON() ([]byte, error) {
	return marshalJSONString(string(*f))
}

func (f *Filepath) UnmarshalJSON(data []byte) error {
	*f = Filepath(unmarshalJSONString(data))
	return nil
}
