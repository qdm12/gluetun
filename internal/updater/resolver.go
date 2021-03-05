package updater

import (
	"bytes"
	"context"
	"net"
	"sort"
	"time"
)

func newResolver(resolverAddress string) *net.Resolver {
	d := net.Dialer{}
	resolverAddress = net.JoinHostPort(resolverAddress, "53")
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return d.DialContext(ctx, "udp", resolverAddress)
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

func parallelResolve(ctx context.Context, lookupIP lookupIPFunc, hosts []string,
	repetition int, timeBetween time.Duration, failOnErr bool) (
	hostToIPs map[string][]net.IP, warnings []string, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	type result struct {
		host string
		ips  []net.IP
	}

	results := make(chan result)
	defer close(results)
	errors := make(chan error)
	defer close(errors)

	for _, host := range hosts {
		go func(host string) {
			ips, err := resolveRepeat(ctx, lookupIP, host, repetition, timeBetween)
			if err != nil {
				errors <- err
				return
			}
			results <- result{
				host: host,
				ips:  ips,
			}
		}(host)
	}

	hostToIPs = make(map[string][]net.IP, len(hosts))

	for range hosts {
		select {
		case newErr := <-errors:
			if !failOnErr {
				warnings = append(warnings, newErr.Error())
			} else if err == nil {
				err = newErr
				cancel()
			}
		case r := <-results:
			hostToIPs[r.host] = r.ips
		}
	}

	return hostToIPs, warnings, err
}

func resolveRepeat(ctx context.Context, lookupIP lookupIPFunc, host string,
	repetition int, timeBetween time.Duration) (ips []net.IP, err error) {
	uniqueIPs := make(map[string]struct{})

	i := 0
	for {
		newIPs, newErr := lookupIP(ctx, host)
		if err == nil {
			err = newErr // it's fine to fail some of the resolutions
		}
		for _, ip := range newIPs {
			key := ip.String()
			uniqueIPs[key] = struct{}{}
		}

		i++
		if i == repetition {
			break
		}

		timer := time.NewTimer(timeBetween)
		select {
		case <-timer.C:
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return nil, ctx.Err()
		}
	}

	if len(uniqueIPs) == 0 {
		return nil, err
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
