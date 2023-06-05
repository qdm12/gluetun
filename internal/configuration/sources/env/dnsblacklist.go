package env

import (
	"errors"
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readDNSBlacklist() (blacklist settings.DNSBlacklist, err error) {
	blacklist.BlockMalicious, err = s.env.BoolPtr("BLOCK_MALICIOUS")
	if err != nil {
		return blacklist, err
	}

	blacklist.BlockSurveillance, err = s.env.BoolPtr("BLOCK_SURVEILLANCE",
		env.RetroKeys("BLOCK_NSA"))
	if err != nil {
		return blacklist, err
	}

	blacklist.BlockAds, err = s.env.BoolPtr("BLOCK_ADS")
	if err != nil {
		return blacklist, err
	}

	blacklist.AddBlockedIPs, blacklist.AddBlockedIPPrefixes,
		err = s.readDoTPrivateAddresses() // TODO v4 split in 2
	if err != nil {
		return blacklist, err
	}

	blacklist.AllowedHosts = s.env.CSV("UNBLOCK") // TODO v4 change name

	return blacklist, nil
}

var (
	ErrPrivateAddressNotValid = errors.New("private address is not a valid IP or CIDR range")
)

func (s *Source) readDoTPrivateAddresses() (ips []netip.Addr,
	ipPrefixes []netip.Prefix, err error) {
	privateAddresses := s.env.CSV("DOT_PRIVATE_ADDRESS")
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
