package openvpn

import (
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/loopstate"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/openvpn/state"
	"github.com/qdm12/gluetun/internal/portforward"
	"github.com/qdm12/gluetun/internal/routing"
	"github.com/qdm12/golibs/logging"
)

var _ Looper = (*Loop)(nil)

type Looper interface {
	Runner
	loopstate.Getter
	loopstate.Applier
	SettingsGetSetter
	ServersGetterSetter
}

type Loop struct {
	statusManager loopstate.Manager
	state         state.Manager
	// Fixed parameters
	username       string
	puid           int
	pgid           int
	targetConfPath string
	// Configurators
	conf        StarterAuthWriter
	fw          firewallConfigurer
	routing     routing.VPNLocalGatewayIPGetter
	portForward portforward.StartStopper
	// Other objects
	logger      logging.Logger
	client      *http.Client
	tunnelReady chan<- struct{}
	// Internal channels and values
	stop        <-chan struct{}
	stopped     chan<- struct{}
	start       <-chan struct{}
	running     chan<- models.LoopStatus
	userTrigger bool
	startPFCh   chan struct{}
	// Internal constant values
	backoffTime time.Duration
}

type firewallConfigurer interface {
	firewall.VPNConnectionSetter
	firewall.PortAllower
}

const (
	defaultBackoffTime = 15 * time.Second
)

func NewLoop(settings configuration.OpenVPN, username string,
	puid, pgid int, allServers models.AllServers, conf Configurator,
	fw firewallConfigurer, routing routing.VPNLocalGatewayIPGetter,
	portForward portforward.StartStopper, logger logging.Logger,
	client *http.Client, tunnelReady chan<- struct{}) *Loop {
	start := make(chan struct{})
	running := make(chan models.LoopStatus)
	stop := make(chan struct{})
	stopped := make(chan struct{})

	statusManager := loopstate.New(constants.Stopped, start, running, stop, stopped)
	state := state.New(statusManager, settings, allServers)

	return &Loop{
		statusManager:  statusManager,
		state:          state,
		username:       username,
		puid:           puid,
		pgid:           pgid,
		targetConfPath: constants.OpenVPNConf,
		conf:           conf,
		fw:             fw,
		routing:        routing,
		portForward:    portForward,
		logger:         logger,
		client:         client,
		tunnelReady:    tunnelReady,
		start:          start,
		running:        running,
		stop:           stop,
		stopped:        stopped,
		userTrigger:    true,
		startPFCh:      make(chan struct{}),
		backoffTime:    defaultBackoffTime,
	}
}
