package configuration

import (
	"errors"
	"fmt"

	"github.com/qdm12/golibs/params"
	"inet.af/netaddr"
)

func (settings *DNS) readBlacklistBuilding(r reader) (err error) {
	settings.BlacklistBuild.BlockMalicious, err = r.env.OnOff("BLOCK_MALICIOUS", params.Default("on"))
	if err != nil {
		return fmt.Errorf("environment variable BLOCK_MALICIOUS: %w", err)
	}

	settings.BlacklistBuild.BlockSurveillance, err = r.env.OnOff("BLOCK_SURVEILLANCE", params.Default("on"),
		params.RetroKeys([]string{"BLOCK_NSA"}, r.onRetroActive))
	if err != nil {
		return fmt.Errorf("environment variable BLOCK_SURVEILLANCE (or BLOCK_NSA): %w", err)
	}

	settings.BlacklistBuild.BlockAds, err = r.env.OnOff("BLOCK_ADS", params.Default("off"))
	if err != nil {
		return fmt.Errorf("environment variable BLOCK_ADS: %w", err)
	}

	if err := settings.readPrivateAddresses(r.env); err != nil {
		return err
	}

	return settings.readBlacklistUnblockedHostnames(r)
}

var (
	ErrInvalidPrivateAddress = errors.New("private address is not a valid IP or CIDR range")
)

func (settings *DNS) readPrivateAddresses(env params.Interface) (err error) {
	privateAddresses, err := env.CSV("DOT_PRIVATE_ADDRESS")
	if err != nil {
		return fmt.Errorf("environment variable DOT_PRIVATE_ADDRESS: %w", err)
	} else if len(privateAddresses) == 0 {
		return nil
	}

	ips := make([]netaddr.IP, 0, len(privateAddresses))
	ipPrefixes := make([]netaddr.IPPrefix, 0, len(privateAddresses))

	for _, address := range privateAddresses {
		ip, err := netaddr.ParseIP(address)
		if err == nil {
			ips = append(ips, ip)
			continue
		}

		ipPrefix, err := netaddr.ParseIPPrefix(address)
		if err == nil {
			ipPrefixes = append(ipPrefixes, ipPrefix)
			continue
		}

		return fmt.Errorf("%w: %s", ErrInvalidPrivateAddress, address)
	}

	settings.BlacklistBuild.AddBlockedIPs = append(settings.BlacklistBuild.AddBlockedIPs, ips...)
	settings.BlacklistBuild.AddBlockedIPPrefixes = append(settings.BlacklistBuild.AddBlockedIPPrefixes, ipPrefixes...)

	return nil
}

func (settings *DNS) readBlacklistUnblockedHostnames(r reader) (err error) {
	hostnames, err := r.env.CSV("UNBLOCK")
	if err != nil {
		return fmt.Errorf("environment variable UNBLOCK: %w", err)
	} else if len(hostnames) == 0 {
		return nil
	}
	for _, hostname := range hostnames {
		if !r.regex.MatchHostname(hostname) {
			return fmt.Errorf("%w: %s", ErrInvalidHostname, hostname)
		}
	}

	settings.BlacklistBuild.AllowedHosts = append(settings.BlacklistBuild.AllowedHosts, hostnames...)
	return nil
}
