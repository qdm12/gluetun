package dns

import "context"

func (l *looper) updateFiles(ctx context.Context) (err error) {
	l.logger.Info("downloading DNS over TLS cryptographic files")
	if err := l.conf.SetupFiles(ctx); err != nil {
		return err
	}
	settings := l.GetSettings()

	l.logger.Info("downloading hostnames and IP block lists")
	blockedHostnames, blockedIPs, blockedIPPrefixes, errs := l.blockBuilder.All(
		ctx, settings.BlacklistBuild)
	for _, err := range errs {
		l.logger.Warn(err.Error())
	}

	// TODO change to BlockHostnames() when migrating to qdm12/dns v2
	settings.Unbound.Blacklist.FqdnHostnames = blockedHostnames
	settings.Unbound.Blacklist.IPs = blockedIPs
	settings.Unbound.Blacklist.IPPrefixes = blockedIPPrefixes

	return l.conf.MakeUnboundConf(settings.Unbound)
}
