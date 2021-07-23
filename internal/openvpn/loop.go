package openvpn

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/loopstate"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/openvpn/state"
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
	statusManager loopstate.Manager
	state         state.Manager
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

	statusManager := loopstate.New(constants.Stopped, start, running, stop, stopped)
	state := state.New(statusManager, settings, allServers)

	return &looper{
		statusManager:      statusManager,
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
