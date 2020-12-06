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
	GetStatus() (status models.LoopStatus)
	SetStatus(status models.LoopStatus) (outcome string, err error)
	GetSettings() (settings settings.OpenVPN)
	SetSettings(settings settings.OpenVPN) (outcome string)
	GetServers() (servers models.AllServers)
	SetServers(servers models.AllServers)
	GetPortForwarded() (port uint16)
	PortForward(vpnGatewayIP net.IP)
}

type looper struct {
	state state
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
	// Internal channels and locks
	loopLock           sync.Mutex
	running            chan models.LoopStatus
	stop, stopped      chan struct{}
	start              chan struct{}
	portForwardSignals chan net.IP
}

func NewLooper(settings settings.OpenVPN,
	uid, gid int, allServers models.AllServers,
	conf Configurator, fw firewall.Configurator, routing routing.Routing,
	logger logging.Logger, client *http.Client, fileManager files.FileManager,
	streamMerger command.StreamMerger, cancel context.CancelFunc) Looper {
	return &looper{
		state: state{
			status:     constants.Stopped,
			settings:   settings,
			allServers: allServers,
		},
		uid:                uid,
		gid:                gid,
		conf:               conf,
		fw:                 fw,
		routing:            routing,
		logger:             logger.WithPrefix("openvpn: "),
		pfLogger:           logger.WithPrefix("port forwarding: "),
		client:             client,
		fileManager:        fileManager,
		streamMerger:       streamMerger,
		cancel:             cancel,
		start:              make(chan struct{}),
		running:            make(chan models.LoopStatus),
		stop:               make(chan struct{}),
		stopped:            make(chan struct{}),
		portForwardSignals: make(chan net.IP),
	}
}

func (l *looper) PortForward(vpnGateway net.IP) { l.portForwardSignals <- vpnGateway }

func (l *looper) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	crashed := false
	select {
	case <-l.start:
	case <-ctx.Done():
		return
	}
	defer l.logger.Warn("loop exited")

	for ctx.Err() == nil {
		settings, allServers := l.state.getSettingsAndServers()
		providerConf := provider.New(settings.Provider.Name, allServers, time.Now)
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
			if !crashed {
				l.running <- constants.Crashed
				crashed = true
			}
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

		if !crashed {
			l.running <- constants.Running
			crashed = false
		} else {
			l.state.setStatusWithLock(constants.Running)
		}

		stayHere := true
		for stayHere {
			select {
			case <-ctx.Done():
				l.logger.Warn("context canceled: exiting loop")
				openvpnCancel()
				<-waitError
				close(waitError)
				return
			case <-l.stop:
				l.logger.Info("stopping")
				openvpnCancel()
				<-waitError
				close(waitError)
				l.stopped <- struct{}{}
			case <-l.start:
				l.logger.Info("starting")
				stayHere = false
			case err := <-waitError: // unexpected error
				openvpnCancel()
				close(waitError)
				l.state.setStatusWithLock(constants.Crashed)
				l.logAndWait(ctx, err)
				crashed = true
				stayHere = false
			}
		}
		openvpnCancel() // just for the linter
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
	l.state.portForwardedMu.RLock()
	settings := l.state.settings
	l.state.portForwardedMu.RUnlock()
	if !settings.Provider.PortForwarding.Enabled {
		return
	}
	syncState := func(port uint16) (pfFilepath models.Filepath) {
		l.state.portForwardedMu.Lock()
		defer l.state.portForwardedMu.Unlock()
		l.state.portForwarded = port
		l.state.settingsMu.RLock()
		defer l.state.settingsMu.RUnlock()
		return settings.Provider.PortForwarding.Filepath
	}
	providerConf.PortForward(ctx,
		client, l.fileManager, l.pfLogger,
		gateway, l.fw, syncState)
}
