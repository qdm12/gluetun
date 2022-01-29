package env

import (
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/govalid/binary"
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
	key, value := r.getEnvWithRetro("BLOCK_SURVEILLANCE", "BLOCK_NSA")
	if value == "" {
		return nil, nil //nolint:nilnil
	}

	blocked = new(bool)
	*blocked, err = binary.Validate(key)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s: %w", key, err)
	}

	return blocked, nil
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
