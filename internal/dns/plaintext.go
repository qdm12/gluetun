package dns

import (
	"net/netip"
	"time"

	"github.com/qdm12/dns/v2/pkg/nameserver"
)

func (l *Loop) useUnencryptedDNS(fallback bool) {
	settings := l.GetSettings()

	// Try with user provided plaintext ip address
	// if it's not 127.0.0.1 (default for DoT), otherwise
	// use the first DoT provider ipv4 address found.
	var targetIP netip.Addr
	if settings.ServerAddress.Compare(netip.AddrFrom4([4]byte{127, 0, 0, 1})) != 0 {
		targetIP = settings.ServerAddress
	} else {
		targetIP = settings.DoT.GetFirstPlaintextIPv4()
	}

	if fallback {
		l.logger.Info("falling back on plaintext DNS at address " + targetIP.String())
	} else {
		l.logger.Info("using plaintext DNS at address " + targetIP.String())
	}

	const dialTimeout = 3 * time.Second
	settingsInternalDNS := nameserver.SettingsInternalDNS{
		IP:      targetIP,
		Timeout: dialTimeout,
	}
	nameserver.UseDNSInternally(settingsInternalDNS)

	settingsSystemWide := nameserver.SettingsSystemDNS{
		IP:         targetIP,
		ResolvPath: l.resolvConf,
	}
	err := nameserver.UseDNSSystemWide(settingsSystemWide)
	if err != nil {
		l.logger.Error(err.Error())
	}
}
