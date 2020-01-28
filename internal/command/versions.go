package command

import (
	"fmt"
	"strings"
)

// VersionOpenVPN obtains the version of the installed OpenVPN
func (c *command) VersionOpenVPN() (string, error) {
	output, err := c.command.Run("openvpn", "--version")
	if err != nil {
		return "", err
	}
	firstLine := strings.Split(output, "\n")[0]
	words := strings.Split(firstLine, " ")
	if len(words) < 2 {
		return "", fmt.Errorf("openvpn --version: first line is too short: %q", firstLine)
	}
	return words[1], nil
}

// VersionUnbound obtains the version of the installed Unbound
func (c *command) VersionUnbound() (string, error) {
	output, err := c.command.Run("unbound", "-h")
	if err != nil {
		return "", err
	}
	var version string
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "Version ") {
			words := strings.Split(line, " ")
			if len(words) < 2 {
				continue
			}
			version = words[1]
		}
	}
	if version == "" {
		return "", fmt.Errorf("unbound -h: version was not found in %q", output)
	}
	return version, nil
}

// VersionIptables obtains the version of the installed iptables
func (c *command) VersionIptables() (string, error) {
	output, err := c.command.Run("iptables", "--version")
	if err != nil {
		return "", err
	}
	words := strings.Split(output, " ")
	if len(words) < 2 {
		return "", fmt.Errorf("iptables --version: output is too short: %q", output)
	}
	return words[1], nil
}

// VersionShadowSocks obtains the version of the installed shadowsocks server
func (c *command) VersionShadowSocks() (string, error) {
	output, err := c.command.Run("ss-server", "-h")
	if err != nil {
		return "", err
	}
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("ss-server -h: not enough lines in %q", output)
	}
	words := strings.Split(lines[1], " ")
	if len(words) < 2 {
		return "", fmt.Errorf("ss-server -h: line 2 is too short: %q", lines[1])
	}
	return words[1], nil
}

// VersionTinyProxy obtains the version of the installed shadowsocks server
func (c *command) VersionTinyProxy() (string, error) {
	output, err := c.command.Run("tinyproxy", "-v")
	if err != nil {
		return "", err
	}
	words := strings.Split(output, " ")
	if len(words) < 2 {
		return "", fmt.Errorf("tinyproxy -v: output is too short: %q", output)
	}
	return words[1], nil
}
