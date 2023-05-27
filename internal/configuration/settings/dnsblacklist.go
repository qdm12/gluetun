package settings

import (
	"errors"
	"fmt"
	"net/netip"
	"regexp"

	"github.com/qdm12/dns/pkg/blacklist"
	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gotree"
)

// DNSBlacklist is settings for the DNS blacklist building.
type DNSBlacklist struct {
	BlockMalicious       *bool
	BlockAds             *bool
	BlockSurveillance    *bool
	AllowedHosts         []string
	AddBlockedHosts      []string
	AddBlockedIPs        []netip.Addr
	AddBlockedIPPrefixes []netip.Prefix
}

func (b *DNSBlacklist) setDefaults() {
	b.BlockMalicious = gosettings.DefaultPointer(b.BlockMalicious, true)
	b.BlockAds = gosettings.DefaultPointer(b.BlockAds, false)
	b.BlockSurveillance = gosettings.DefaultPointer(b.BlockSurveillance, true)
}

var hostRegex = regexp.MustCompile(`^([a-zA-Z0-9]|[a-zA-Z0-9_][a-zA-Z0-9\-_]{0,61}[a-zA-Z0-9_])(\.([a-zA-Z0-9]|[a-zA-Z0-9_][a-zA-Z0-9\-_]{0,61}[a-zA-Z0-9]))*$`) //nolint:lll

var (
	ErrAllowedHostNotValid = errors.New("allowed host is not valid")
	ErrBlockedHostNotValid = errors.New("blocked host is not valid")
)

func (b DNSBlacklist) validate() (err error) {
	for _, host := range b.AllowedHosts {
		if !hostRegex.MatchString(host) {
			return fmt.Errorf("%w: %s", ErrAllowedHostNotValid, host)
		}
	}

	for _, host := range b.AddBlockedHosts {
		if !hostRegex.MatchString(host) {
			return fmt.Errorf("%w: %s", ErrBlockedHostNotValid, host)
		}
	}

	return nil
}

func (b DNSBlacklist) copy() (copied DNSBlacklist) {
	return DNSBlacklist{
		BlockMalicious:       gosettings.CopyPointer(b.BlockMalicious),
		BlockAds:             gosettings.CopyPointer(b.BlockAds),
		BlockSurveillance:    gosettings.CopyPointer(b.BlockSurveillance),
		AllowedHosts:         gosettings.CopySlice(b.AllowedHosts),
		AddBlockedHosts:      gosettings.CopySlice(b.AddBlockedHosts),
		AddBlockedIPs:        gosettings.CopySlice(b.AddBlockedIPs),
		AddBlockedIPPrefixes: gosettings.CopySlice(b.AddBlockedIPPrefixes),
	}
}

func (b *DNSBlacklist) mergeWith(other DNSBlacklist) {
	b.BlockMalicious = gosettings.MergeWithPointer(b.BlockMalicious, other.BlockMalicious)
	b.BlockAds = gosettings.MergeWithPointer(b.BlockAds, other.BlockAds)
	b.BlockSurveillance = gosettings.MergeWithPointer(b.BlockSurveillance, other.BlockSurveillance)
	b.AllowedHosts = gosettings.MergeWithSlice(b.AllowedHosts, other.AllowedHosts)
	b.AddBlockedHosts = gosettings.MergeWithSlice(b.AddBlockedHosts, other.AddBlockedHosts)
	b.AddBlockedIPs = gosettings.MergeWithSlice(b.AddBlockedIPs, other.AddBlockedIPs)
	b.AddBlockedIPPrefixes = gosettings.MergeWithSlice(b.AddBlockedIPPrefixes, other.AddBlockedIPPrefixes)
}

func (b *DNSBlacklist) overrideWith(other DNSBlacklist) {
	b.BlockMalicious = gosettings.OverrideWithPointer(b.BlockMalicious, other.BlockMalicious)
	b.BlockAds = gosettings.OverrideWithPointer(b.BlockAds, other.BlockAds)
	b.BlockSurveillance = gosettings.OverrideWithPointer(b.BlockSurveillance, other.BlockSurveillance)
	b.AllowedHosts = gosettings.OverrideWithSlice(b.AllowedHosts, other.AllowedHosts)
	b.AddBlockedHosts = gosettings.OverrideWithSlice(b.AddBlockedHosts, other.AddBlockedHosts)
	b.AddBlockedIPs = gosettings.OverrideWithSlice(b.AddBlockedIPs, other.AddBlockedIPs)
	b.AddBlockedIPPrefixes = gosettings.OverrideWithSlice(b.AddBlockedIPPrefixes, other.AddBlockedIPPrefixes)
}

func (b DNSBlacklist) ToBlacklistFormat() (settings blacklist.BuilderSettings, err error) {
	return blacklist.BuilderSettings{
		BlockMalicious:       *b.BlockMalicious,
		BlockAds:             *b.BlockAds,
		BlockSurveillance:    *b.BlockSurveillance,
		AllowedHosts:         b.AllowedHosts,
		AddBlockedHosts:      b.AddBlockedHosts,
		AddBlockedIPs:        netipAddressesToNetaddrIPs(b.AddBlockedIPs),
		AddBlockedIPPrefixes: netipPrefixesToNetaddrIPPrefixes(b.AddBlockedIPPrefixes),
	}, nil
}

func (b DNSBlacklist) String() string {
	return b.toLinesNode().String()
}

func (b DNSBlacklist) toLinesNode() (node *gotree.Node) {
	node = gotree.New("DNS filtering settings:")

	node.Appendf("Block malicious: %s", helpers.BoolPtrToYesNo(b.BlockMalicious))
	node.Appendf("Block ads: %s", helpers.BoolPtrToYesNo(b.BlockAds))
	node.Appendf("Block surveillance: %s", helpers.BoolPtrToYesNo(b.BlockSurveillance))

	if len(b.AllowedHosts) > 0 {
		allowedHostsNode := node.Appendf("Allowed hosts:")
		for _, host := range b.AllowedHosts {
			allowedHostsNode.Appendf(host)
		}
	}

	if len(b.AddBlockedHosts) > 0 {
		blockedHostsNode := node.Appendf("Blocked hosts:")
		for _, host := range b.AddBlockedHosts {
			blockedHostsNode.Appendf(host)
		}
	}

	if len(b.AddBlockedIPs) > 0 {
		blockedIPsNode := node.Appendf("Blocked IP addresses:")
		for _, ip := range b.AddBlockedIPs {
			blockedIPsNode.Appendf(ip.String())
		}
	}

	if len(b.AddBlockedIPPrefixes) > 0 {
		blockedIPPrefixesNode := node.Appendf("Blocked IP networks:")
		for _, ipNetwork := range b.AddBlockedIPPrefixes {
			blockedIPPrefixesNode.Appendf(ipNetwork.String())
		}
	}

	return node
}
