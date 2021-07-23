package dns

import "github.com/qdm12/dns/pkg/nameserver"

func (l *looper) useUnencryptedDNS(fallback bool) {
	settings := l.GetSettings()

	// Try with user provided plaintext ip address
	targetIP := settings.PlaintextAddress
	if targetIP != nil {
		if fallback {
			l.logger.Info("falling back on plaintext DNS at address " + targetIP.String())
		} else {
			l.logger.Info("using plaintext DNS at address " + targetIP.String())
		}
		nameserver.UseDNSInternally(targetIP)
		err := nameserver.UseDNSSystemWide(l.resolvConf, targetIP, settings.KeepNameserver)
		if err != nil {
			l.logger.Error(err.Error())
		}
		return
	}

	provider := settings.Unbound.Providers[0]
	targetIP = provider.DoT().IPv4[0]
	if fallback {
		l.logger.Info("falling back on plaintext DNS at address " + targetIP.String())
	} else {
		l.logger.Info("using plaintext DNS at address " + targetIP.String())
	}
	nameserver.UseDNSInternally(targetIP)
	err := nameserver.UseDNSSystemWide(l.resolvConf, targetIP, settings.KeepNameserver)
	if err != nil {
		l.logger.Error(err.Error())
	}
}
