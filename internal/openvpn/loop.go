package openvpn

import (
	"context"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/gluetun/internal/routing"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
)

type Looper interface {
	Run(ctx context.Context, done chan<- struct{})
	GetStatus() (status models.LoopStatus)
	SetStatus(status models.LoopStatus) (outcome string, err error)
	GetSettings() (settings configuration.OpenVPN)
	SetSettings(settings configuration.OpenVPN) (outcome string)
	GetServers() (servers models.AllServers)
	SetServers(servers models.AllServers)
	GetPortForwarded() (port uint16)
	PortForward(vpnGatewayIP net.IP)
}

type looper struct {
	state state
	// Fixed parameters
	username string
	puid     int
	pgid     int
	// Configurators
	conf    Configurator
	fw      firewall.Configurator
	routing routing.Routing
	// Other objects
	logger, pfLogger logging.Logger
	client           *http.Client
	openFile         os.OpenFileFunc
	tunnelReady      chan<- struct{}
	healthy          <-chan bool
	// Internal channels and locks
	loopLock           sync.Mutex
	running            chan models.LoopStatus
	stop, stopped      chan struct{}
	start              chan struct{}
	portForwardSignals chan net.IP
	crashed            bool
	backoffTime        time.Duration
	healthWaitTime     time.Duration
}

const (
	defaultBackoffTime    = 15 * time.Second
	defaultHealthWaitTime = 6 * time.Second
)

func NewLooper(settings configuration.OpenVPN,
	username string, puid, pgid int, allServers models.AllServers,
	conf Configurator, fw firewall.Configurator, routing routing.Routing,
	logger logging.Logger, client *http.Client, openFile os.OpenFileFunc,
	tunnelReady chan<- struct{}, healthy <-chan bool) Looper {
	return &looper{
		state: state{
			status:     constants.Stopped,
			settings:   settings,
			allServers: allServers,
		},
		username:           username,
		puid:               puid,
		pgid:               pgid,
		conf:               conf,
		fw:                 fw,
		routing:            routing,
		logger:             logger.NewChild(logging.SetPrefix("openvpn: ")),
		pfLogger:           logger.NewChild(logging.SetPrefix("port forwarding: ")),
		client:             client,
		openFile:           openFile,
		tunnelReady:        tunnelReady,
		healthy:            healthy,
		start:              make(chan struct{}),
		running:            make(chan models.LoopStatus),
		stop:               make(chan struct{}),
		stopped:            make(chan struct{}),
		portForwardSignals: make(chan net.IP),
		backoffTime:        defaultBackoffTime,
		healthWaitTime:     defaultHealthWaitTime,
	}
}

func (l *looper) PortForward(vpnGateway net.IP) { l.portForwardSignals <- vpnGateway }

func (l *looper) signalCrashedStatus() {
	if !l.crashed {
		l.crashed = true
		l.running <- constants.Crashed
	}
}

