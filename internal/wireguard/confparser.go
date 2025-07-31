package wireguard

import (
	"bufio"
	"os"
	"strings"
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
