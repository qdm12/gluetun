package dns

import (
	"context"
	"fmt"

	"github.com/qdm12/dns/v2/pkg/blockbuilder"
	"github.com/qdm12/dns/v2/pkg/filter/update"
)

func (l *Loop) updateFiles(ctx context.Context) (err error) {
	settings := l.GetSettings()

	l.logger.Info("downloading hostnames and IP block lists")
	blacklistSettings := settings.DoT.Blacklist.ToBlockBuilderSettings()

	blockBuilder := blockbuilder.New(blacklistSettings)
	result := blockBuilder.BuildAll(ctx)
	for _, resultErr := range result.Errors {
		if err != nil {
			err = fmt.Errorf("%w, %w", err, resultErr)
			continue
		}
		err = resultErr
	}

	if err != nil {
		return err
	}

	updateSettings := update.Settings{
		IPs:        result.BlockedIPs,
		IPPrefixes: result.BlockedIPPrefixes,
	}
	updateSettings.BlockHostnames(result.BlockedHostnames)
	l.filter.Update(updateSettings)

	return nil
}
