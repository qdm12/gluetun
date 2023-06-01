package env

import (
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *Source) readDNS() (dns settings.DNS, err error) {
	dns.ServerAddress, err = s.readDNSServerAddress()
	if err != nil {
		return dns, err
	}

	dns.KeepNameserver, err = s.env.BoolPtr("DNS_KEEP_NAMESERVER")
	if err != nil {
		return dns, err
	}

	dns.DoT, err = s.readDoT()
	if err != nil {
		return dns, fmt.Errorf("DoT settings: %w", err)
	}

	return dns, nil
}

func (s *Source) readDNSServerAddress() (address netip.Addr, err error) {
	key, value := s.getEnvWithRetro("DNS_ADDRESS", []string{"DNS_PLAINTEXT_ADDRESS"})
	if value == nil {
		return address, nil
	}

	address, err = netip.ParseAddr(*value)
	if err != nil {
		return address, fmt.Errorf("environment variable %s: %w", key, err)
	}

	// TODO remove in v4
	if address.Unmap().Compare(netip.AddrFrom4([4]byte{127, 0, 0, 1})) != 0 {
		s.warner.Warn(key + " is set to " + *value +
			" so the DNS over TLS (DoT) server will not be used." +
			" The default value changed to 127.0.0.1 so it uses the internal DoT serves." +
			" If the DoT server fails to start, the IPv4 address of the first plaintext DNS server" +
			" corresponding to the first DoT provider chosen is used.")
	}

	return address, nil
}
