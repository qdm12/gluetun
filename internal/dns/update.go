package dns

import (
	"context"
	"fmt"

	"github.com/qdm12/dns/v2/pkg/blockbuilder"
	"github.com/qdm12/dns/v2/pkg/middlewares/filter/update"
)

func (l *Loop) updateFiles(ctx context.Context) (err error) {
	settings := l.GetSettings()

	l.logger.Info("downloading hostnames and IP block lists")
	blacklistSettings := settings.DoT.Blacklist.ToBlockBuilderSettings(l.client)

	blockBuilder, err := blockbuilder.New(blacklistSettings)
	if err != nil {
		return fmt.Errorf("creating block builder: %w", err)
	}

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
	err = l.filter.Update(updateSettings)
	if err != nil {
		return fmt.Errorf("updating filter: %w", err)
	}

	return nil
}
