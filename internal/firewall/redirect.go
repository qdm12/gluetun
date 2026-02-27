package firewall

import (
	"context"
	"fmt"
)

// RedirectPort redirects a source port to a destination port on the interface
// intf. If intf is empty, it is set to "*" which means all interfaces.
// If a redirection for the source port given already exists, it is removed first.
// If the destination port is zero, the redirection for the source port is removed
// and no new redirection is added.
func (c *Config) RedirectPort(ctx context.Context, intf string, sourcePort,
	destinationPort uint16,
) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if sourcePort == 0 {
		panic("source port cannot be 0")
	}

	newRedirection := portRedirection{
		interfaceName:   intf,
		sourcePort:      sourcePort,
		destinationPort: destinationPort,
	}

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating redirected ports internal state")
		if destinationPort == 0 {
			c.portRedirections.remove(intf, sourcePort)
			return nil
		}
		exists, conflict := c.portRedirections.check(newRedirection)
		switch {
		case exists:
			return nil
		case conflict != nil:
			c.portRedirections.remove(conflict.interfaceName,
				conflict.sourcePort)
		}
		c.portRedirections.append(newRedirection)
		return nil
	}

	exists, conflict := c.portRedirections.check(newRedirection)
	switch {
	case exists:
		return nil
	case conflict != nil:
		const remove = true
		err = c.impl.RedirectPort(ctx, conflict.interfaceName, conflict.sourcePort,
			conflict.destinationPort, remove)
		if err != nil {
			return fmt.Errorf("removing conflicting redirection: %w", err)
		}
		c.portRedirections.remove(conflict.interfaceName,
			conflict.sourcePort)
	}

	const remove = false
	err = c.impl.RedirectPort(ctx, intf, sourcePort, destinationPort, remove)
	if err != nil {
		return fmt.Errorf("redirecting port: %w", err)
	}
	c.portRedirections.append(newRedirection)

	return nil
}

type portRedirection struct {
	interfaceName   string
	sourcePort      uint16
	destinationPort uint16
}

type portRedirections []portRedirection

func (p *portRedirections) remove(intf string, sourcePort uint16) {
	slice := *p
	for i, redirection := range slice {
		interfaceMatch := intf == "" || intf == redirection.interfaceName
		if redirection.sourcePort == sourcePort && interfaceMatch {
			// Remove redirection - note: order does not matter
			slice[i] = slice[len(slice)-1]
			slice = slice[:len(slice)-1]
		}
	}
	*p = slice
}

func (p *portRedirections) check(dryRun portRedirection) (alreadyExists bool,
	conflict *portRedirection,
) {
	slice := *p
	for _, redirection := range slice {
		interfaceMatch := redirection.interfaceName == "" ||
			redirection.interfaceName == dryRun.interfaceName

		if redirection.sourcePort == dryRun.sourcePort &&
			redirection.destinationPort == dryRun.destinationPort &&
			interfaceMatch {
			return true, nil
		}

		if redirection.sourcePort == dryRun.sourcePort &&
			interfaceMatch {
			// Source port has a redirection already for the same interface or all interfaces
			return false, &redirection
		}
	}
	return false, nil
}

// append should be called after running `check` to avoid rule conflicts.
func (p *portRedirections) append(newRedirection portRedirection) {
	slice := *p
	slice = append(slice, newRedirection)
	*p = slice
}
