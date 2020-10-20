package openvpn

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/gluetun/internal/routing"
	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
)

type Looper interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	Restart()
	PortForward(vpnGatewayIP net.IP)
	GetSettings() (settings settings.OpenVPN)
	SetSettings(settings settings.OpenVPN)
	GetPortForwarded() (portForwarded uint16)
	SetAllServers(allServers models.AllServers)
}

type looper struct {
	// Variable parameters
	provider           models.VPNProvider
	settings           settings.OpenVPN
	settingsMutex      sync.RWMutex
	portForwarded      uint16
	portForwardedMutex sync.RWMutex
	allServers         models.AllServers
	allServersMutex    sync.RWMutex
	// Fixed parameters
	uid int
	gid int
	// Configurators
	conf    Configurator
	fw      firewall.Configurator
	routing routing.Routing
	// Other objects
	logger, pfLogger logging.Logger
	client           *http.Client
	fileManager      files.FileManager
	streamMerger     command.StreamMerger
	cancel           context.CancelFunc
	// Internal channels
	restart            chan struct{}
	portForwardSignals chan net.IP
}

func NewLooper(provider models.VPNProvider, settings settings.OpenVPN,
	uid, gid int, allServers models.AllServers,
	conf Configurator, fw firewall.Configurator, routing routing.Routing,
	logger logging.Logger, client *http.Client, fileManager files.FileManager,
	streamMerger command.StreamMerger, cancel context.CancelFunc) Looper {
	return &looper{
		provider:           provider,
		settings:           settings,
		uid:                uid,
		gid:                gid,
		allServers:         allServers,
		conf:               conf,
		fw:                 fw,
		routing:            routing,
		logger:             logger.WithPrefix("openvpn: "),
		pfLogger:           logger.WithPrefix("port forwarding: "),
		client:             client,
		fileManager:        fileManager,
		streamMerger:       streamMerger,
		cancel:             cancel,
		restart:            make(chan struct{}),
		portForwardSignals: make(chan net.IP),
	}
}

func (l *looper) Restart()                      { l.restart <- struct{}{} }
func (l *looper) PortForward(vpnGateway net.IP) { l.portForwardSignals <- vpnGateway }

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

func (l *looper) SetAllServers(allServers models.AllServers) {
	l.allServersMutex.Lock()
	defer l.allServersMutex.Unlock()
	l.allServers = allServers
}

func (l *looper) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-l.restart:
	case <-ctx.Done():
		return
	}
	defer l.logger.Warn("loop exited")

	for ctx.Err() == nil {
		settings := l.GetSettings()
		l.allServersMutex.RLock()
		providerConf := provider.New(l.provider, l.allServers, time.Now)
		l.allServersMutex.RUnlock()
		connection, err := providerConf.GetOpenVPNConnection(settings.Provider.ServerSelection)
		if err != nil {
			l.logger.Error(err)
			l.cancel()
			return
		}
		lines := providerConf.BuildConf(
			connection,
			settings.Verbosity,
			l.uid,
			l.gid,
			settings.Root,
			settings.Cipher,
			settings.Auth,
			settings.Provider.ExtraConfigOptions,
		)
		if err := l.fileManager.WriteLinesToFile(string(constants.OpenVPNConf), lines,
			files.Ownership(l.uid, l.gid), files.Permissions(constants.UserReadPermission)); err != nil {
			l.logger.Error(err)
			l.cancel()
			return
		}

		if err := l.conf.WriteAuthFile(settings.User, settings.Password, l.uid, l.gid); err != nil {
			l.logger.Error(err)
			l.cancel()
			return
		}

		if err := l.fw.SetVPNConnection(ctx, connection); err != nil {
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

		// Needs the stream line from main.go to know when the tunnel is up
		go func(ctx context.Context) {
			for {
				select {
				// TODO have a way to disable pf with a context
				case <-ctx.Done():
					return
				case gateway := <-l.portForwardSignals:
					wg.Add(1)
					go l.portForward(ctx, wg, providerConf, l.client, gateway)
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
	const waitTime = 30 * time.Second
	l.logger.Info("retrying in %s", waitTime)
	timer := time.NewTimer(waitTime)
	select {
	case <-timer.C:
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
	}
}

// portForward is a blocking operation which may or may not be infinite.
// You should therefore always call it in a goroutine.
func (l *looper) portForward(ctx context.Context, wg *sync.WaitGroup,
	providerConf provider.Provider, client *http.Client, gateway net.IP) {
	defer wg.Done()
	settings := l.GetSettings()
	if !settings.Provider.PortForwarding.Enabled {
		return
	}
	syncState := func(port uint16) (pfFilepath models.Filepath) {
		l.portForwardedMutex.Lock()
		l.portForwarded = port
		l.portForwardedMutex.Unlock()
		settings := l.GetSettings()
		return settings.Provider.PortForwarding.Filepath
	}
	providerConf.PortForward(ctx,
		client, l.fileManager, l.pfLogger,
		gateway, l.fw, syncState)
}

func (l *looper) GetPortForwarded() (portForwarded uint16) {
	l.portForwardedMutex.RLock()
	defer l.portForwardedMutex.RUnlock()
	return l.portForwarded
}
