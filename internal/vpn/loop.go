package vpn

import (
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/dns"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/loopstate"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/qdm12/gluetun/internal/openvpn"
	"github.com/qdm12/gluetun/internal/portforward"
	"github.com/qdm12/gluetun/internal/publicip"
	"github.com/qdm12/gluetun/internal/routing"
	"github.com/qdm12/gluetun/internal/vpn/state"
	"github.com/qdm12/golibs/command"
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
	buildInfo     models.BuildInformation
	versionInfo   bool
	vpnInputPorts []uint16 // TODO make changeable through stateful firewall
	// Configurators
	openvpnConf openvpn.Interface
	netLinker   netlink.NetLinker
	fw          firewallConfigurer
	routing     routing.VPNGetter
	portForward portforward.StartStopper
	publicip    publicip.Looper
	dnsLooper   dns.Looper
	// Other objects
	starter command.Starter // for OpenVPN
	logger  logging.ParentLogger
	client  *http.Client
	// Internal channels and values
	stop        <-chan struct{}
	stopped     chan<- struct{}
	start       <-chan struct{}
	running     chan<- models.LoopStatus
	userTrigger bool
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

func NewLoop(vpnSettings settings.VPN, vpnInputPorts []uint16,
	allServers models.AllServers, openvpnConf openvpn.Interface,
	netLinker netlink.NetLinker, fw firewallConfigurer, routing routing.VPNGetter,
	portForward portforward.StartStopper, starter command.Starter,
	publicip publicip.Looper, dnsLooper dns.Looper,
	logger logging.ParentLogger, client *http.Client,
	buildInfo models.BuildInformation, versionInfo bool) *Loop {
	start := make(chan struct{})
	running := make(chan models.LoopStatus)
	stop := make(chan struct{})
	stopped := make(chan struct{})

	statusManager := loopstate.New(constants.Stopped, start, running, stop, stopped)
	state := state.New(statusManager, vpnSettings, allServers)

	return &Loop{
		statusManager: statusManager,
		state:         state,
		buildInfo:     buildInfo,
		versionInfo:   versionInfo,
		vpnInputPorts: vpnInputPorts,
		openvpnConf:   openvpnConf,
		netLinker:     netLinker,
		fw:            fw,
		routing:       routing,
		portForward:   portForward,
		publicip:      publicip,
		dnsLooper:     dnsLooper,
		starter:       starter,
		logger:        logger,
		client:        client,
		start:         start,
		running:       running,
		stop:          stop,
		stopped:       stopped,
		userTrigger:   true,
		backoffTime:   defaultBackoffTime,
	}
}
