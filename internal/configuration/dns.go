package configuration

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	unboundmodels "github.com/qdm12/dns/pkg/models"
	unbound "github.com/qdm12/dns/pkg/unbound"
	"github.com/qdm12/golibs/params"
)

// DNS contains settings to configure Unbound for DNS over TLS operation.
type DNS struct { //nolint:maligned
	Enabled           bool
	PlaintextAddress  net.IP
	KeepNameserver    bool
	BlockMalicious    bool
	BlockAds          bool
	BlockSurveillance bool
	UpdatePeriod      time.Duration
	Unbound           unboundmodels.Settings
}

func (settings *DNS) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *DNS) lines() (lines []string) {
	lines = append(lines, lastIndent+"DNS:")

	if settings.PlaintextAddress != nil {
		lines = append(lines, indent+lastIndent+"Plaintext address: "+settings.PlaintextAddress.String())
	}

	if settings.KeepNameserver {
		lines = append(lines, indent+lastIndent+"Keep nameserver (disabled blocking): yes")
	}

	if !settings.Enabled {
		return lines
	}

	lines = append(lines, indent+lastIndent+"DNS over TLS:")

	lines = append(lines, indent+indent+lastIndent+"Unbound:")
	for _, line := range settings.Unbound.Lines() {
		lines = append(lines, indent+indent+indent+line)
	}

	if settings.BlockMalicious {
		lines = append(lines, indent+indent+lastIndent+"Block malicious: enabled")
	}

	if settings.BlockAds {
		lines = append(lines, indent+indent+lastIndent+"Block ads: enabled")
	}

	if settings.BlockSurveillance {
		lines = append(lines, indent+indent+lastIndent+"Block surveillance: enabled")
	}

	if settings.UpdatePeriod > 0 {
		lines = append(lines, indent+indent+lastIndent+"Update: every "+settings.UpdatePeriod.String())
	}

	return lines
}

var (
	ErrUnboundSettings   = errors.New("failed getting Unbound settings")
	ErrDNSProviderNoData = errors.New("DNS provider has no associated data")
	ErrDNSProviderNoTLS  = errors.New("DNS provider does not support DNS over TLS")
	ErrDNSNoIPv6Support  = errors.New("no DNS provider supports IPv6")
)

func (settings *DNS) read(r reader) (err error) {
	settings.Enabled, err = r.env.OnOff("DOT", params.Default("on"))
	if err != nil {
		return err
	}

	// Plain DNS settings
	if err := settings.readDNSPlaintext(r.env); err != nil {
		return err
	}
	settings.KeepNameserver, err = r.env.OnOff("DNS_KEEP_NAMESERVER", params.Default("off"))
	if err != nil {
		return err
	}

	// DNS over TLS external settings
	settings.BlockMalicious, err = r.env.OnOff("BLOCK_MALICIOUS", params.Default("on"))
	if err != nil {
		return err
	}
	settings.BlockSurveillance, err = r.env.OnOff("BLOCK_SURVEILLANCE", params.Default("on"),
		params.RetroKeys([]string{"BLOCK_NSA"}, r.onRetroActive))
	if err != nil {
		return err
	}
	settings.BlockAds, err = r.env.OnOff("BLOCK_ADS", params.Default("off"))
	if err != nil {
		return err
	}
	settings.UpdatePeriod, err = r.env.Duration("DNS_UPDATE_PERIOD", params.Default("24h"))
	if err != nil {
		return err
	}

	if err := settings.readUnbound(r); err != nil {
		return fmt.Errorf("%w: %s", ErrUnboundSettings, err)
	}

	// Consistency check
	IPv6Support := false
	for _, provider := range settings.Unbound.Providers {
		providerData, ok := unbound.GetProviderData(provider)
		switch {
		case !ok:
			return fmt.Errorf("%w: %s", ErrDNSProviderNoData, provider)
		case !providerData.SupportsTLS:
			return fmt.Errorf("%w: %s", ErrDNSProviderNoTLS, provider)
		case providerData.SupportsIPv6:
			IPv6Support = true
		}
	}

	if settings.Unbound.IPv6 && !IPv6Support {
		return ErrDNSNoIPv6Support
	}

	return nil
}

var (
	ErrDNSAddressNotAnIP = errors.New("DNS plaintext address is not an IP address")
)

func (settings *DNS) readDNSPlaintext(env params.Env) error {
	s, err := env.Get("DNS_PLAINTEXT_ADDRESS", params.Default("1.1.1.1"))
	if err != nil {
		return err
	}

	settings.PlaintextAddress = net.ParseIP(s)
	if settings.PlaintextAddress == nil {
		return fmt.Errorf("%w: %s", ErrDNSAddressNotAnIP, s)
	}

	return nil
}
