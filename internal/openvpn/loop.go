package openvpn

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
)

type Looper interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	Restart()
	PortForward()
	GetSettings() (settings settings.OpenVPN)
	SetSettings(settings settings.OpenVPN)
	GetPortForwarded() (portForwarded uint16)
}

type looper struct {
	// Variable parameters
	provider           models.VPNProvider
	settings           settings.OpenVPN
	settingsMutex      sync.RWMutex
	portForwarded      uint16
	portForwardedMutex sync.RWMutex
	// Fixed parameters
	uid        int
	gid        int
	allServers models.AllServers
	// Configurators
	conf Configurator
	fw   firewall.Configurator
	// Other objects
	logger       logging.Logger
	client       network.Client
	fileManager  files.FileManager
	streamMerger command.StreamMerger
	cancel       context.CancelFunc
	// Internal channels
	restart            chan struct{}
	portForwardSignals chan struct{}
}

func NewLooper(provider models.VPNProvider, settings settings.OpenVPN,
	uid, gid int, allServers models.AllServers,
	conf Configurator, fw firewall.Configurator,
	logger logging.Logger, client network.Client, fileManager files.FileManager,
	streamMerger command.StreamMerger, cancel context.CancelFunc) Looper {
	return &looper{
		provider:           provider,
		settings:           settings,
		uid:                uid,
		gid:                gid,
		allServers:         allServers,
		conf:               conf,
		fw:                 fw,
		logger:             logger.WithPrefix("openvpn: "),
		client:             client,
		fileManager:        fileManager,
		streamMerger:       streamMerger,
		cancel:             cancel,
		restart:            make(chan struct{}),
		portForwardSignals: make(chan struct{}),
	}
}

func (l *looper) Restart()     { l.restart <- struct{}{} }
func (l *looper) PortForward() { l.portForwardSignals <- struct{}{} }

func (l *looper) GetSettings() (settings settings.OpenVPN) {
	l.settingsMutex.RLock()
	defer l.settingsMutex.RUnlock()
	return l.settings
}

func (l *looper) SetSettings(settings settings.OpenVPN) {
	l.settingsMutex.Lock()
	defer l.settingsMutex.Unlock()
	l.settings = settings
}

func (l *looper) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	select {
	case <-l.restart:
	case <-ctx.Done():
		return
	}
	defer l.logger.Warn("loop exited")

	for ctx.Err() == nil {
		settings := l.GetSettings()
		providerConf := provider.New(l.provider, l.allServers)
		connections, err := providerConf.GetOpenVPNConnections(settings.Provider.ServerSelection)
		if err != nil {
			l.logger.Error(err)
			l.cancel()
			return
		}
		lines := providerConf.BuildConf(
			connections,
			settings.Verbosity,
			l.uid,
			l.gid,
			settings.Root,
			settings.Cipher,
			settings.Auth,
			settings.Provider.ExtraConfigOptions,
		)
		if err := l.fileManager.WriteLinesToFile(string(constants.OpenVPNConf), lines, files.Ownership(l.uid, l.gid), files.Permissions(0400)); err != nil {
			l.logger.Error(err)
			l.cancel()
			return
		}

		if err := l.conf.WriteAuthFile(settings.User, settings.Password, l.uid, l.gid); err != nil {
			l.logger.Error(err)
			l.cancel()
			return
		}

		if err := l.fw.SetVPNConnections(ctx, connections); err != nil {
			l.logger.Error(err)
			l.cancel()
			return
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
				case <-l.portForwardSignals:
					l.portForward(ctx, providerConf, l.client)
				}
			}
		}(openvpnCtx)

		go l.streamMerger.Merge(openvpnCtx, stream, command.MergeName("openvpn"))
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
		case <-l.restart: // triggered restart
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

func (l *looper) portForward(ctx context.Context, providerConf provider.Provider, client network.Client) {
	settings := l.GetSettings()
	if !settings.Provider.PortForwarding.Enabled {
		return
	}
	var port uint16
	err := fmt.Errorf("")
	for err != nil {
		if ctx.Err() != nil {
			return
		}
		port, err = providerConf.GetPortForward(client)
		if err != nil {
			l.logAndWait(ctx, err)
		}
	}

	l.logger.Info("port forwarded is %d", port)
	l.portForwardedMutex.Lock()
	if err := l.fw.RemoveAllowedPort(ctx, l.portForwarded); err != nil {
		l.logger.Error(err)
	}
	if err := l.fw.SetAllowedPort(ctx, port, string(constants.TUN)); err != nil {
		l.logger.Error(err)
	}
	l.portForwarded = port
	l.portForwardedMutex.Unlock()

	filepath := settings.Provider.PortForwarding.Filepath
	l.logger.Info("writing forwarded port to %s", filepath)
	err = l.fileManager.WriteLinesToFile(
		string(filepath), []string{fmt.Sprintf("%d", port)},
		files.Ownership(l.uid, l.gid), files.Permissions(0400),
	)
	if err != nil {
		l.logger.Error(err)
	}
}

func (l *looper) GetPortForwarded() (portForwarded uint16) {
	l.portForwardedMutex.RLock()
	defer l.portForwardedMutex.RUnlock()
	return l.portForwarded
}
