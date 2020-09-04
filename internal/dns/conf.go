package dns

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
)

func (c *configurator) MakeUnboundConf(settings settings.DNS, uid, gid int) (err error) {
	c.logger.Info("generating Unbound configuration")
	lines, warnings := generateUnboundConf(settings, c.client, c.logger)
	for _, warning := range warnings {
		c.logger.Warn(warning)
	}
	return c.fileManager.WriteLinesToFile(
		string(constants.UnboundConf),
		lines,
		files.Ownership(uid, gid),
		files.Permissions(0400))
}

// MakeUnboundConf generates an Unbound configuration from the user provided settings
func generateUnboundConf(settings settings.DNS, client network.Client, logger logging.Logger) (lines []string, warnings []error) {
	doIPv6 := "no"
	if settings.IPv6 {
		doIPv6 = "yes"
	}
	serverSection := map[string]string{
		// Logging
		"verbosity":     fmt.Sprintf("%d", settings.VerbosityLevel),
		"val-log-level": fmt.Sprintf("%d", settings.ValidationLogLevel),
		"use-syslog":    "no",
		// Performance
		"num-threads":       "1",
		"prefetch":          "yes",
		"prefetch-key":      "yes",
		"key-cache-size":    "16m",
		"key-cache-slabs":   "4",
		"msg-cache-size":    "4m",
		"msg-cache-slabs":   "4",
		"rrset-cache-size":  "4m",
		"rrset-cache-slabs": "4",
		"cache-min-ttl":     "3600",
		"cache-max-ttl":     "9000",
		// Privacy
		"rrset-roundrobin": "yes",
		"hide-identity":    "yes",
		"hide-version":     "yes",
		// Security
		"tls-cert-bundle":       fmt.Sprintf("%q", constants.CACertificates),
		"root-hints":            fmt.Sprintf("%q", constants.RootHints),
		"trust-anchor-file":     fmt.Sprintf("%q", constants.RootKey),
		"harden-below-nxdomain": "yes",
		"harden-referral-path":  "yes",
		"harden-algo-downgrade": "yes",
		// Network
		"do-ip4":    "yes",
		"do-ip6":    doIPv6,
		"interface": "127.0.0.1",
		"port":      "53",
		// Other
		"username": "\"nonrootuser\"",
	}

	// Block lists
	hostnamesLines, ipsLines, warnings := buildBlocked(client,
		settings.BlockMalicious, settings.BlockAds, settings.BlockSurveillance,
		settings.AllowedHostnames, settings.PrivateAddresses,
	)
	logger.Info("%d hostnames blocked overall", len(hostnamesLines))
	logger.Info("%d IP addresses blocked overall", len(ipsLines))
	sort.Slice(hostnamesLines, func(i, j int) bool { // for unit tests really
		return hostnamesLines[i] < hostnamesLines[j]
	})
	sort.Slice(ipsLines, func(i, j int) bool { // for unit tests really
		return ipsLines[i] < ipsLines[j]
	})

	// Server
	lines = append(lines, "server:")
	serverLines := make([]string, len(serverSection))
	i := 0
	for k, v := range serverSection {
		serverLines[i] = "  " + k + ": " + v
		i++
	}
	sort.Slice(serverLines, func(i, j int) bool {
		return serverLines[i] < serverLines[j]
	})
	lines = append(lines, serverLines...)
	lines = append(lines, hostnamesLines...)
	lines = append(lines, ipsLines...)

	// Forward zone
	lines = append(lines, "forward-zone:")
	forwardZoneSection := map[string]string{
		"name":                 "\".\"",
		"forward-tls-upstream": "yes",
	}
	if settings.Caching {
		forwardZoneSection["forward-no-cache"] = "no"
	} else {
		forwardZoneSection["forward-no-cache"] = "yes"
	}
	forwardZoneLines := make([]string, len(forwardZoneSection))
	i = 0
	for k, v := range forwardZoneSection {
		forwardZoneLines[i] = "  " + k + ": " + v
		i++
	}
	sort.Slice(forwardZoneLines, func(i, j int) bool {
		return forwardZoneLines[i] < forwardZoneLines[j]
	})
	for _, provider := range settings.Providers {
		providerData := constants.DNSProviderMapping()[provider]
		for _, IP := range providerData.IPs {
			forwardZoneLines = append(forwardZoneLines,
				fmt.Sprintf("  forward-addr: %s@853#%s", IP, providerData.Host))
		}
	}
	lines = append(lines, forwardZoneLines...)
	return lines, warnings
}

