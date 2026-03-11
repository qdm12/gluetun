package iptables

import (
	"context"
	"sync"

	"github.com/qdm12/gluetun/internal/mod"
)

type Config struct {
	runner         CmdRunner
	logger         Logger
	iptablesMutex  sync.Mutex
	ip6tablesMutex sync.Mutex

	// Fixed state
	ipTables  string
	ip6Tables string
	nftables  bool
	xtMark    bool
}

func New(ctx context.Context, runner CmdRunner, logger Logger) (*Config, error) {
	iptables, err := checkIptablesSupport(ctx, runner, "iptables", "iptables-nft", "iptables-legacy")
	if err != nil {
		return nil, err
	}

	ip6tables, err := findIP6tablesSupported(ctx, runner)
	if err != nil {
		return nil, err
	}

	modules := map[string]bool{
		"xt_mark":   false,
		"nf_tables": false,
	}
	for module := range modules {
		err := mod.Probe(module)
		modules[module] = err == nil
	}

	return &Config{
		runner:    runner,
		logger:    logger,
		ipTables:  iptables,
		ip6Tables: ip6tables,
		nftables:  modules["nf_tables"],
		xtMark:    modules["xt_mark"],
	}, nil
}
