package env

import (
	"errors"
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
	"github.com/qdm12/govalid/binary"
)

func (s *Source) readDNSBlacklist() (blacklist settings.DNSBlacklist, err error) {
	blacklist.BlockMalicious, err = env.BoolPtr("BLOCK_MALICIOUS")
	if err != nil {
		return blacklist, fmt.Errorf("environment variable BLOCK_MALICIOUS: %w", err)
	}

	blacklist.BlockSurveillance, err = s.readBlockSurveillance()
	if err != nil {
		return blacklist, err
	}

	blacklist.BlockAds, err = env.BoolPtr("BLOCK_ADS")
	if err != nil {
		return blacklist, fmt.Errorf("environment variable BLOCK_ADS: %w", err)
	}

	blacklist.AddBlockedIPs, blacklist.AddBlockedIPPrefixes,
		err = readDoTPrivateAddresses() // TODO v4 split in 2
	if err != nil {
		return blacklist, err
	}

	blacklist.AllowedHosts = env.CSV("UNBLOCK") // TODO v4 change name

	return blacklist, nil
}

func (s *Source) readBlockSurveillance() (blocked *bool, err error) {
	key, value := s.getEnvWithRetro("BLOCK_SURVEILLANCE", []string{"BLOCK_NSA"})
	blocked, err = binary.Validate(value)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s: %w", key, err)
	}

	return blocked, nil
}

var (
	ErrPrivateAddressNotValid = errors.New("private address is not a valid IP or CIDR range")
)

func readDoTPrivateAddresses() (ips []netip.Addr,
	ipPrefixes []netip.Prefix, err error) {
	privateAddresses := env.CSV("DOT_PRIVATE_ADDRESS")
	if len(privateAddresses) == 0 {
		return nil, nil, nil
	}

	ips = make([]netip.Addr, 0, len(privateAddresses))
	ipPrefixes = make([]netip.Prefix, 0, len(privateAddresses))

	for _, privateAddress := range privateAddresses {
		ip, err := netip.ParseAddr(privateAddress)
		if err == nil {
			ips = append(ips, ip)
			continue
		}

		ipPrefix, err := netip.ParsePrefix(privateAddress)
		if err == nil {
			ipPrefixes = append(ipPrefixes, ipPrefix)
			continue
		}

		return nil, nil, fmt.Errorf(
			"environment variable DOT_PRIVATE_ADDRESS: %w: %s",
			ErrPrivateAddressNotValid, privateAddress)
	}

	return ips, ipPrefixes, nil
}