func (l *looper) Run(ctx context.Context, done chan<- struct{}) { //nolint:gocognit
	defer close(done)
	select {
	case <-l.start:
	case <-ctx.Done():
		return
	}

	for ctx.Err() == nil {
		settings, allServers := l.state.getSettingsAndServers()

		providerConf := provider.New(settings.Provider.Name, allServers, time.Now)

		var connection models.OpenVPNConnection
		var lines []string
		var err error
		if len(settings.Config) == 0 {
			connection, err = providerConf.GetOpenVPNConnection(settings.Provider.ServerSelection)
			if err != nil {
				l.signalCrashedStatus()
				l.logAndWait(ctx, err)
				continue
			}
			lines = providerConf.BuildConf(connection, l.username, settings)
		} else {
			lines, connection, err = l.processCustomConfig(settings)
			if err != nil {
				l.signalCrashedStatus()
				l.logAndWait(ctx, err)
				continue
			}
		}

		if err := writeOpenvpnConf(lines, l.openFile); err != nil {
			l.signalCrashedStatus()
			l.logAndWait(ctx, err)
			continue
		}

		if err := l.conf.WriteAuthFile(settings.User, settings.Password, l.puid, l.pgid); err != nil {
			l.signalCrashedStatus()
			l.logAndWait(ctx, err)
			continue
		}

		if err := l.fw.SetVPNConnection(ctx, connection); err != nil {
			l.signalCrashedStatus()
			l.logAndWait(ctx, err)
			continue
		}

		openvpnCtx, openvpnCancel := context.WithCancel(context.Background())

		stdoutLines, stderrLines, waitError, err := l.conf.Start(openvpnCtx)
		if err != nil {
			openvpnCancel()
			l.signalCrashedStatus()
			l.logAndWait(ctx, err)
			continue
		}

		lineCollectionDone := make(chan struct{})
		go l.collectLines(stdoutLines, stderrLines, lineCollectionDone)

		// Needs the stream line from main.go to know when the tunnel is up
		portForwardDone := make(chan struct{})
		go func(ctx context.Context) {
			defer close(portForwardDone)
			select {
			// TODO have a way to disable pf with a context
			case <-ctx.Done():
				return
			case gateway := <-l.portForwardSignals:
				l.portForward(ctx, providerConf, l.client, gateway)
			}
		}(openvpnCtx)

		if l.crashed {
			l.crashed = false
			l.backoffTime = defaultBackoffTime
			l.state.setStatusWithLock(constants.Running)
		} else {
			l.running <- constants.Running
		}

		stayHere := true
		for stayHere {
			select {
			case <-ctx.Done():
				openvpnCancel()
				<-waitError
				close(waitError)
				close(stdoutLines)
				close(stderrLines)
				<-lineCollectionDone
				<-portForwardDone
				return
			case <-l.stop:
				l.logger.Info("stopping")
				openvpnCancel()
				<-waitError
				l.stopped <- struct{}{}
			case <-l.start:
				l.logger.Info("starting")
				stayHere = false
			case err := <-waitError: // unexpected error
				openvpnCancel()
				l.state.setStatusWithLock(constants.Crashed)
				l.logAndWait(ctx, err)
				l.crashed = true
				stayHere = false
			case healthy := <-l.healthy:
				if healthy {
					continue
				}
				// ensure it stays unhealthy for some time before restarting it
				healthy = l.waitForHealth(ctx)
				if healthy || ctx.Err() != nil {
					continue
				}
				l.crashed = true // flag as crashed
				l.state.setStatusWithLock(constants.Stopping)
				l.logger.Warn("unhealthy program: restarting openvpn")
				openvpnCancel()
				<-waitError
				l.state.setStatusWithLock(constants.Stopped)
				stayHere = false
			}
		}
		close(waitError)
		close(stdoutLines)
		close(stderrLines)
		openvpnCancel() // just for the linter
	}
}

func (l *looper) logAndWait(ctx context.Context, err error) {
	if err != nil {
		l.logger.Error(err)
	}
	l.logger.Info("retrying in %s", l.backoffTime)
	timer := time.NewTimer(l.backoffTime)
	l.backoffTime *= 2
	select {
	case <-timer.C:
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
	}
}

// waitForHealth waits for a true healthy signal
// after restarting openvpn in order to avoid restarting
// openvpn in a loop as it requires a few seconds to connect.
func (l *looper) waitForHealth(ctx context.Context) (healthy bool) {
	l.logger.Info("unhealthy program: waiting %s for it to change to healthy", l.healthWaitTime)
	timer := time.NewTimer(l.healthWaitTime)
	l.healthWaitTime *= 2
	for {
		select {
		case healthy = <-l.healthy:
			if !healthy {
				break
			}
			if !timer.Stop() {
				<-timer.C
			}
			l.healthWaitTime = defaultHealthWaitTime
			return true
		case <-timer.C:
			return false
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return false
		}
	}
}

// portForward is a blocking operation which may or may not be infinite.
// You should therefore always call it in a goroutine.
func (l *looper) portForward(ctx context.Context,
	providerConf provider.Provider, client *http.Client, gateway net.IP) {
	l.state.portForwardedMu.RLock()
	settings := l.state.settings
	l.state.portForwardedMu.RUnlock()
	if !settings.Provider.PortForwarding.Enabled {
		return
	}
	syncState := func(port uint16) (pfFilepath string) {
		l.state.portForwardedMu.Lock()
		defer l.state.portForwardedMu.Unlock()
		l.state.portForwarded = port
		l.state.settingsMu.RLock()
		defer l.state.settingsMu.RUnlock()
		return settings.Provider.PortForwarding.Filepath
	}
	providerConf.PortForward(ctx,
		client, l.openFile, l.pfLogger,
		gateway, l.fw, syncState)
}

func writeOpenvpnConf(lines []string, openFile os.OpenFileFunc) error {
	file, err := openFile(constants.OpenVPNConf, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(strings.Join(lines, "\n"))
	if err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return nil
}
