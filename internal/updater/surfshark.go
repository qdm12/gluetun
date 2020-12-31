package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/network"
)

func (u *updater) updateSurfshark(ctx context.Context) (err error) {
	servers, warnings, err := findSurfsharkServersFromZip(ctx, u.client, u.lookupIP)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("Surfshark: %s", warning)
		}
	}
	if err != nil {
		return fmt.Errorf("cannot update Surfshark servers: %w", err)
	}
	if u.options.Stdout {
		u.println(stringifySurfsharkServers(servers))
	}
	u.servers.Surfshark.Timestamp = u.timeNow().Unix()
	u.servers.Surfshark.Servers = servers
	return nil
}

//nolint:deadcode,unused
func findSurfsharkServersFromAPI(ctx context.Context, client network.Client, lookupIP lookupIPFunc) (
	servers []models.SurfsharkServer, warnings []string, err error) {
	const url = "https://my.surfshark.com/vpn/api/v1/server/clusters"
	b, status, err := client.Get(ctx, url)
	if err != nil {
		return nil, nil, err
	} else if status != http.StatusOK {
		return nil, nil, fmt.Errorf("HTTP status code %d", status)
	}
	var jsonServers []struct {
		Host     string `json:"connectionName"`
		Country  string `json:"country"`
		Location string `json:"location"`
	}
	if err := json.Unmarshal(b, &jsonServers); err != nil {
		return nil, nil, err
	}
	for _, jsonServer := range jsonServers {
		host := jsonServer.Host
		const repetition = 5
		IPs, err := resolveRepeat(ctx, lookupIP, host, repetition)
		if err != nil {
			return nil, warnings, err
		} else if len(IPs) == 0 {
			warning := fmt.Sprintf("no IP address found for host %q", host)
			warnings = append(warnings, warning)
			continue
		}
		server := models.SurfsharkServer{
			Region: jsonServer.Country + " " + jsonServer.Location,
			IPs:    uniqueSortedIPs(IPs),
		}
		servers = append(servers, server)
	}
	return servers, warnings, nil
}

func findSurfsharkServersFromZip(ctx context.Context, client network.Client, lookupIP lookupIPFunc) (
	servers []models.SurfsharkServer, warnings []string, err error) {
	const zipURL = "https://my.surfshark.com/vpn/api/v1/server/configurations"
	contents, err := fetchAndExtractFiles(ctx, client, zipURL)
	if err != nil {
		return nil, nil, err
	}
	mapping := surfsharkSubdomainToRegion()
	for fileName, content := range contents {
		if err := ctx.Err(); err != nil {
			return nil, warnings, err
		}
		if strings.HasSuffix(fileName, "_tcp.ovpn") {
			continue // only parse UDP files
		}
		host, warning, err := extractHostFromOVPN(content)
		if len(warning) > 0 {
			warnings = append(warnings, warning)
		}
		if err != nil {
			return nil, warnings, fmt.Errorf("%w in %s", err, fileName)
		}
		const repetition = 5
		IPs, err := resolveRepeat(ctx, lookupIP, host, repetition)
		if err != nil {
			return nil, warnings, err
		} else if len(IPs) == 0 {
			warning := fmt.Sprintf("no IP address found for host %q", host)
			warnings = append(warnings, warning)
			continue
		}
		subdomain := strings.TrimSuffix(host, ".prod.surfshark.com")
		region, ok := mapping[subdomain]
		if ok {
			delete(mapping, subdomain)
		} else {
			region = strings.TrimSuffix(host, ".prod.surfshark.com")
			warning := fmt.Sprintf("subdomain %q not found in Surfshark mapping", subdomain)
			warnings = append(warnings, warning)
		}
		server := models.SurfsharkServer{
			Region: region,
			IPs:    uniqueSortedIPs(IPs),
		}
		servers = append(servers, server)
	}

	// process entries in mapping that were not in zip file
	remainingServers, newWarnings, err := getRemainingServers(ctx, mapping, lookupIP)
	warnings = append(warnings, newWarnings...)
	if err != nil {
		return nil, warnings, err
	}
	servers = append(servers, remainingServers...)

	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Region < servers[j].Region
	})
	return servers, warnings, nil
}

