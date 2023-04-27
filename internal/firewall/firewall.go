package firewall

import (
	"context"
	"net/netip"
	"sync"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/routing"
	"github.com/qdm12/golibs/command"
)

type Config struct { //nolint:maligned
	runner         command.Runner
	logger         Logger
	iptablesMutex  sync.Mutex
	ip6tablesMutex sync.Mutex
	defaultRoutes  []routing.DefaultRoute
	localNetworks  []routing.LocalNetwork

	// Fixed state
	ipTables        string
	ip6Tables       string
	customRulesPath string

	// State
	enabled           bool
	vpnConnection     models.Connection
	vpnIntf           string
	outboundSubnets   []netip.Prefix
	allowedInputPorts map[uint16]map[string]struct{} // port to interfaces set mapping
	stateMutex        sync.Mutex
}

// NewConfig creates a new Config instance and returns an error
// if no iptables implementation is available.
func NewConfig(ctx context.Context, logger Logger,
	runner command.Runner, defaultRoutes []routing.DefaultRoute,
	localNetworks []routing.LocalNetwork) (config *Config, err error) {
	iptables, err := checkIptablesSupport(ctx, runner, "iptables", "iptables-nft")
	if err != nil {
		return nil, err
	}

	ip6tables, err := findIP6tablesSupported(ctx, runner)
	if err != nil {
		return nil, err
	}

	return &Config{
		runner:            runner,
		logger:            logger,
		allowedInputPorts: make(map[uint16]map[string]struct{}),
		ipTables:          iptables,
		ip6Tables:         ip6tables,
		customRulesPath:   "/iptables/post-rules.txt",
		// Obtained from routing
		defaultRoutes: defaultRoutes,
		localNetworks: localNetworks,
	}, nil
}
