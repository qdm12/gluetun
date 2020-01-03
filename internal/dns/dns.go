package dns

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/settings"
)

// MakeUnboundConf generates an Unbound configuration from the user provided settings
func MakeUnboundConf(settings *settings.DNS) (conf string) {
	serverSection := map[string][]string{
		// Logging
		"verbosity":     []string{fmt.Sprintf("%d", settings.Verbosity)},
		"val-log-level": []string{fmt.Sprintf("%d", settings.LogLevel)},
		"use-syslog":    []string{"no"},
		// Performance
		"num-threads":       []string{"1"},
		"prefetch":          []string{"yes"},
		"prefetch-key":      []string{"yes"},
		"key-cache-size":    []string{"16m"},
		"key-cache-slabs":   []string{"4"},
		"msg-cache-size":    []string{"4m"},
		"msg-cache-slabs":   []string{"4"},
		"rrset-cache-size":  []string{"4m"},
		"rrset-cache-slabs": []string{"4"},
		"cache-min-ttl":     []string{"3600"},
		"cache-max-ttl":     []string{"9000"},
		// Privacy
		"rrset-roundrobin": []string{"yes"},
		"hide-identity":    []string{"yes"},
		"hide-version":     []string{"yes"},
		// Security
		"tls-cert-bundle":       []string{"\"/etc/ssl/certs/ca-certificates.crt\""},
		"root-hints":            []string{"\"/etc/unbound/root.hints\""},
		"trust-anchor-file":     []string{"\"/etc/unbound/root.key\""},
		"harden-below-nxdomain": []string{"yes"},
		"harden-referral-path":  []string{"yes"},
		"harden-algo-downgrade": []string{"yes"},
		// Network
		"do-ip4":    []string{"yes"},
		"do-ip6":    []string{"no"},
		"interface": []string{"127.0.0.1"},
		"port":      []string{"53"},
		// Other
		"username": []string{"\"nonrootuser\""},
	}
	serverSection["private-address"] = settings.PrivateAddresses
	forwardZoneSection := map[string][]string{
		"name":                 []string{"\".\""},
		"forward-tls-upstream": []string{"yes"},
	}
	forwardZoneSection["forward-addr"] = settings.Provider.GetForwardAddresses()

	// Block lists
	blockConf := buildBlockedHostnames(settings.BlockMalicious, settings.BlockAds, settings.BlockSurveillance)

	// Make configuration string
	conf = "server:\n"
	for k, arr := range serverSection {
		for i := range arr {
			conf += "  " + k + ": " + arr[i] + "\n"
		}
	}
	conf += blockConf
	conf += "forward-zone:\n"
	for k, arr := range forwardZoneSection {
		for i := range arr {
			conf += "  " + k + ": " + arr[i] + "\n"
		}
	}
	return conf
}

func getList(client *http.Client, URL string, chResults chan []string, chError chan error) {
	content, err := network.GetContent(client, URL)
	if err != nil {
		chError <- err
		return
	}
	chResults <- strings.Split(string(content), "\n")
	chError <- nil
}

func buildBlockedHostnames(blockMalicious, blockAds, blockSurveillance bool) (conf string) {
	client := &http.Client{Timeout: 5 * time.Second}
	chHostnames := make(chan []string)
	chIPs := make(chan []string)
	chError := make(chan error)
	listsLeftToFetch := 0
	if blockMalicious {
		listsLeftToFetch += 2
		go getList(client, constants.MaliciousBlockListHostnamesURL, chHostnames, chError)
		go getList(client, constants.MaliciousBlockListIPsURL, chIPs, chError)
	}
	if blockAds {
		listsLeftToFetch += 2
		go getList(client, constants.AdsBlockListHostnamesURL, chHostnames, chError)
		go getList(client, constants.AdsBlockListIPsURL, chIPs, chError)
	}
	if blockSurveillance {
		listsLeftToFetch += 2
		go getList(client, constants.SurveillanceBlockListHostnamesURL, chHostnames, chError)
		go getList(client, constants.SurveillanceBlockListIPsURL, chIPs, chError)
	}
	uniqueHostnames := make(map[string]struct{})
	uniqueIPs := make(map[string]struct{})
	for listsLeftToFetch > 0 {
		select {
		case hostnames := <-chHostnames:
			for _, hostname := range hostnames {
				uniqueHostnames[hostname] = struct{}{}
			}
		case IPs := <-chIPs:
			for _, IP := range IPs {
				uniqueIPs[IP] = struct{}{}
			}
		case err := <-chError:
			listsLeftToFetch--
			if err != nil {
				logging.Warn(err.Error())
			}
		}
	}
	for hostname := range uniqueHostnames {
		conf += "local-zone: \"" + hostname + "\" static\n"
	}
	for IP := range uniqueIPs {
		conf += "private-address: " + IP + "\n"
	}
	return conf
}
