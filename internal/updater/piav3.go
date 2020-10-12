package updater

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"

	"github.com/qdm12/gluetun/internal/models"
)

func (u *updater) updatePIAOld(ctx context.Context) (err error) {
	const zipURL = "https://www.privateinternetaccess.com/openvpn/openvpn.zip"
	contents, err := fetchAndExtractFiles(zipURL)
	if err != nil {
		return err
	}
	const maxGoroutines = 10
	guard := make(chan struct{}, maxGoroutines)
	errors := make(chan error)
	serversCh := make(chan models.PIAOldServer)
	servers := make([]models.PIAOldServer, 0, len(contents))
	ctx, cancel := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}
	defer func() {
		cancel()
		wg.Wait()
		defer close(guard)
		defer close(errors)
		defer close(serversCh)
	}()
	for fileName, content := range contents {
		remoteLines := extractRemoteLinesFromOpenvpn(content)
		if len(remoteLines) == 0 {
			return fmt.Errorf("cannot find any remote lines in %s", fileName)
		}
		hosts := extractHostnamesFromRemoteLines(remoteLines)
		if len(hosts) == 0 {
			return fmt.Errorf("cannot find any hosts in %s", fileName)
		}
		region := strings.TrimSuffix(fileName, ".ovpn")
		wg.Add(1)
		go resolvePIAv3Hostname(ctx, wg, region, hosts, u.lookupIP, errors, serversCh, guard)
	}
	for range contents {
		select {
		case err := <-errors:
			return err
		case server := <-serversCh:
			servers = append(servers, server)
		}
	}
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Region < servers[j].Region
	})
	if u.options.Stdout {
		u.println(stringifyPIAOldServers(servers))
	}
	u.servers.PiaOld.Timestamp = u.timeNow().Unix()
	u.servers.PiaOld.Servers = servers
	return nil
}

func resolvePIAv3Hostname(ctx context.Context, wg *sync.WaitGroup,
	region string, hosts []string, lookupIP lookupIPFunc,
	errors chan<- error, serversCh chan<- models.PIAOldServer, guard chan struct{}) {
	guard <- struct{}{}
	defer func() {
		<-guard
		wg.Done()
	}()
	var IPs []net.IP //nolint:prealloc
	// usually one single host in this case
	// so no need to run in goroutines the for loop below
	for _, host := range hosts {
		const repetition = 5
		newIPs, err := resolveRepeat(ctx, lookupIP, host, repetition)
		if err != nil {
			errors <- err
			return
		}
		IPs = append(IPs, newIPs...)
	}
	serversCh <- models.PIAOldServer{
		Region: region,
		IPs:    uniqueSortedIPs(IPs),
	}
}

func stringifyPIAOldServers(servers []models.PIAOldServer) (s string) {
	s = "func PIAOldServers() []models.PIAOldServer {\n"
	s += "	return []models.PIAOldServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
