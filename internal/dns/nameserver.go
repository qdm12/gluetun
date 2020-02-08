package dns

import (
	"context"
	"net"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// UseDNSInternally is to change the Go program DNS only
func (c *configurator) UseDNSInternally(IP net.IP) {
	c.logger.Info("%s: using DNS address %s internally", logPrefix, IP.String())
	net.DefaultResolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", net.JoinHostPort(IP.String(), "53"))
		},
	}
}

// UseDNSSystemWide changes the nameserver to use for DNS system wide
func (c *configurator) UseDNSSystemWide(IP net.IP) error {
	c.logger.Info("%s: using DNS address %s system wide", logPrefix, IP.String())
	data, err := c.fileManager.ReadFile(string(constants.ResolvConf))
	if err != nil {
		return err
	}
	s := strings.TrimSuffix(string(data), "\n")
	lines := strings.Split(s, "\n")
	if len(lines) == 1 && lines[0] == "" {
		lines = nil
	}
	found := false
	for i := range lines {
		if strings.HasPrefix(lines[i], "nameserver ") {
			lines[i] = "nameserver " + IP.String()
			found = true
		}
	}
	if !found {
		lines = append(lines, "nameserver "+IP.String())
	}
	data = []byte(strings.Join(lines, "\n"))
	return c.fileManager.WriteToFile(string(constants.ResolvConf), data)
}
