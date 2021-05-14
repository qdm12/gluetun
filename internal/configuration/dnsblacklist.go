package configuration

import (
	"errors"
	"fmt"
	"net"

	"github.com/qdm12/golibs/params"
)

func (settings *DNS) readBlacklistBuilding(r reader) (err error) {
	settings.BlacklistBuild.BlockMalicious, err = r.env.OnOff("BLOCK_MALICIOUS", params.Default("on"))
	if err != nil {
		return err
	}

	settings.BlacklistBuild.BlockSurveillance, err = r.env.OnOff("BLOCK_SURVEILLANCE", params.Default("on"),
		params.RetroKeys([]string{"BLOCK_NSA"}, r.onRetroActive))
	if err != nil {
		return err
	}

	settings.BlacklistBuild.BlockAds, err = r.env.OnOff("BLOCK_ADS", params.Default("off"))
	if err != nil {
		return err
	}

	if err := settings.readPrivateAddresses(r.env); err != nil {
		return err
	}

	if err := settings.readBlacklistUnblockedHostnames(r); err != nil {
		return err
	}

	return nil
}

var (
	ErrInvalidPrivateAddress = errors.New("private address is not a valid IP or CIDR range")
)

func (settings *DNS) readPrivateAddresses(env params.Env) (err error) {
	privateAddresses, err := env.CSV("DOT_PRIVATE_ADDRESS")
	if err != nil {
		return err
	} else if len(privateAddresses) == 0 {
		return nil
	}

	ips := make([]net.IP, 0, len(privateAddresses))
	ipNets := make([]*net.IPNet, 0, len(privateAddresses))

	for _, address := range privateAddresses {
		ip := net.ParseIP(address)
		if ip != nil {
			ips = append(ips, ip)
			continue
		}

		_, ipNet, err := net.ParseCIDR(address)
		if err == nil && ipNet != nil {
			ipNets = append(ipNets, ipNet)
			continue
		}

		return fmt.Errorf("%w: %s", ErrInvalidPrivateAddress, address)
	}

	settings.BlacklistBuild.AddBlockedIPs = append(settings.BlacklistBuild.AddBlockedIPs, ips...)
	settings.BlacklistBuild.AddBlockedIPNets = append(settings.BlacklistBuild.AddBlockedIPNets, ipNets...)

	return nil
}

func (settings *DNS) readBlacklistUnblockedHostnames(r reader) (err error) {
	hostnames, err := r.env.CSV("UNBLOCK")
	if err != nil {
		return err
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
