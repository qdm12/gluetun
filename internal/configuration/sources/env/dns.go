package env

import (
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *Source) readDNS() (dns settings.DNS, err error) {
	dns.ServerAddress, err = s.readDNSServerAddress()
	if err != nil {
		return dns, err
	}

	dns.KeepNameserver, err = envToBoolPtr("DNS_KEEP_NAMESERVER")
	if err != nil {
		return dns, fmt.Errorf("environment variable DNS_KEEP_NAMESERVER: %w", err)
	}

	dns.DoT, err = s.readDoT()
	if err != nil {
		return dns, fmt.Errorf("DoT settings: %w", err)
	}

	return dns, nil
}

func (s *Source) readDNSServerAddress() (address net.IP, err error) {
	key, value := s.getEnvWithRetro("DNS_ADDRESS", "DNS_PLAINTEXT_ADDRESS")
	if value == "" {
		return nil, nil
	}

	address = net.ParseIP(value)
	if address == nil {
		return nil, fmt.Errorf("environment variable %s: %w: %s", key, ErrIPAddressParse, value)
	}

	// TODO remove in v4
	if !address.Equal(net.IPv4(127, 0, 0, 1)) { //nolint:gomnd
		s.warner.Warn(key + " is set to " + value +
			" so the DNS over TLS (DoT) server will not be used." +
			" The default value changed to 127.0.0.1 so it uses the internal DoT serves." +
			" If the DoT server fails to start, the IPv4 address of the first plaintext DNS server" +
			" corresponding to the first DoT provider chosen is used.")
	}

	return address, nil
}