func buildBlocked(client network.Client, blockMalicious, blockAds, blockSurveillance bool,
	allowedHostnames, privateAddresses []string) (hostnamesLines, ipsLines []string, errs []error) {
	chHostnames := make(chan []string)
	chIPs := make(chan []string)
	chErrors := make(chan []error)
	go func() {
		lines, errs := buildBlockedHostnames(client, blockMalicious, blockAds, blockSurveillance, allowedHostnames)
		chHostnames <- lines
		chErrors <- errs
	}()
	go func() {
		lines, errs := buildBlockedIPs(client, blockMalicious, blockAds, blockSurveillance, privateAddresses)
		chIPs <- lines
		chErrors <- errs
	}()
	n := 2
	for n > 0 {
		select {
		case lines := <-chHostnames:
			hostnamesLines = append(hostnamesLines, lines...)
		case lines := <-chIPs:
			ipsLines = append(ipsLines, lines...)
		case routineErrs := <-chErrors:
			errs = append(errs, routineErrs...)
			n--
		}
	}
	return hostnamesLines, ipsLines, errs
}

func getList(client network.Client, url string) (results []string, err error) {
	content, status, err := client.GetContent(url)
	if err != nil {
		return nil, err
	} else if status != http.StatusOK {
		return nil, fmt.Errorf("HTTP status code is %d and not 200", status)
	}
	results = strings.Split(string(content), "\n")

	// remove empty lines
	last := len(results) - 1
	for i := range results {
		if len(results[i]) == 0 {
			results[i] = results[last]
			last--
		}
	}
	results = results[:last+1]

	if len(results) == 0 {
		return nil, nil
	}
	return results, nil
}

func buildBlockedHostnames(client network.Client, blockMalicious, blockAds, blockSurveillance bool,
	allowedHostnames []string) (lines []string, errs []error) {
	chResults := make(chan []string)
	chError := make(chan error)
	listsLeftToFetch := 0
	if blockMalicious {
		listsLeftToFetch++
		go func() {
			results, err := getList(client, string(constants.MaliciousBlockListHostnamesURL))
			chResults <- results
			chError <- err
		}()
	}
	if blockAds {
		listsLeftToFetch++
		go func() {
			results, err := getList(client, string(constants.AdsBlockListHostnamesURL))
			chResults <- results
			chError <- err
		}()
	}
	if blockSurveillance {
		listsLeftToFetch++
		go func() {
			results, err := getList(client, string(constants.SurveillanceBlockListHostnamesURL))
			chResults <- results
			chError <- err
		}()
	}
	uniqueResults := make(map[string]struct{})
	for listsLeftToFetch > 0 {
		select {
		case results := <-chResults:
			for _, result := range results {
				uniqueResults[result] = struct{}{}
			}
		case err := <-chError:
			listsLeftToFetch--
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	for _, allowedHostname := range allowedHostnames {
		delete(uniqueResults, allowedHostname)
	}
	for result := range uniqueResults {
		lines = append(lines, "  local-zone: \""+result+"\" static")
	}
	return lines, errs
}

func buildBlockedIPs(client network.Client, blockMalicious, blockAds, blockSurveillance bool,
	privateAddresses []string) (lines []string, errs []error) {
	chResults := make(chan []string)
	chError := make(chan error)
	listsLeftToFetch := 0
	if blockMalicious {
		listsLeftToFetch++
		go func() {
			results, err := getList(client, string(constants.MaliciousBlockListIPsURL))
			chResults <- results
			chError <- err
		}()
	}
	if blockAds {
		listsLeftToFetch++
		go func() {
			results, err := getList(client, string(constants.AdsBlockListIPsURL))
			chResults <- results
			chError <- err
		}()
	}
	if blockSurveillance {
		listsLeftToFetch++
		go func() {
			results, err := getList(client, string(constants.SurveillanceBlockListIPsURL))
			chResults <- results
			chError <- err
		}()
	}
	uniqueResults := make(map[string]struct{})
	for listsLeftToFetch > 0 {
		select {
		case results := <-chResults:
			for _, result := range results {
				uniqueResults[result] = struct{}{}
			}
		case err := <-chError:
			listsLeftToFetch--
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	for _, privateAddress := range privateAddresses {
		uniqueResults[privateAddress] = struct{}{}
	}
	for result := range uniqueResults {
		lines = append(lines, "  private-address: "+result)
	}
	return lines, errs
}
