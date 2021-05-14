package configuration

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/qdm12/dns/pkg/blacklist"
	"github.com/qdm12/dns/pkg/unbound"
	"github.com/qdm12/golibs/params"
)

// DNS contains settings to configure Unbound for DNS over TLS operation.
type DNS struct { //nolint:maligned
	Enabled          bool
	PlaintextAddress net.IP
	KeepNameserver   bool
	UpdatePeriod     time.Duration
	Unbound          unbound.Settings
	BlacklistBuild   blacklist.BuilderSettings
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

	lines = append(lines, indent+indent+lastIndent+"Blacklist:")
	for _, line := range settings.BlacklistBuild.Lines(indent, lastIndent) {
		lines = append(lines, indent+indent+indent+line)
	}

	if settings.UpdatePeriod > 0 {
		lines = append(lines, indent+indent+lastIndent+"Update: every "+settings.UpdatePeriod.String())
	}

	return lines
}

var (
	ErrUnboundSettings   = errors.New("failed getting Unbound settings")
	ErrBlacklistSettings = errors.New("failed getting DNS blacklist settings")
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
	if err := settings.readBlacklistBuilding(r); err != nil {
		return fmt.Errorf("%w: %s", ErrBlacklistSettings, err)
	}

	settings.UpdatePeriod, err = r.env.Duration("DNS_UPDATE_PERIOD", params.Default("24h"))
	if err != nil {
		return err
	}

	// Unbound settings
	if err := settings.readUnbound(r); err != nil {
		return fmt.Errorf("%w: %s", ErrUnboundSettings, err)
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
