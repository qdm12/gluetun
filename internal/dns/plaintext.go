package dns

import (
	"time"

	"github.com/qdm12/dns/v2/pkg/nameserver"
)

func (l *Loop) useUnencryptedDNS(fallback bool) {
	settings := l.GetSettings()

	targetIP := settings.GetFirstPlaintextIPv4()

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
