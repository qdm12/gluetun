package mod

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	errModuleNameUnknown     = errors.New("unknown module name")
	errKernelFeatureIsModule = errors.New("kernel feature is a module, not built-in")
	errKernelFeatureNotSet   = errors.New("kernel feature not set")
	errKernelFeatureNotFound = errors.New("kernel feature not found")
)

// checkProcConfig checks /proc/config.gz for a the kernel feature corresponding
// to the given module name. If the kernel feature is found and set to "y", it returns nil.
// If the kernel feature is found and set to "m", it returns an error indicating that the kernel
// feature is a module, not built-in.
// If the kernel feature is found and not set, it returns an error indicating that the kernel
// feature is not set. If the kernel feature is not found, it returns an error indicating that the kernel
// feature is not found.
func checkProcConfig(moduleName string) error {
	f, err := os.Open("/proc/config.gz")
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("creating gzip reader: %w", err)
	}
	defer gz.Close()

	// If any group of kernel features is satisfied, then the module is considered supported.
	kernelFeatureGroups, ok := moduleNameToKernelFeatureGroups(moduleName)
	if !ok {
		return fmt.Errorf("%w: %s", errModuleNameUnknown, moduleName)
	}
	groups := make([]map[string]bool, len(kernelFeatureGroups))
	for i, group := range kernelFeatureGroups {
		featureToOK := make(map[string]bool)
		for _, feature := range group {
			featureToOK[feature] = false
		}
		groups[i] = featureToOK
	}

	scanner := bufio.NewScanner(gz)
	for scanner.Scan() {
		line := scanner.Text()
		for _, featureToOK := range groups {
			for name, ok := range featureToOK {
				switch {
				case ok:
				case strings.HasPrefix(line, name+"=m"):
					return fmt.Errorf("%w: %s", errKernelFeatureIsModule, name)
				case strings.HasPrefix(line, name+"=y"):
					featureToOK[name] = true
					if allFeaturesOK(featureToOK) {
						return nil
					}
				case strings.HasPrefix(line, "# "+name+" is not set"):
					return fmt.Errorf("%w: %s", errKernelFeatureNotSet, name)
				}
			}
		}
	}

	return fmt.Errorf("%w: for module name %s", errKernelFeatureNotFound, moduleName)
}

func moduleNameToKernelFeatureGroups(moduleName string) (featureGroups [][]string, ok bool) {
	moduleMap := map[string][][]string{
		"nf_tables": {{"CONFIG_NF_TABLES"}},

		// Netfilter Matches
		"xt_conntrack": {{"CONFIG_NETFILTER_XT_MATCH_CONNTRACK"}},
		"xt_connmark": {
			{"CONFIG_NETFILTER_XT_CONNMARK"},
			{"CONFIG_NETFILTER_XT_MATCH_CONNMARK", "CONFIG_NETFILTER_XT_TARGET_CONNMARK"},
		},
		"xt_mark": {
			{"CONFIG_NETFILTER_XT_MARK"},
			{"CONFIG_NETFILTER_XT_MATCH_MARK", "CONFIG_NETFILTER_XT_TARGET_MARK"},
		},
		"nf_conntrack_netlink": {{"CONFIG_NF_CT_NETLINK"}},
		"nf_reject_ipv4":       {{"CONFIG_NF_REJECT_IPV4"}},

		// Common Netfilter Targets
		"xt_log": {{"CONFIG_NETFILTER_XT_TARGET_LOG"}},
		"xt_reject": {
			{"CONFIG_IP_NF_TARGET_REJECT", "CONFIG_NF_REJECT_IPV4"},
			{"CONFIG_NETFILTER_XT_TARGET_REJECT", "CONFIG_NF_REJECT_IPV4"},
		},
		"xt_masquerade": {{"CONFIG_NETFILTER_XT_TARGET_MASQUERADE"}},

		// Additional Netfilter Matches
		"xt_addrtype":  {{"CONFIG_NETFILTER_XT_MATCH_ADDRTYPE"}},
		"xt_comment":   {{"CONFIG_NETFILTER_XT_MATCH_COMMENT"}},
		"xt_multiport": {{"CONFIG_NETFILTER_XT_MATCH_MULTIPORT"}},
		"xt_state":     {{"CONFIG_NETFILTER_XT_MATCH_STATE"}},
		"xt_tcpudp":    {{"CONFIG_NETFILTER_XT_MATCH_TCPUDP"}},

		// Tunneling and Virtualization
		"tun":       {{"CONFIG_TUN"}},
		"bridge":    {{"CONFIG_BRIDGE"}},
		"veth":      {{"CONFIG_VETH"}},
		"vxlan":     {{"CONFIG_VXLAN"}},
		"wireguard": {{"CONFIG_WIREGUARD"}},

		// Filesystems
		"overlay": {{"CONFIG_OVERLAY_FS"}},
		"fuse":    {{"CONFIG_FUSE_FS"}},
	}

	featureGroups, ok = moduleMap[strings.ToLower(moduleName)]
	return featureGroups, ok
}

func allFeaturesOK(featureToOK map[string]bool) bool {
	for _, ok := range featureToOK {
		if !ok {
			return false
		}
	}
	return true
}
