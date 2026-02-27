package iptables

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/mod"
)

type kernelModules struct {
	nfConntrack  kernelModule
	nfRejectIPv4 kernelModule
	xtConnmark   kernelModule
	xtConntrack  kernelModule
	xtReject     kernelModule
}

type kernelModule struct {
	name string
	ok   bool
}

func newKernelModules() kernelModules {
	var m kernelModules
	nameToFieldPtr := map[string]*kernelModule{
		"nf_conntrack_netlink": &m.nfConntrack,
		"nf_reject_ipv4":       &m.nfRejectIPv4,
		"xt_connmark":          &m.xtConnmark,
		"xt_conntrack":         &m.xtConntrack,
		"xt_REJECT":            &m.xtReject,
	}
	for name, fieldPtr := range nameToFieldPtr {
		fieldPtr.name = name
		err := mod.Probe(name)
		fieldPtr.ok = err == nil
	}
	return m
}

func checkKernelModulesAreOK(modules ...kernelModule) error {
	missing := make([]string, 0, len(modules))
	for _, module := range modules {
		if !module.ok {
			missing = append(missing, module.name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("%w: %s", ErrKernelModuleMissing, strings.Join(missing, ", "))
	}
	return nil
}
