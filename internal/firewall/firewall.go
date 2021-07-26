// Package firewall defines a configurator used to change the state
// of the firewall as well as do some light routing changes.
package firewall

import (
	"context"
	"net"
	"sync"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/routing"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
)

var _ Configurator = (*Config)(nil)

// Configurator allows to change firewall rules and modify network routes.
type Configurator interface {
	Enabler
	VPNConnectionSetter
	PortAllower
	OutboundSubnetsSetter
}

type Config struct { //nolint:maligned
	runner           command.Runner
	logger           logging.Logger
	iptablesMutex    sync.Mutex
	ip6tablesMutex   sync.Mutex
	defaultInterface string
	defaultGateway   net.IP
	localNetworks    []routing.LocalNetwork
	localIP          net.IP

	// Fixed state
	ip6Tables       bool
	customRulesPath string

	// State
	enabled           bool
	vpnConnection     models.OpenVPNConnection
	outboundSubnets   []net.IPNet
	allowedInputPorts map[uint16]string // port to interface mapping
	stateMutex        sync.Mutex
}

// NewConfig creates a new Config instance.
func NewConfig(logger logging.Logger, runner command.Runner,
	defaultInterface string, defaultGateway net.IP,
	localNetworks []routing.LocalNetwork, localIP net.IP) *Config {
	return &Config{
		runner:            runner,
		logger:            logger,
		allowedInputPorts: make(map[uint16]string),
		ip6Tables:         ip6tablesSupported(context.Background(), runner),
		customRulesPath:   "/iptables/post-rules.txt",
		// Obtained from routing
		defaultInterface: defaultInterface,
		defaultGateway:   defaultGateway,
		localNetworks:    localNetworks,
		localIP:          localIP,
	}
}
