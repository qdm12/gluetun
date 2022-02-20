package env

import (
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (r *Reader) readDNS() (dns settings.DNS, err error) {
	dns.ServerAddress, err = r.readDNSServerAddress()
	if err != nil {
		return dns, err
	}

	dns.KeepNameserver, err = envToBoolPtr("DNS_KEEP_NAMESERVER")
	if err != nil {
		return dns, fmt.Errorf("environment variable DNS_KEEP_NAMESERVER: %w", err)
	}

	dns.DoT, err = r.readDoT()
	if err != nil {
		return dns, fmt.Errorf("DoT settings: %w", err)
	}

	return dns, nil
}

func (r *Reader) readDNSServerAddress() (address net.IP, err error) {
	key, s := r.getEnvWithRetro("DNS_ADDRESS", "DNS_PLAINTEXT_ADDRESS")
	if s == "" {
		return nil, nil
	}

	address = net.ParseIP(s)
	if address == nil {
		return nil, fmt.Errorf("environment variable %s: %w: %s", key, ErrIPAddressParse, s)
	}

	// TODO remove in v4
	if !address.Equal(net.IPv4(127, 0, 0, 1)) { //nolint:gomnd
		r.warner.Warn(key + " is set to " + s +
			" so the DNS over TLS (DoT) server will not be used." +
			" The default value changed to 127.0.0.1 so it uses the internal DoT server." +
			" If the DoT server fails to start, the IPv4 address of the first plaintext DNS server" +
			" corresponding to the first DoT provider chosen is used.")
	}

	return address, nil
}