func getRemainingServers(ctx context.Context, mapping map[string]string, lookupIP lookupIPFunc) (
	servers []models.SurfsharkServer, warnings []string, err error) {
	for subdomain, region := range mapping {
		if err := ctx.Err(); err != nil {
			return servers, warnings, err
		}
		host := subdomain + ".prod.surfshark.com"
		const repetition = 3
		IPs, err := resolveRepeat(ctx, lookupIP, host, repetition)
		if err != nil {
			warning := fmt.Sprintf("subdomain %q for region %q from mapping: %s", subdomain, region, err)
			warnings = append(warnings, warning)
			continue
		} else if len(IPs) == 0 {
			warning := fmt.Sprintf("subdomain %q for region %q from mapping did not resolve to any IP address",
				subdomain, region)
			warnings = append(warnings, warning)
			continue
		}
		server := models.SurfsharkServer{
			Region: region,
			IPs:    uniqueSortedIPs(IPs),
		}
		servers = append(servers, server)
	}
	return servers, warnings, nil
}

func stringifySurfsharkServers(servers []models.SurfsharkServer) (s string) {
	s = "func SurfsharkServers() []models.SurfsharkServer {\n"
	s += "	return []models.SurfsharkServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}

func surfsharkSubdomainToRegion() (mapping map[string]string) {
	return map[string]string{
		"ae-dub":       "United Arab Emirates",
		"al-tia":       "Albania",
		"at-vie":       "Austria",
		"au-adl":       "Australia Adelaide",
		"au-bne":       "Australia Brisbane",
		"au-mel":       "Australia Melbourne",
		"au-per":       "Australia Perth",
		"au-syd":       "Australia Sydney",
		"au-us":        "Australia US",
		"az-bak":       "Azerbaijan",
		"ba-sjj":       "Bosnia and Herzegovina",
		"be-bru":       "Belgium",
		"bg-sof":       "Bulgaria",
		"br-sao":       "Brazil",
		"ca-mon":       "Canada Montreal",
		"ca-tor":       "Canada Toronto",
		"ca-us":        "Canada US",
		"ca-van":       "Canada Vancouver",
		"ch-zur":       "Switzerland",
		"cl-san":       "Chile",
		"co-bog":       "Colombia",
		"cr-sjn":       "Costa Rica",
		"cy-nic":       "Cyprus",
		"cz-prg":       "Czech Republic",
		"de-ber":       "Germany Berlin",
		"de-fra":       "Germany Frankfurt am Main",
		"de-fra-st001": "Germany Frankfurt am Main st001",
		"de-fra-st002": "Germany Frankfurt am Main st002",
		"de-fra-st003": "Germany Frankfurt am Main st003",
		"de-muc":       "Germany Munich",
		"de-nue":       "Germany Nuremberg",
		"de-sg":        "Germany Singapour",
		"de-uk":        "Germany UK",
		"dk-cph":       "Denmark",
		"ee-tll":       "Estonia",
		"es-bcn":       "Spain Barcelona",
		"es-mad":       "Spain Madrid",
		"es-vlc":       "Spain Valencia",
		"fi-hel":       "Finland",
		"fr-bod":       "France Bordeaux",
		"fr-mrs":       "France Marseilles",
		"fr-par":       "France Paris",
		"fr-se":        "France Sweden",
		"gr-ath":       "Greece",
		"hk-hkg":       "Hong Kong",
		"hr-zag":       "Croatia",
		"hu-bud":       "Hungary",
		"id-jak":       "Indonesia",
		"ie-dub":       "Ireland",
		"il-tlv":       "Israel",
		"in-chn":       "India Chennai",
		"in-idr":       "India Indore",
		"in-mum":       "India Mumbai",
		"in-uk":        "India UK",
		"is-rkv":       "Iceland",
		"it-mil":       "Italy Milan",
		"it-rom":       "Italy Rome",
		"jp-tok":       "Japan Tokyo",
		"jp-tok-st001": "Japan Tokyo st001",
		"jp-tok-st002": "Japan Tokyo st002",
		"jp-tok-st003": "Japan Tokyo st003",
		"jp-tok-st004": "Japan Tokyo st004",
		"jp-tok-st005": "Japan Tokyo st005",
		"jp-tok-st006": "Japan Tokyo st006",
		"jp-tok-st007": "Japan Tokyo st007",
		"kr-seo":       "Korea",
		"kz-ura":       "Kazakhstan",
		"lu-ste":       "Luxembourg",
		"lv-rig":       "Latvia",
		"ly-tip":       "Libya",
		"md-chi":       "Moldova",
		"mk-skp":       "North Macedonia",
		"my-kul":       "Malaysia",
		"ng-lag":       "Nigeria",
		"nl-ams":       "Netherlands Amsterdam",
		"nl-ams-st001": "Netherlands Amsterdam st001",
		"nl-us":        "Netherlands US",
		"no-osl":       "Norway",
		"nz-akl":       "New Zealand",
		"ph-mnl":       "Philippines",
		"pl-gdn":       "Poland Gdansk",
		"pl-waw":       "Poland Warsaw",
		"pt-lis":       "Portugal Lisbon",
		"pt-lou":       "Portugal Loule",
		"pt-opo":       "Portugal Porto",
		"py-asu":       "Paraguay",
		"ro-buc":       "Romania",
		"rs-beg":       "Serbia",
		"ru-mos":       "Russia Moscow",
		"ru-spt":       "Russia St. Petersburg",
		"se-sto":       "Sweden",
		"sg-hk":        "Singapore Hong Kong",
		"sg-nl":        "Singapore Netherlands",
		"sg-sng":       "Singapore",
		"sg-in":        "Singapore in",
		"sg-sng-st001": "Singapore st001",
		"sg-sng-st002": "Singapore st002",
		"sg-sng-st003": "Singapore st003",
		"sg-sng-st004": "Singapore st004",
		"sg-sng-mp001": "Singapore mp001",
		"si-lju":       "Slovenia",
		"sk-bts":       "Slovekia",
		"th-bkk":       "Thailand",
		"tr-bur":       "Turkey",
		"tw-tai":       "Taiwan",
		"ua-iev":       "Ukraine",
		"uk-de":        "UK Germany",
		"uk-fr":        "UK France",
		"uk-gla":       "UK Glasgow",
		"uk-lon":       "UK London",
		"uk-lon-mp001": "UK London mp001",
		"uk-lon-st001": "UK London st001",
		"uk-lon-st002": "UK London st002",
		"uk-lon-st003": "UK London st003",
		"uk-lon-st004": "UK London st004",
		"uk-lon-st005": "UK London st005",
		"uk-man":       "UK Manchester",
		"us-atl":       "US Atlanta",
		"us-bdn":       "US Bend",
		"us-bos":       "US Boston",
		"us-buf":       "US Buffalo",
		"us-chi":       "US Chicago",
		"us-clt":       "US Charlotte",
		"us-dal":       "US Dallas",
		"us-den":       "US Denver",
		"us-dtw":       "US Gahanna",
		"us-hou":       "US Houston",
		"us-kan":       "US Kansas City",
		"us-las":       "US Las Vegas",
		"us-lax":       "US Los Angeles",
		"us-ltm":       "US Latham",
		"us-mia":       "US Miami",
		"us-mnz":       "US Maryland",
		"us-nl":        "US Netherlands",
		"us-nyc":       "US New York City",
		"us-nyc-mp001": "US New York City mp001",
		"us-nyc-st001": "US New York City st001",
		"us-nyc-st002": "US New York City st002",
		"us-nyc-st003": "US New York City st003",
		"us-nyc-st004": "US New York City st004",
		"us-nyc-st005": "US New York City st005",
		"us-orl":       "US Orlando",
		"us-phx":       "US Phoenix",
		"us-pt":        "US Portugal",
		"us-sea":       "US Seatle",
		"us-sfo":       "US San Francisco",
		"us-slc":       "US Salt Lake City",
		"us-stl":       "US Saint Louis",
		"us-tpa":       "US Tampa",
		"vn-hcm":       "Vietnam",
		"za-jnb":       "South Africa",
		"ar-bua":       "Argentina Buenos Aires",
		"tr-ist":       "Turkey Istanbul",
		"mx-mex":       "Mexico City Mexico",
		"ca-tor-mp001": "Canada Toronto mp001",
		"de-fra-mp001": "Germany Frankfurt mp001",
		"nl-ams-mp001": "Netherlands Amsterdam mp001",
		"us-sfo-mp001": "US San Francisco mp001",
	}
}
