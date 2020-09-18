package updater

import (
	"bytes"
	"context"
	"net"
	"sort"
)

func newResolver(resolverAddress string) *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", net.JoinHostPort(resolverAddress, "53"))
		},
	}
}

func newLookupIP(r *net.Resolver) lookupIPFunc {
	return func(ctx context.Context, host string) (ips []net.IP, err error) {
		addresses, err := r.LookupIPAddr(ctx, host)
		if err != nil {
			return nil, err
		}
		ips = make([]net.IP, len(addresses))
		for i := range addresses {
			ips[i] = addresses[i].IP
		}
		return ips, nil
	}
}

func resolveRepeat(ctx context.Context, lookupIP lookupIPFunc, host string, n int) (ips []net.IP, err error) {
	foundIPs := make(chan []net.IP)
	errors := make(chan error)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for i := 0; i < n; i++ {
		go func() {
			newIPs, err := lookupIP(ctx, host)
			if err != nil {
				errors <- err
			} else {
				foundIPs <- newIPs
			}
		}()
	}

	uniqueIPs := make(map[string]struct{})
	for i := 0; i < n; i++ {
		select {
		case newIPs := <-foundIPs:
			for _, ip := range newIPs {
				key := ip.String()
				uniqueIPs[key] = struct{}{}
			}
		case newErr := <-errors:
			if err == nil {
				err = newErr
				cancel()
			}
		}
	}

	ips = make([]net.IP, 0, len(uniqueIPs))
	for key := range uniqueIPs {
		ip := net.ParseIP(key)
		if ipv4 := ip.To4(); ipv4 != nil {
			ip = ipv4
		}
		ips = append(ips, ip)
	}

	sort.Slice(ips, func(i, j int) bool {
		return bytes.Compare(ips[i], ips[j]) < 1
	})

	return ips, err
}
