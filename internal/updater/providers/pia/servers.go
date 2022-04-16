// Package pia contains code to obtain the server information
// for the Private Internet Access provider.
package pia

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/models"
)

var ErrNotEnoughServers = errors.New("not enough servers found")

func GetServers(ctx context.Context, client *http.Client, minServers int) (
	servers []models.Server, err error) {
	nts := make(nameToServer)

	noChangeCounter := 0
	const maxNoChange = 10
	const betweenDuration = 200 * time.Millisecond
	const maxDuration = time.Minute

	maxTimer := time.NewTimer(maxDuration)

	for {
		data, err := fetchAPI(ctx, client)
		if err != nil {
			return nil, err
		}

		change := addData(data.Regions, nts)

		if !change {
			noChangeCounter++
			if noChangeCounter == maxNoChange {
				break
			}
		} else {
			noChangeCounter = 0
		}

		timer := time.NewTimer(betweenDuration)
		maxTimeout := false
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			if !maxTimer.Stop() {
				<-timer.C
			}
			return nil, ctx.Err()
		case <-timer.C:
		case <-maxTimer.C:
			if !timer.Stop() {
				<-timer.C
			}
			maxTimeout = true
		}

		if maxTimeout {
			break
		}
	}

	servers = nts.toServersSlice()

	if len(servers) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(servers), minServers)
	}

	sortServers(servers)

	return servers, nil
}

func addData(regions []regionData, nts nameToServer) (change bool) {
	for _, region := range regions {
		for _, server := range region.Servers.UDP {
			const tcp, udp = false, true
			if nts.add(server.CN, region.DNS, region.Name, tcp, udp, region.PortForward, server.IP) {
				change = true
			}
		}

		for _, server := range region.Servers.TCP {
			const tcp, udp = true, false
			if nts.add(server.CN, region.DNS, region.Name, tcp, udp, region.PortForward, server.IP) {
				change = true
			}
		}
	}

	return change
}
