package dns

import (
	"fmt"
	"sort"
	"strings"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/settings"
)

func (c *configurator) MakeUnboundConf(settings settings.DNS) (err error) {
	lines, warnings, err := generateUnboundConf(settings, c.client, c.logger)
	for _, warning := range warnings {
		c.logger.Warn(warning)
	}
	if err != nil {
		return err
	}
	return c.fileManager.WriteLinesToFile(string(constants.UnboundConf), lines)
}

// MakeUnboundConf generates an Unbound configuration from the user provided settings
func generateUnboundConf(settings settings.DNS, client network.Client, logger logging.Logger) (lines []string, warnings []error, err error) {
	serverSection := map[string]string{
		// Logging
		"verbosity":     fmt.Sprintf("%d", settings.Verbosity),
		"val-log-level": fmt.Sprintf("%d", settings.LogLevel),
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
		"tls-cert-bundle":       "\"/etc/ssl/certs/ca-certificates.crt\"",
		"root-hints":            "\"/etc/unbound/root.hints\"",
		"trust-anchor-file":     "\"/etc/unbound/root.key\"",
		"harden-below-nxdomain": "yes",
		"harden-referral-path":  "yes",
		"harden-algo-downgrade": "yes",
		// Network
		"do-ip4":    "yes",
		"do-ip6":    "no",
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
	var serverLines []string
	for k, v := range serverSection {
		serverLines = append(serverLines, "  "+k+": "+v)
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
	var forwardZoneLines []string
	for k, v := range forwardZoneSection {
		forwardZoneLines = append(forwardZoneLines, "  "+k+": "+v)
	}
	sort.Slice(forwardZoneLines, func(i, j int) bool {
		return forwardZoneLines[i] < forwardZoneLines[j]
	})
	lines = append(lines, forwardZoneLines...)
	forwardAddresses, ok := constants.DNSAddressesMapping[settings.Provider]
	if !ok || len(forwardAddresses) == 0 {
		return nil, warnings, fmt.Errorf("DNS provider %q does not have any matching forward addresses", settings.Provider)
	}
	for _, forwardAddress := range forwardAddresses {
		lines = append(lines, fmt.Sprintf("  forward-addr: %s", forwardAddress))
	}
	return lines, warnings, nil
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

func getList(client network.Client, URL string) (results []string, err error) {
	content, status, err := client.GetContent(URL)
	if err != nil {
		return nil, err
	} else if status != 200 {
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
