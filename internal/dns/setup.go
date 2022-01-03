package dns

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/qdm12/dns/pkg/check"
	"github.com/qdm12/dns/pkg/nameserver"
)

var errUpdateFiles = errors.New("cannot update files")

// Returning cancel == nil signals we want to re-run setupUnbound
// Returning err == errUpdateFiles signals we should not fall back
// on the plaintext DNS as DOT is still up and running.
func (l *Loop) setupUnbound(ctx context.Context) (
	cancel context.CancelFunc, waitError chan error, closeStreams func(), err error) {
	err = l.updateFiles(ctx)
	if err != nil {
		return nil, nil, nil,
			fmt.Errorf("%w: %s", errUpdateFiles, err)
	}

	settings := l.GetSettings()

	unboundCtx, cancel := context.WithCancel(context.Background())
	stdoutLines, stderrLines, waitError, err := l.conf.Start(unboundCtx,
		*settings.DoT.Unbound.VerbosityDetailsLevel)
	if err != nil {
		cancel()
		return nil, nil, nil, err
	}

	linesCollectionCtx, linesCollectionCancel := context.WithCancel(context.Background())
	lineCollectionDone := make(chan struct{})
	go l.collectLines(linesCollectionCtx, lineCollectionDone,
		stdoutLines, stderrLines)
	closeStreams = func() {
		linesCollectionCancel()
		<-lineCollectionDone
	}

	// use Unbound
	nameserver.UseDNSInternally(settings.ServerAddress)
	err = nameserver.UseDNSSystemWide(l.resolvConf, settings.ServerAddress,
		*settings.KeepNameserver)
	if err != nil {
		l.logger.Error(err.Error())
	}

	if err := check.WaitForDNS(ctx, net.DefaultResolver); err != nil {
		cancel()
		<-waitError
		close(waitError)
		closeStreams()
		return nil, nil, nil, err
	}

	return cancel, waitError, closeStreams, nil
}
