package wireguard

import (
	"bufio"
	"fmt"
	"net/netip"
	"os"
	"strings"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type ConfFile struct {
	Interface map[string]string
	Peer      map[string]string
}

// ParseConfFile parses a WireGuard .conf file into a ConfFile struct.
func ParseConfFile(path string) (*ConfFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	section := ""
	conf := &ConfFile{
		Interface: make(map[string]string),
		Peer:      make(map[string]string),
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.ToLower(strings.Trim(line, "[]"))
			continue
		}
		if section == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch section {
		case "interface":
			conf.Interface[key] = val
		case "peer":
			conf.Peer[key] = val
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return conf, nil
}

// DetectProvider tries to guess the VPN provider from a parsed WireGuard config.
func DetectProvider(conf *ConfFile) string {
	endpoint := conf.Peer["Endpoint"]
	endpoint = strings.ToLower(endpoint)
	// Check for ProtonVPN (IP or domain patterns, or comment markers)
	if strings.Contains(endpoint, "protonvpn") || strings.HasPrefix(endpoint, "185.") || strings.HasPrefix(endpoint, "45.") {
		return "protonvpn"
	}
	// Check for Mullvad
	if strings.Contains(endpoint, "mullvad.net") {
		return "mullvad"
	}
	// Add more providers as needed
	return "unknown"
}

// ValidateConfFile validates a parsed WireGuard config file for correctness.
func ValidateConfFile(conf *ConfFile) error {
	// Validate Interface section
	if conf.Interface == nil {
		return fmt.Errorf("missing [Interface] section")
	}

	// Validate PrivateKey
	privateKey, exists := conf.Interface["PrivateKey"]
	if !exists || privateKey == "" {
		return fmt.Errorf("missing or empty PrivateKey in [Interface] section")
	}
	if _, err := wgtypes.ParseKey(privateKey); err != nil {
		return fmt.Errorf("invalid PrivateKey: %w", err)
	}

	// Validate Address (optional but should be valid if present)
	if address, exists := conf.Interface["Address"]; exists && address != "" {
		addresses := strings.Split(address, ",")
		for _, addr := range addresses {
			addr = strings.TrimSpace(addr)
			if addr == "" {
				continue
			}
			if _, err := netip.ParsePrefix(addr); err != nil {
				return fmt.Errorf("invalid Address '%s': %w", addr, err)
			}
		}
	}

	// Validate Peer section
	if conf.Peer == nil {
		return fmt.Errorf("missing [Peer] section")
	}

	// Validate PublicKey
	publicKey, exists := conf.Peer["PublicKey"]
	if !exists || publicKey == "" {
		return fmt.Errorf("missing or empty PublicKey in [Peer] section")
	}
	if _, err := wgtypes.ParseKey(publicKey); err != nil {
		return fmt.Errorf("invalid PublicKey: %w", err)
	}

	// Validate PreSharedKey (optional but should be valid if present)
	if preSharedKey, exists := conf.Peer["PreSharedKey"]; exists && preSharedKey != "" {
		if _, err := wgtypes.ParseKey(preSharedKey); err != nil {
			return fmt.Errorf("invalid PreSharedKey: %w", err)
		}
	}

	// Validate Endpoint
	endpoint, exists := conf.Peer["Endpoint"]
	if !exists || endpoint == "" {
		return fmt.Errorf("missing or empty Endpoint in [Peer] section")
	}
	// Basic endpoint format validation (host:port)
	if !strings.Contains(endpoint, ":") {
		return fmt.Errorf("invalid Endpoint format '%s': must be in host:port format", endpoint)
	}

	// Validate AllowedIPs
	allowedIPs, exists := conf.Peer["AllowedIPs"]
	if !exists || allowedIPs == "" {
		return fmt.Errorf("missing or empty AllowedIPs in [Peer] section")
	}
	ips := strings.Split(allowedIPs, ",")
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		if ip == "" {
			continue
		}
		if _, err := netip.ParsePrefix(ip); err != nil {
			return fmt.Errorf("invalid AllowedIPs entry '%s': %w", ip, err)
		}
	}

	return nil
}
