package env

import (
	"fmt"
	"net"
	"os"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (r *Reader) readDNS() (dns settings.DNS, err error) {
	dns.ServerAddress, err = readDNSServerAddress()
	if err != nil {
		return dns, err
	}

	dns.KeepNameserver, err = envToBoolPtr("DNS_KEEP_NAMESERVER")
	if err != nil {
		return dns, fmt.Errorf("environment variable DNS_KEEP_NAMESERVER: %w", err)
	}

	dns.DoT, err = r.readDoT()
	if err != nil {
		return dns, fmt.Errorf("cannot read DoT settings: %w", err)
	}

	return dns, nil
}

func readDNSServerAddress() (address net.IP, err error) {
	s := os.Getenv("DNS_PLAINTEXT_ADDRESS")
	if s == "" {
		return nil, nil
	}

	address = net.ParseIP(s)
	if address == nil {
		return nil, fmt.Errorf("environment variable DNS_PLAINTEXT_ADDRESS: %w: %s", ErrIPAddressParse, s)
	}

	return address, nil
}
