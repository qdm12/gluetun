package openvpn

import (
	"context"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/gluetun/internal/routing"
	"github.com/qdm12/golibs/logging"
)

type Looper interface {
	Run(ctx context.Context, done chan<- struct{})
	GetStatus() (status models.LoopStatus)
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
	GetSettings() (settings configuration.OpenVPN)
	SetSettings(ctx context.Context, settings configuration.OpenVPN) (
		outcome string)
	GetServers() (servers models.AllServers)
	SetServers(servers models.AllServers)
	GetPortForwarded() (port uint16)
	PortForward(vpnGatewayIP net.IP)
}

type looper struct {
	state *state
	// Fixed parameters
	username       string
	puid           int
	pgid           int
	targetConfPath string
	// Configurators
	conf    Configurator
	fw      firewall.Configurator
	routing routing.Routing
	// Other objects
	logger, pfLogger logging.Logger
	client           *http.Client
	tunnelReady      chan<- struct{}
	// Internal channels and values
	stop               <-chan struct{}
	stopped            chan<- struct{}
	start              <-chan struct{}
	running            chan<- models.LoopStatus
	portForwardSignals chan net.IP
	userTrigger        bool
	// Internal constant values
	backoffTime time.Duration
}

const (
	defaultBackoffTime = 15 * time.Second
)

func NewLooper(settings configuration.OpenVPN,
	username string, puid, pgid int, allServers models.AllServers,
	conf Configurator, fw firewall.Configurator, routing routing.Routing,
	logger logging.ParentLogger, client *http.Client,
	tunnelReady chan<- struct{}) Looper {
	start := make(chan struct{})
	running := make(chan models.LoopStatus)
	stop := make(chan struct{})
	stopped := make(chan struct{})

	state := newState(constants.Stopped, settings, allServers,
		start, running, stop, stopped)

	return &looper{
		state:              state,
		username:           username,
		puid:               puid,
		pgid:               pgid,
		targetConfPath:     constants.OpenVPNConf,
		conf:               conf,
		fw:                 fw,
		routing:            routing,
		logger:             logger.NewChild(logging.Settings{Prefix: "openvpn: "}),
		pfLogger:           logger.NewChild(logging.Settings{Prefix: "port forwarding: "}),
		client:             client,
		tunnelReady:        tunnelReady,
		start:              start,
		running:            running,
		stop:               stop,
		stopped:            stopped,
		portForwardSignals: make(chan net.IP),
		userTrigger:        true,
		backoffTime:        defaultBackoffTime,
	}
}

func (l *looper) PortForward(vpnGateway net.IP) { l.portForwardSignals <- vpnGateway }

func (l *looper) signalOrSetStatus(status models.LoopStatus) {
	if l.userTrigger {
		l.userTrigger = false
		select {
		case l.running <- status:
		default: // receiver calling ApplyStatus droppped out
		}
	} else {
		l.state.SetStatus(status)
	}
}

func (l *looper) Run(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	select {
	case <-l.start:
	case <-ctx.Done():
		return
	}

	for ctx.Err() == nil {
		settings, allServers := l.state.GetSettingsAndServers()

		providerConf := provider.New(settings.Provider.Name, allServers, time.Now)

		var connection models.OpenVPNConnection
		var lines []string
		var err error
		if settings.Config == "" {
			connection, err = providerConf.GetOpenVPNConnection(settings.Provider.ServerSelection)
			if err != nil {
				l.signalOrSetStatus(constants.Crashed)
				l.logAndWait(ctx, err)
				continue
			}
			lines = providerConf.BuildConf(connection, l.username, settings)
		} else {
			lines, connection, err = l.processCustomConfig(settings)
			if err != nil {
				l.signalOrSetStatus(constants.Crashed)
				l.logAndWait(ctx, err)
				continue
			}
		}

		if err := l.writeOpenvpnConf(lines); err != nil {
			l.signalOrSetStatus(constants.Crashed)
			l.logAndWait(ctx, err)
			continue
		}

		if settings.User != "" {
			err := l.conf.WriteAuthFile(
				settings.User, settings.Password, l.puid, l.pgid)
			if err != nil {
				l.signalOrSetStatus(constants.Crashed)
				l.logAndWait(ctx, err)
				continue
			}
		}

		if err := l.fw.SetVPNConnection(ctx, connection); err != nil {
			l.signalOrSetStatus(constants.Crashed)
			l.logAndWait(ctx, err)
			continue
		}

		openvpnCtx, openvpnCancel := context.WithCancel(context.Background())

		stdoutLines, stderrLines, waitError, err := l.conf.Start(
			openvpnCtx, settings.Version, settings.Flags)
		if err != nil {
			openvpnCancel()
			l.signalOrSetStatus(constants.Crashed)
			l.logAndWait(ctx, err)
			continue
		}

		lineCollectionDone := make(chan struct{})
		go l.collectLines(stdoutLines, stderrLines, lineCollectionDone)
		closeStreams := func() {
			close(stdoutLines)
			close(stderrLines)
			<-lineCollectionDone
		}

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

		l.backoffTime = defaultBackoffTime
		l.signalOrSetStatus(constants.Running)

		stayHere := true
		for stayHere {
			select {
			case <-ctx.Done():
				openvpnCancel()
				<-waitError
				close(waitError)
				closeStreams()
				<-portForwardDone
				return
			case <-l.stop:
				l.userTrigger = true
				l.logger.Info("stopping")
				openvpnCancel()
				<-waitError
				// do not close waitError or the waitError
				// select case will trigger
				closeStreams()
				<-portForwardDone
				l.stopped <- struct{}{}
			case <-l.start:
				l.userTrigger = true
				l.logger.Info("starting")
				stayHere = false
			case err := <-waitError: // unexpected error
				close(waitError)
				closeStreams()

				l.state.Lock() // prevent SetStatus from running in parallel

				openvpnCancel()
				l.state.SetStatus(constants.Crashed)
				<-portForwardDone
				l.logAndWait(ctx, err)
				stayHere = false

				l.state.Unlock()
			}
		}
		openvpnCancel()
	}
}

func (l *looper) logAndWait(ctx context.Context, err error) {
	if err != nil {
		l.logger.Error(err.Error())
	}
	l.logger.Info("retrying in " + l.backoffTime.String())
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
	providerConf.PortForward(ctx, client, l.pfLogger,
		gateway, l.fw, syncState)
}

func (l *looper) writeOpenvpnConf(lines []string) error {
	file, err := os.OpenFile(l.targetConfPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(strings.Join(lines, "\n"))
	if err != nil {
		return err
	}
	return file.Close()
}

func (l *looper) GetStatus() (status models.LoopStatus) {
	return l.state.GetStatus()
}
func (l *looper) ApplyStatus(ctx context.Context, status models.LoopStatus) (
	outcome string, err error) {
	return l.state.ApplyStatus(ctx, status)
}
func (l *looper) GetSettings() (settings configuration.OpenVPN) {
	return l.state.GetSettings()
}
func (l *looper) SetSettings(ctx context.Context, settings configuration.OpenVPN) (
	outcome string) {
	return l.state.SetSettings(ctx, settings)
}
func (l *looper) GetServers() (servers models.AllServers) {
	return l.state.GetServers()
}
func (l *looper) SetServers(servers models.AllServers) {
	l.state.SetServers(servers)
}
func (l *looper) GetPortForwarded() (port uint16) {
	return l.state.GetPortForwarded()
}
