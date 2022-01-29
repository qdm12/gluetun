package env

import (
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"inet.af/netaddr"
)

func (r *Reader) readDNSBlacklist() (blacklist settings.DNSBlacklist, err error) {
	blacklist.BlockMalicious, err = envToBoolPtr("BLOCK_MALICIOUS")
	if err != nil {
		return blacklist, fmt.Errorf("environment variable BLOCK_MALICIOUS: %w", err)
	}

	blacklist.BlockSurveillance, err = r.readBlockSurveillance()
	if err != nil {
		return blacklist, fmt.Errorf("environment variable BLOCK_MALICIOUS: %w", err)
	}

	blacklist.BlockAds, err = envToBoolPtr("BLOCK_ADS")
	if err != nil {
		return blacklist, fmt.Errorf("environment variable BLOCK_ADS: %w", err)
	}

	blacklist.AddBlockedIPs, blacklist.AddBlockedIPPrefixes,
		err = readDoTPrivateAddresses() // TODO v4 split in 2
	if err != nil {
		return blacklist, err
	}

	blacklist.AllowedHosts = envToCSV("UNBLOCK") // TODO v4 change name

	return blacklist, nil
}

func (r *Reader) readBlockSurveillance() (blocked *bool, err error) {
	blocked, err = envToBoolPtr("BLOCK_NSA")
	if err != nil {
		r.onRetroActive("BLOCK_NSA", "BLOCK_SURVEILLANCE")
		return nil, fmt.Errorf("environment variable BLOCK_NSA: %w", err)
	} else if blocked != nil {
		r.onRetroActive("BLOCK_NSA", "BLOCK_SURVEILLANCE")
		return blocked, nil
	}

	blocked, err = envToBoolPtr("BLOCK_SURVEILLANCE")
	if err != nil {
		return nil, fmt.Errorf("environment variable BLOCK_SURVEILLANCE: %w", err)
	}
		return blocked, nil
	}

	return nil, nil //nolint:nilnil
}

var (
	ErrPrivateAddressNotValid = errors.New("private address is not a valid IP or CIDR range")
)

func readDoTPrivateAddresses() (ips []netaddr.IP,
	ipPrefixes []netaddr.IPPrefix, err error) {
	privateAddresses := envToCSV("DOT_PRIVATE_ADDRESS")
	if len(privateAddresses) == 0 {
		return nil, nil, nil
	}

	ips = make([]netaddr.IP, 0, len(privateAddresses))
	ipPrefixes = make([]netaddr.IPPrefix, 0, len(privateAddresses))

	for _, privateAddress := range privateAddresses {
		ip, err := netaddr.ParseIP(privateAddress)
		if err == nil {
			ips = append(ips, ip)
			continue
		}

		ipPrefix, err := netaddr.ParseIPPrefix(privateAddress)
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
