package firewall

import (
	"context"
	"net/netip"
	"strings"
	"sync"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/routing"
)

type Config struct { //nolint:maligned
	runner         CmdRunner
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
	portRedirections  portRedirections
	appliedPostRules  []string
	stateMutex        sync.Mutex
}

// NewConfig creates a new Config instance and returns an error
// if no iptables implementation is available.
func NewConfig(ctx context.Context, logger Logger,
	runner CmdRunner, defaultRoutes []routing.DefaultRoute,
	localNetworks []routing.LocalNetwork,
) (config *Config, err error) {
	iptables, err := checkIptablesSupport(ctx, runner, "iptables", "iptables-nft", "iptables-legacy")
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

// clearAppliedPostRules removes all previously applied post-rules
func (c *Config) clearAppliedPostRules(ctx context.Context) error {
	for _, rule := range c.appliedPostRules {
		flippedRule := flipRule(rule)
		if strings.Contains(rule, "ip6tables") {
			if err := c.runIP6tablesInstruction(ctx, flippedRule); err != nil {
				c.logger.Debug("failed to remove post-rule (may not exist): " + err.Error())
			}
		} else {
			if err := c.runIptablesInstruction(ctx, flippedRule); err != nil {
				c.logger.Debug("failed to remove post-rule (may not exist): " + err.Error())
			}
		}
	}
	c.appliedPostRules = nil
	return nil
}

// applyPostRulesOnce applies post-rules only if they haven't been applied yet
func (c *Config) applyPostRulesOnce(ctx context.Context) error {
	if len(c.appliedPostRules) > 0 {
		c.logger.Debug("post-rules already applied, skipping")
		return nil
	}
	return c.runUserPostRules(ctx, c.customRulesPath, false)
}
