package updater

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/resolver"
)

func (u *updater) updateHideMyAss(ctx context.Context) (err error) {
	servers, warnings, err := findHideMyAssServers(ctx, u.client, u.presolver)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("HideMyAss: %s", warning)
		}
	}
	if err != nil {
		return fmt.Errorf("%w: HideMyAss: %s", ErrUpdateServerInformation, err)
	}
	if u.options.Stdout {
		u.println(stringifyHideMyAssServers(servers))
	}
	u.servers.HideMyAss.Timestamp = u.timeNow().Unix()
	u.servers.HideMyAss.Servers = servers
	return nil
}

func findHideMyAssServers(ctx context.Context, client *http.Client, presolver resolver.Parallel) (
	servers []models.HideMyAssServer, warnings []string, err error) {
	TCPhostToURL, err := findHideMyAssHostToURLForProto(ctx, client, "TCP")
	if err != nil {
		return nil, nil, err
	}

	UDPhostToURL, err := findHideMyAssHostToURLForProto(ctx, client, "UDP")
	if err != nil {
		return nil, nil, err
	}

	uniqueHosts := make(map[string]struct{}, len(TCPhostToURL))
	for host := range TCPhostToURL {
		uniqueHosts[host] = struct{}{}
	}
	for host := range UDPhostToURL {
		uniqueHosts[host] = struct{}{}
	}

	hosts := make([]string, len(uniqueHosts))
	i := 0
	for host := range uniqueHosts {
		hosts[i] = host
		i++
	}

	const (
		maxFailRatio    = 0.1
		maxDuration     = 15 * time.Second
		betweenDuration = 2 * time.Second
		maxNoNew        = 2
		maxFails        = 2
	)
	settings := resolver.ParallelSettings{
		MaxFailRatio: maxFailRatio,
		Repeat: resolver.RepeatSettings{
			MaxDuration:     maxDuration,
			BetweenDuration: betweenDuration,
			MaxNoNew:        maxNoNew,
			MaxFails:        maxFails,
			SortIPs:         true,
		},
	}
	hostToIPs, warnings, err := presolver.Resolve(ctx, hosts, settings)
	if err != nil {
		return nil, warnings, err
	}

	servers = make([]models.HideMyAssServer, 0, len(hostToIPs))
	for host, IPs := range hostToIPs {
		tcpURL, tcp := TCPhostToURL[host]
		udpURL, udp := UDPhostToURL[host]

		var url, protocol string
		if tcp {
			url = tcpURL
			protocol = "TCP"
		} else if udp {
			url = udpURL
			protocol = "UDP"
		}
		country, region, city := parseHideMyAssURL(url, protocol)

		server := models.HideMyAssServer{
			Country:  country,
			Region:   region,
			City:     city,
			Hostname: host,
			IPs:      IPs,
			TCP:      tcp,
			UDP:      udp,
		}
		servers = append(servers, server)
	}

	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Country+servers[i].Region+servers[i].City+servers[i].Hostname <
			servers[j].Country+servers[j].Region+servers[j].City+servers[j].Hostname
	})

	return servers, warnings, nil
}

func findHideMyAssHostToURLForProto(ctx context.Context, client *http.Client, protocol string) (
	hostToURL map[string]string, err error) {
	indexURL := "https://vpn.hidemyass.com/vpn-config/" + strings.ToUpper(protocol) + "/"

	urls, err := fetchHideMyAssHTTPIndex(ctx, client, indexURL)
	if err != nil {
		return nil, err
	}

	return fetchMultiOvpnFiles(ctx, client, urls)
}

func parseHideMyAssURL(url, protocol string) (country, region, city string) {
	lastSlashIndex := strings.LastIndex(url, "/")
	url = url[lastSlashIndex+1:]

	suffix := "." + strings.ToUpper(protocol) + ".ovpn"
	url = strings.TrimSuffix(url, suffix)

	parts := strings.Split(url, ".")

	switch len(parts) {
	case 1:
		country = parts[0]
		return country, "", ""
	case 2: //nolint:gomnd
		country = parts[0]
		city = parts[1]
	default:
		country = parts[0]
		region = parts[1]
		city = parts[2]
	}

	return camelCaseToWords(country), camelCaseToWords(region), camelCaseToWords(city)
}

func camelCaseToWords(camelCase string) (words string) {
	wasLowerCase := false
	for _, r := range camelCase {
		if wasLowerCase && unicode.IsUpper(r) {
			words += " "
		}
		wasLowerCase = unicode.IsLower(r)
		words += string(r)
	}
	return words
}

var hideMyAssIndexRegex = regexp.MustCompile(`<a[ ]+href=".+\.ovpn">.+\.ovpn</a>`)

func fetchHideMyAssHTTPIndex(ctx context.Context, client *http.Client, indexURL string) (urls []string, err error) {
	htmlCode, err := fetchFile(ctx, client, indexURL)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(indexURL, "/") {
		indexURL += "/"
	}

	lines := strings.Split(string(htmlCode), "\n")
	for _, line := range lines {
		found := hideMyAssIndexRegex.FindString(line)
		if len(found) == 0 {
			continue
		}
		const prefix = `.ovpn">`
		const suffix = `</a>`
		startIndex := strings.Index(found, prefix) + len(prefix)
		endIndex := strings.Index(found, suffix)
		filename := found[startIndex:endIndex]
		url := indexURL + filename
		if !strings.HasSuffix(url, ".ovpn") {
			continue
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func fetchMultiOvpnFiles(ctx context.Context, client *http.Client, urls []string) (
	hostToURL map[string]string, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	hostToURL = make(map[string]string, len(urls))

	type Result struct {
		url  string
		host string
	}
	results := make(chan Result)
	errors := make(chan error)
	for _, url := range urls {
		go func(url string) {
			host, err := fetchOvpnFile(ctx, client, url)
			if err != nil {
				errors <- fmt.Errorf("%w: for %s", err, url)
				return
			}
			results <- Result{
				url:  url,
				host: host,
			}
		}(url)
	}

	for range urls {
		select {
		case newErr := <-errors:
			if err == nil { // only assign to the first error
				err = newErr
				cancel() // stop other operations, this will trigger other errors we ignore
			}
		case result := <-results:
			hostToURL[result.host] = result.url
		}
	}

	if err != nil {
		return nil, err
	}

	return hostToURL, nil
}

func fetchOvpnFile(ctx context.Context, client *http.Client, url string) (hostname string, err error) {
	b, err := fetchFile(ctx, client, url)
	if err != nil {
		return "", err
	}

	const rejectIP = true
	const rejectDomain = false
	hosts := extractRemoteHostsFromOpenvpn(b, rejectIP, rejectDomain)
	if len(hosts) == 0 {
		return "", errRemoteHostNotFound
	}

	return hosts[0], nil
}

func fetchFile(ctx context.Context, client *http.Client, url string) (b []byte, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

func stringifyHideMyAssServers(servers []models.HideMyAssServer) (s string) {
	s = "func HideMyAssServers() []models.HideMyAssServer {\n"
	s += "	return []models.HideMyAssServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
