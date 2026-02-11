package dns

import (
	"context"
	"net/netip"
	"time"

	"github.com/qdm12/dns/v2/pkg/nameserver"
)

func (l *Loop) useUnencryptedDNS(ctx context.Context, fallback bool) {
	settings := l.GetSettings()

	targetIP := settings.GetFirstPlaintextIPv4()

	if fallback {
		l.logger.Info("falling back on plaintext DNS at address " + targetIP.String())
	} else {
		l.logger.Info("using plaintext DNS at address " + targetIP.String())
	}

	const dialTimeout = 3 * time.Second
	const defaultDNSPort = 53
	addrPort := netip.AddrPortFrom(targetIP, defaultDNSPort)
	settingsInternalDNS := nameserver.SettingsInternalDNS{
		AddrPort: addrPort,
		Timeout:  dialTimeout,
	}
	nameserver.UseDNSInternally(settingsInternalDNS)

	settingsSystemWide := nameserver.SettingsSystemDNS{
		IPs:        []netip.Addr{targetIP},
		ResolvPath: l.resolvConf,
	}
	err := nameserver.UseDNSSystemWide(settingsSystemWide)
	if err != nil {
		l.logger.Error(err.Error())
	}

	err = l.firewall.RestrictOutputAddrPort(ctx, addrPort)
	if err != nil {
		l.logger.Error("restricting plain DNS traffic to " + targetIP.String() + ": " + err.Error())
	}
}
