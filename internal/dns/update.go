package dns

import "context"

func (l *Loop) updateFiles(ctx context.Context) (err error) {
	l.logger.Info("downloading DNS over TLS cryptographic files")
	if err := l.conf.SetupFiles(ctx); err != nil {
		return err
	}
	settings := l.GetSettings()

	unboundSettings, err := settings.DoT.Unbound.ToUnboundFormat()
	if err != nil {
		return err
	}

	l.logger.Info("downloading hostnames and IP block lists")
	blacklistSettings, err := settings.DoT.Blacklist.ToBlacklistFormat()
	if err != nil {
		return err
	}

	blockedHostnames, blockedIPs, blockedIPPrefixes, errs :=
		l.blockBuilder.All(ctx, blacklistSettings)
	for _, err := range errs {
		l.logger.Warn(err.Error())
	}

	// TODO change to BlockHostnames() when migrating to qdm12/dns v2
	unboundSettings.Blacklist.FqdnHostnames = blockedHostnames
	unboundSettings.Blacklist.IPs = blockedIPs
	unboundSettings.Blacklist.IPPrefixes = blockedIPPrefixes

	return l.conf.MakeUnboundConf(unboundSettings)
}
