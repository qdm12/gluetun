package dns

import (
	"context"
	"io/ioutil"
	"net"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/os"
)

// UseDNSInternally is to change the Go program DNS only.
func (c *configurator) UseDNSInternally(ip net.IP) {
	c.logger.Info("using DNS address %s internally", ip.String())
	net.DefaultResolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", net.JoinHostPort(ip.String(), "53"))
		},
	}
}

// UseDNSSystemWide changes the nameserver to use for DNS system wide.
func (c *configurator) UseDNSSystemWide(ip net.IP, keepNameserver bool) error {
	c.logger.Info("using DNS address %s system wide", ip.String())
	const filepath = string(constants.ResolvConf)
	file, err := c.openFile(filepath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}

	s := strings.TrimSuffix(string(data), "\n")

	lines := []string{
		"nameserver " + ip.String(),
	}
	for _, line := range strings.Split(s, "\n") {
		if line == "" ||
			(!keepNameserver && strings.HasPrefix(line, "nameserver ")) {
			continue
		}
		lines = append(lines, line)
	}

	s = strings.Join(lines, "\n") + "\n"

	file, err = c.openFile(filepath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(s)
	if err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}
