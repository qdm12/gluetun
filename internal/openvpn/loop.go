package openvpn

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/firewall"
	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/qdm12/private-internet-access-docker/internal/provider"
	"github.com/qdm12/private-internet-access-docker/internal/settings"
)

type Looper interface {
	Run(ctx context.Context, restart, portForward <-chan struct{}, wg *sync.WaitGroup)
}

type looper struct {
	// Variable parameters
	provider models.VPNProvider
	settings settings.OpenVPN
	// Fixed parameters
	uid int
	gid int
	// Configurators
	conf Configurator
	fw   firewall.Configurator
	// Other objects
	logger       logging.Logger
	client       network.Client
	fileManager  files.FileManager
	streamMerger command.StreamMerger
	fatalOnError func(err error)
}

func NewLooper(provider models.VPNProvider, settings settings.OpenVPN,
	uid, gid int,
	conf Configurator, fw firewall.Configurator,
	logger logging.Logger, client network.Client, fileManager files.FileManager,
	streamMerger command.StreamMerger, fatalOnError func(err error)) Looper {
	return &looper{
		provider:     provider,
		settings:     settings,
		uid:          uid,
		gid:          gid,
		conf:         conf,
		fw:           fw,
		logger:       logger.WithPrefix("openvpn: "),
		client:       client,
		fileManager:  fileManager,
		streamMerger: streamMerger,
		fatalOnError: fatalOnError,
	}
}

func (l *looper) Run(ctx context.Context, restart, portForward <-chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	select {
	case <-restart:
	case <-ctx.Done():
		return
	}
	defer l.logger.Warn("loop exited")

	for ctx.Err() == nil {
		providerConf := provider.New(l.provider, l.client, l.fileManager)
		connections, err := providerConf.GetOpenVPNConnections(l.settings.Provider.ServerSelection)
		l.fatalOnError(err)
		err = providerConf.BuildConf(
			connections,
			l.settings.Verbosity,
			l.uid,
			l.gid,
			l.settings.Root,
			l.settings.Cipher,
			l.settings.Auth,
			l.settings.Provider.ExtraConfigOptions,
		)
		l.fatalOnError(err)

		err = l.conf.WriteAuthFile(l.settings.User, l.settings.Password, l.uid, l.gid)
		l.fatalOnError(err)

		if err := l.fw.SetVPNConnections(ctx, connections); err != nil {
			l.fatalOnError(err)
		}

		openvpnCtx, openvpnCancel := context.WithCancel(context.Background())

		stream, waitFn, err := l.conf.Start(openvpnCtx)
		if err != nil {
			openvpnCancel()
			l.logAndWait(ctx, err)
			continue
		}

		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-portForward:
					l.portForward(ctx, providerConf)
				}
			}
		}(openvpnCtx)

		go l.streamMerger.Merge(openvpnCtx, stream,
			command.MergeName("openvpn"), command.MergeColor(constants.ColorOpenvpn()))
		waitError := make(chan error)
		go func() {
			err := waitFn() // blocking
			waitError <- err
		}()
		select {
		case <-ctx.Done():
			l.logger.Warn("context canceled: exiting loop")
			openvpnCancel()
			<-waitError
			close(waitError)
			return
		case <-restart: // triggered restart
			l.logger.Info("restarting")
			openvpnCancel()
			<-waitError
			close(waitError)
		case err := <-waitError: // unexpected error
			openvpnCancel()
			close(waitError)
			l.logAndWait(ctx, err)
		}
	}
}

func (l *looper) logAndWait(ctx context.Context, err error) {
	l.logger.Error(err)
	l.logger.Info("retrying in 30 seconds")
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel() // just for the linter
	<-ctx.Done()
}

func (l *looper) portForward(ctx context.Context, providerConf provider.Provider) {
	if !l.settings.Provider.PortForwarding.Enabled {
		return
	}
	var port uint16
	err := fmt.Errorf("")
	for err != nil {
		if ctx.Err() != nil {
			return
		}
		port, err = providerConf.GetPortForward()
		if err != nil {
			l.logAndWait(ctx, err)
			continue
		}
		l.logger.Info("port forwarded is %d", port)
	}

	filepath := l.settings.Provider.PortForwarding.Filepath
	l.logger.Info("writing forwarded port to %s", filepath)
	err = l.fileManager.WriteLinesToFile(
		string(filepath), []string{fmt.Sprintf("%d", port)},
		files.Ownership(l.uid, l.gid), files.Permissions(0400),
	)
	if err != nil {
		l.logger.Error(err)
	}

	if err := l.fw.SetPortForward(ctx, port); err != nil {
		l.logger.Error(err)
	}
}
