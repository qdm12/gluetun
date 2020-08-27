package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
)

func main() {
	ctx := context.Background()
	os.Exit(_main(ctx))
}

func _main(ctx context.Context) int {
	resolverAddress := flag.String("resolver", "1.1.1.1", "DNS Resolver IP address to use")
	provider := flag.String("provider", "pia", "VPN provider to resolve for, 'pia', 'windscribe', 'cyberghost', 'vyprvpn' or 'purevpn'")
	region := flag.String("region", "all", "Comma separated list of VPN provider region names to resolve for, use 'all' to resolve all")
	flag.Parse()

	resolver := newResolver(*resolverAddress)
	lookupIP := newLookupIP(resolver)

	var domain string
	var servers []server
	switch *provider {
	case "windscribe":
		domain = "windscribe.com"
		servers = windscribeServers()
	case "surfshark":
		domain = "prod.surfshark.com"
		servers = surfsharkServers()
	case "cyberghost":
		domain = "cg-dialup.net"
		servers = cyberghostServers()
	case "vyprvpn":
		domain = "vyprvpn.com"
		servers = vyprvpnServers()
	case "purevpn":
		domain = "pointtoserver.com"
		servers = purevpnServers()
	default:
		fmt.Printf("Provider %q is not supported\n", *provider)
		return 1
	}
	if *region != "all" {
		regions := strings.Split(*region, ",")
		uniqueRegions := make(map[string]struct{})
		for _, r := range regions {
			uniqueRegions[r] = struct{}{}
		}
		for i := range servers {
			if _, ok := uniqueRegions[servers[i].region]; !ok {
				servers[i] = servers[len(servers)-1]
				servers = servers[:len(servers)-1]
			}
		}
	}

	stringChannel := make(chan string)
	errorChannel := make(chan error)
	const maxGoroutines = 10
	guard := make(chan struct{}, maxGoroutines)
	for _, s := range servers {
		go func(s server) {
			guard <- struct{}{}
			ips, err := resolveRepeat(ctx, lookupIP, s.subdomain+"."+domain, 3)
			<-guard
			if err != nil {
				errorChannel <- err
				return
			}
			stringChannel <- formatLine(*provider, s, ips)
		}(s)
	}
	var lines []string
	var errors []error
	for range servers {
		select {
		case err := <-errorChannel:
			errors = append(errors, err)
		case s := <-stringChannel:
			lines = append(lines, s)
		}
	}
	sort.Slice(lines, func(i, j int) bool {
		return lines[i] < lines[j]
	})
	for _, s := range lines {
		fmt.Println(s)
	}
	if len(errors) > 0 {
		fmt.Printf("\n%d errors occurred, described below\n\n", len(errors))
		for _, err := range errors {
			fmt.Println(err)
		}
		return 1
	}
	return 0
}

func formatLine(provider string, s server, ips []net.IP) string {
	ipStrings := make([]string, len(ips))
	for i := range ips {
		ipStrings[i] = fmt.Sprintf("{%s}", strings.ReplaceAll(ips[i].String(), ".", ", "))
	}
	ipString := strings.Join(ipStrings, ", ")
	switch provider {
	case "pia":
		return fmt.Sprintf(
			"{Region: %q, IPs: []net.IP{%s}},",
			s.region, ipString,
		)
	case "windscribe":
		return fmt.Sprintf(
			"{Region: %q, IPs: []net.IP{%s}},",
			s.region, ipString,
		)
	case "surfshark":
		return fmt.Sprintf(
			"{Region: %q, IPs: []net.IP{%s}},",
			s.region, ipString,
		)
	case "cyberghost":
		return fmt.Sprintf(
			"{Region: %q, Group: %q, IPs: []net.IP{%s}},",
			s.region, s.group, ipString,
		)
	case "vyprvpn":
		return fmt.Sprintf(
			"{Region: %q, IPs: []net.IP{%s}},",
			s.region, ipString,
		)
	case "purevpn":
		return fmt.Sprintf(
			"{Region: %q, Country: %q, City: %q, IPs: []net.IP{%s}},",
			s.region, s.country, s.city, ipString,
		)
	}
	return ""
}

type lookupIPFunc func(ctx context.Context, host string) (ips []net.IP, err error)

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

func newResolver(ip string) *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", net.JoinHostPort(ip, "53"))
		},
	}
}

func resolveRepeat(ctx context.Context, lookupIP lookupIPFunc, host string, n int) (ips []net.IP, err error) {
	for i := 0; i < n; i++ {
		newIPs, err := lookupIP(ctx, host)
		if err != nil {
			return nil, err
		}
		ips = append(ips, newIPs...)
	}
	return uniqueSortedIPs(ips), nil
}

func uniqueSortedIPs(ips []net.IP) []net.IP {
	uniqueIPs := make(map[string]struct{})
	for _, ip := range ips {
		uniqueIPs[ip.String()] = struct{}{}
	}
	ips = make([]net.IP, len(uniqueIPs))
	i := 0
	for ip := range uniqueIPs {
		ips[i] = net.ParseIP(ip)
		i++
	}
	sort.Slice(ips, func(i, j int) bool {
		return bytes.Compare(ips[i], ips[j]) < 0
	})
	return ips
}

type server struct {
	subdomain string
	region    string
	group     string // only for cyberghost
	country   string // only for purevpn
	city      string // only for purevpn
}

func windscribeServers() []server {
	return []server{
		{subdomain: "al", region: "Albania"},
		{subdomain: "ar", region: "Argentina"},
		{subdomain: "au", region: "Australia"},
		{subdomain: "at", region: "Austria"},
		{subdomain: "az", region: "Azerbaijan"},
		{subdomain: "be", region: "Belgium"},
		{subdomain: "ba", region: "Bosnia"},
		{subdomain: "br", region: "Brazil"},
		{subdomain: "bg", region: "Bulgaria"},
		{subdomain: "ca", region: "Canada East"},
		{subdomain: "ca-west", region: "Canada West"},
		{subdomain: "co", region: "Colombia"},
		{subdomain: "hr", region: "Croatia"},
		{subdomain: "cy", region: "Cyprus"},
		{subdomain: "cz", region: "Czech republic"},
		{subdomain: "dk", region: "Denmark"},
		{subdomain: "ee", region: "Estonia"},
		{subdomain: "aq", region: "Fake antarctica"},
		{subdomain: "fi", region: "Finland"},
		{subdomain: "fr", region: "France"},
		{subdomain: "ge", region: "Georgia"},
		{subdomain: "de", region: "Germany"},
		{subdomain: "gr", region: "Greece"},
		{subdomain: "hk", region: "Hong kong"},
		{subdomain: "hu", region: "Hungary"},
		{subdomain: "is", region: "Iceland"},
		{subdomain: "in", region: "India"},
		{subdomain: "id", region: "Indonesia"},
		{subdomain: "ie", region: "Ireland"},
		{subdomain: "il", region: "Israel"},
		{subdomain: "it", region: "Italy"},
		{subdomain: "jp", region: "Japan"},
		{subdomain: "lv", region: "Latvia"},
		{subdomain: "lt", region: "Lithuania"},
		{subdomain: "mk", region: "Macedonia"},
		{subdomain: "my", region: "Malaysia"},
		{subdomain: "mx", region: "Mexico"},
		{subdomain: "md", region: "Moldova"},
		{subdomain: "nl", region: "Netherlands"},
		{subdomain: "nz", region: "New zealand"},
		{subdomain: "no", region: "Norway"},
		{subdomain: "ph", region: "Philippines"},
		{subdomain: "pl", region: "Poland"},
		{subdomain: "pt", region: "Portugal"},
		{subdomain: "ro", region: "Romania"},
		{subdomain: "ru", region: "Russia"},
		{subdomain: "rs", region: "Serbia"},
		{subdomain: "sg", region: "Singapore"},
		{subdomain: "sk", region: "Slovakia"},
		{subdomain: "si", region: "Slovenia"},
		{subdomain: "za", region: "South Africa"},
		{subdomain: "kr", region: "South Korea"},
		{subdomain: "es", region: "Spain"},
		{subdomain: "se", region: "Sweden"},
		{subdomain: "ch", region: "Switzerland"},
		{subdomain: "th", region: "Thailand"},
		{subdomain: "tn", region: "Tunisia"},
		{subdomain: "tr", region: "Turkey"},
		{subdomain: "ua", region: "Ukraine"},
		{subdomain: "ae", region: "United Arab Emirates"},
		{subdomain: "uk", region: "United Kingdom"},
		{subdomain: "us-central", region: "US Central"},
		{subdomain: "us-east", region: "US East"},
		{subdomain: "us-west", region: "US West"},
		{subdomain: "vn", region: "Vietnam"},
		{subdomain: "wf-ca", region: "Windflix CA"},
		{subdomain: "wf-jp", region: "Windflix JP"},
		{subdomain: "wf-uk", region: "Windflix UK"},
		{subdomain: "wf-us", region: "Windflix US"},
	}
}

func surfsharkServers() []server {
	return []server{
		{subdomain: "ae-dub", region: "United Arab Emirates"},
		{subdomain: "al-tia", region: "Albania"},
		{subdomain: "at-vie", region: "Austria"},
		{subdomain: "au-adl", region: "Australia Adelaide"},
		{subdomain: "au-bne", region: "Australia Brisbane"},
		{subdomain: "au-mel", region: "Australia Melbourne"},
		{subdomain: "au-per", region: "Australia Perth"},
		{subdomain: "au-syd", region: "Australia Sydney"},
		{subdomain: "au-us", region: "Australia US"},
		{subdomain: "az-bak", region: "Azerbaijan"},
		{subdomain: "ba-sjj", region: "Bosnia and Herzegovina"},
		{subdomain: "be-bru", region: "Belgium"},
		{subdomain: "bg-sof", region: "Bulgaria"},
		{subdomain: "br-sao", region: "Brazil"},
		{subdomain: "ca-mon", region: "Canada Montreal"},
		{subdomain: "ca-tor", region: "Canada Toronto"},
		{subdomain: "ca-us", region: "Canada US"},
		{subdomain: "ca-van", region: "Canada Vancouver"},
		{subdomain: "ch-zur", region: "Switzerland"},
		{subdomain: "cl-san", region: "Chile"},
		{subdomain: "co-bog", region: "Colombia"},
		{subdomain: "cr-sjn", region: "Costa Rica"},
		{subdomain: "cy-nic", region: "Cyprus"},
		{subdomain: "cz-prg", region: "Czech Republic"},
		{subdomain: "de-ber", region: "Germany Berlin"},
		{subdomain: "de-fra", region: "Germany Frankfurt am Main"},
		{subdomain: "de-fra-st001", region: "Germany Frankfurt am Main st001"},
		{subdomain: "de-fra-st002", region: "Germany Frankfurt am Main st002"},
		{subdomain: "de-fra-st003", region: "Germany Frankfurt am Main st003"},
		{subdomain: "de-muc", region: "Germany Munich"},
		{subdomain: "de-nue", region: "Germany Nuremberg"},
		{subdomain: "de-sg", region: "Germany Singapour"},
		{subdomain: "de-uk", region: "Germany UK"},
		{subdomain: "dk-cph", region: "Denmark"},
		{subdomain: "ee-tll", region: "Estonia"},
		{subdomain: "es-bcn", region: "Spain Barcelona"},
		{subdomain: "es-mad", region: "Spain Madrid"},
		{subdomain: "es-vlc", region: "Spain Valencia"},
		{subdomain: "fi-hel", region: "Finland"},
		{subdomain: "fr-bod", region: "France Bordeaux"},
		{subdomain: "fr-mrs", region: "France Marseilles"},
		{subdomain: "fr-par", region: "France Paris"},
		{subdomain: "fr-se", region: "France Sweden"},
		{subdomain: "gr-ath", region: "Greece"},
		{subdomain: "hk-hkg", region: "Hong Kong"},
		{subdomain: "hr-zag", region: "Croatia"},
		{subdomain: "hu-bud", region: "Hungary"},
		{subdomain: "id-jak", region: "Indonesia"},
		{subdomain: "ie-dub", region: "Ireland"},
		{subdomain: "il-tlv", region: "Israel"},
		{subdomain: "in-chn", region: "India Chennai"},
		{subdomain: "in-idr", region: "India Indore"},
		{subdomain: "in-mum", region: "India Mumbai"},
		{subdomain: "in-uk", region: "India UK"},
		{subdomain: "is-rkv", region: "Iceland"},
		{subdomain: "it-mil", region: "Italy Milan"},
		{subdomain: "it-rom", region: "Italy Rome"},
		{subdomain: "jp-tok", region: "Japan Tokyo"},
		{subdomain: "jp-tok-st001", region: "Japan Tokyo st001"},
		{subdomain: "jp-tok-st002", region: "Japan Tokyo st002"},
		{subdomain: "jp-tok-st003", region: "Japan Tokyo st003"},
		{subdomain: "jp-tok-st004", region: "Japan Tokyo st004"},
		{subdomain: "jp-tok-st005", region: "Japan Tokyo st005"},
		{subdomain: "jp-tok-st006", region: "Japan Tokyo st006"},
		{subdomain: "jp-tok-st007", region: "Japan Tokyo st007"},
		{subdomain: "kr-seo", region: "Korea"},
		{subdomain: "kz-ura", region: "Kazakhstan"},
		{subdomain: "lu-ste", region: "Luxembourg"},
		{subdomain: "lv-rig", region: "Latvia"},
		{subdomain: "ly-tip", region: "Libya"},
		{subdomain: "md-chi", region: "Moldova"},
		{subdomain: "mk-skp", region: "North Macedonia"},
		{subdomain: "my-kul", region: "Malaysia"},
		{subdomain: "ng-lag", region: "Nigeria"},
		{subdomain: "nl-ams", region: "Netherlands Amsterdam"},
		{subdomain: "nl-ams-st001", region: "Netherlands Amsterdam st001"},
		{subdomain: "nl-us", region: "Netherlands US"},
		{subdomain: "no-osl", region: "Norway"},
		{subdomain: "nz-akl", region: "New Zealand"},
		{subdomain: "ph-mnl", region: "Philippines"},
		{subdomain: "pl-gdn", region: "Poland Gdansk"},
		{subdomain: "pl-waw", region: "Poland Warsaw"},
		{subdomain: "pt-lis", region: "Portugal Lisbon"},
		{subdomain: "pt-lou", region: "Portugal Loule"},
		{subdomain: "pt-opo", region: "Portugal Porto"},
		{subdomain: "py-asu", region: "Paraguay"},
		{subdomain: "ro-buc", region: "Romania"},
		{subdomain: "rs-beg", region: "Serbia"},
		{subdomain: "ru-mos", region: "Russia Moscow"},
		{subdomain: "ru-spt", region: "Russia St. Petersburg"},
		{subdomain: "se-sto", region: "Sweden"},
		{subdomain: "sg-hk", region: "Singapore Hong Kong"},
		{subdomain: "sg-nl", region: "Singapore Netherlands"},
		{subdomain: "sg-sng", region: "Singapore"},
		{subdomain: "sg-sng-st001", region: "Singapore st001"},
		{subdomain: "sg-sng-st002", region: "Singapore st002"},
		{subdomain: "sg-sng-st003", region: "Singapore st003"},
		{subdomain: "sg-sng-st004", region: "Singapore st004"},
		{subdomain: "si-lju", region: "Slovenia"},
		{subdomain: "sk-bts", region: "Slovekia"},
		{subdomain: "th-bkk", region: "Thailand"},
		{subdomain: "tr-bur", region: "Turkey"},
		{subdomain: "tw-tai", region: "Taiwan"},
		{subdomain: "ua-iev", region: "Ukraine"},
		{subdomain: "uk-de", region: "UK Germany"},
		{subdomain: "uk-fr", region: "UK France"},
		{subdomain: "uk-gla", region: "UK Glasgow"},
		{subdomain: "uk-lon", region: "UK London"},
		{subdomain: "uk-lon-st001", region: "UK London st001"},
		{subdomain: "uk-lon-st002", region: "UK London st002"},
		{subdomain: "uk-lon-st003", region: "UK London st003"},
		{subdomain: "uk-lon-st004", region: "UK London st004"},
		{subdomain: "uk-lon-st005", region: "UK London st005"},
		{subdomain: "uk-man", region: "UK Manchester"},
		{subdomain: "us-atl", region: "US Atlanta"},
		{subdomain: "us-bdn", region: "US Bend"},
		{subdomain: "us-bos", region: "US Boston"},
		{subdomain: "us-buf", region: "US Buffalo"},
		{subdomain: "us-chi", region: "US Chicago"},
		{subdomain: "us-clt", region: "US Charlotte"},
		{subdomain: "us-dal", region: "US Dallas"},
		{subdomain: "us-den", region: "US Denver"},
		{subdomain: "us-dtw", region: "US Gahanna"},
		{subdomain: "us-hou", region: "US Houston"},
		{subdomain: "us-kan", region: "US Kansas City"},
		{subdomain: "us-las", region: "US Las Vegas"},
		{subdomain: "us-lax", region: "US Los Angeles"},
		{subdomain: "us-ltm", region: "US Latham"},
		{subdomain: "us-mia", region: "US Miami"},
		{subdomain: "us-mnz", region: "US Maryland"},
		{subdomain: "us-nl", region: "US Netherlands"},
		{subdomain: "us-nyc", region: "US New York City"},
		{subdomain: "us-nyc-mp001", region: "US New York City mp001"},
		{subdomain: "us-nyc-st001", region: "US New York City st001"},
		{subdomain: "us-nyc-st002", region: "US New York City st002"},
		{subdomain: "us-nyc-st003", region: "US New York City st003"},
		{subdomain: "us-nyc-st004", region: "US New York City st004"},
		{subdomain: "us-nyc-st005", region: "US New York City st005"},
		{subdomain: "us-orl", region: "US Orlando"},
		{subdomain: "us-phx", region: "US Phoenix"},
		{subdomain: "us-pt", region: "US Portugal"},
		{subdomain: "us-sea", region: "US Seatle"},
		{subdomain: "us-sfo", region: "US San Francisco"},
		{subdomain: "us-slc", region: "US Salt Lake City"},
		{subdomain: "us-stl", region: "US Saint Louis"},
		{subdomain: "us-tpa", region: "US Tampa"},
		{subdomain: "vn-hcm", region: "Vietnam"},
		{subdomain: "za-jnb", region: "South Africa"},
	}
}

func cyberghostServers() []server {
	return []server{
		{subdomain: "97-1-al", region: "Albania", group: "Premium TCP Europe"},
		{subdomain: "87-1-al", region: "Albania", group: "Premium UDP Europe"},
		{subdomain: "87-1-dz", region: "Algeria", group: "Premium UDP Europe"},
		{subdomain: "97-1-dz", region: "Algeria", group: "Premium TCP Europe"},
		{subdomain: "97-1-ad", region: "Andorra", group: "Premium TCP Europe"},
		{subdomain: "87-1-ad", region: "Andorra", group: "Premium UDP Europe"},
		{subdomain: "94-1-ar", region: "Argentina", group: "Premium UDP USA"},
		{subdomain: "93-1-ar", region: "Argentina", group: "Premium TCP USA"},
		{subdomain: "87-1-am", region: "Armenia", group: "Premium UDP Europe"},
		{subdomain: "97-1-am", region: "Armenia", group: "Premium TCP Europe"},
		{subdomain: "95-1-au", region: "Australia", group: "Premium UDP Asia"},
		{subdomain: "96-1-au", region: "Australia", group: "Premium TCP Asia"},
		{subdomain: "97-1-at", region: "Austria", group: "Premium TCP Europe"},
		{subdomain: "87-1-at", region: "Austria", group: "Premium UDP Europe"},
		{subdomain: "93-1-bs", region: "Bahamas", group: "Premium TCP USA"},
		{subdomain: "94-1-bs", region: "Bahamas", group: "Premium UDP USA"},
		{subdomain: "95-1-bd", region: "Bangladesh", group: "Premium UDP Asia"},
		{subdomain: "96-1-bd", region: "Bangladesh", group: "Premium TCP Asia"},
		{subdomain: "97-1-by", region: "Belarus", group: "Premium TCP Europe"},
		{subdomain: "87-1-by", region: "Belarus", group: "Premium UDP Europe"},
		{subdomain: "97-1-be", region: "Belgium", group: "Premium TCP Europe"},
		{subdomain: "87-1-be", region: "Belgium", group: "Premium UDP Europe"},
		{subdomain: "87-1-ba", region: "Bosnia and Herzegovina", group: "Premium UDP Europe"},
		{subdomain: "97-1-ba", region: "Bosnia and Herzegovina", group: "Premium TCP Europe"},
		{subdomain: "94-1-br", region: "Brazil", group: "Premium UDP USA"},
		{subdomain: "93-1-br", region: "Brazil", group: "Premium TCP USA"},
		{subdomain: "87-1-bg", region: "Bulgaria", group: "Premium UDP Europe"},
		{subdomain: "97-1-bg", region: "Bulgaria", group: "Premium TCP Europe"},
		{subdomain: "96-1-kh", region: "Cambodia", group: "Premium TCP Asia"},
		{subdomain: "95-1-kh", region: "Cambodia", group: "Premium UDP Asia"},
		{subdomain: "93-1-ca", region: "Canada", group: "Premium TCP USA"},
		{subdomain: "94-1-ca", region: "Canada", group: "Premium UDP USA"},
		{subdomain: "93-1-cl", region: "Chile", group: "Premium TCP USA"},
		{subdomain: "94-1-cl", region: "Chile", group: "Premium UDP USA"},
		{subdomain: "96-1-cn", region: "China", group: "Premium TCP Asia"},
		{subdomain: "95-1-cn", region: "China", group: "Premium UDP Asia"},
		{subdomain: "94-1-co", region: "Colombia", group: "Premium UDP USA"},
		{subdomain: "93-1-co", region: "Colombia", group: "Premium TCP USA"},
		{subdomain: "93-1-cr", region: "Costa Rica", group: "Premium TCP USA"},
		{subdomain: "94-1-cr", region: "Costa Rica", group: "Premium UDP USA"},
		{subdomain: "87-1-cy", region: "Cyprus", group: "Premium UDP Europe"},
		{subdomain: "97-1-cy", region: "Cyprus", group: "Premium TCP Europe"},
		{subdomain: "97-1-cz", region: "Czech Republic", group: "Premium TCP Europe"},
		{subdomain: "87-1-cz", region: "Czech Republic", group: "Premium UDP Europe"},
		{subdomain: "87-1-dk", region: "Denmark", group: "Premium UDP Europe"},
		{subdomain: "97-1-dk", region: "Denmark", group: "Premium TCP Europe"},
		{subdomain: "87-1-eg", region: "Egypt", group: "Premium UDP Europe"},
		{subdomain: "97-1-eg", region: "Egypt", group: "Premium TCP Europe"},
		{subdomain: "87-1-ee", region: "Estonia", group: "Premium UDP Europe"},
		{subdomain: "97-1-ee", region: "Estonia", group: "Premium TCP Europe"},
		{subdomain: "97-1-fi", region: "Finland", group: "Premium TCP Europe"},
		{subdomain: "87-1-fi", region: "Finland", group: "Premium UDP Europe"},
		{subdomain: "87-1-fr", region: "France", group: "Premium UDP Europe"},
		{subdomain: "97-1-fr", region: "France", group: "Premium TCP Europe"},
		{subdomain: "87-1-ge", region: "Georgia", group: "Premium UDP Europe"},
		{subdomain: "97-1-ge", region: "Georgia", group: "Premium TCP Europe"},
		{subdomain: "97-1-de", region: "Germany", group: "Premium TCP Europe"},
		{subdomain: "87-1-de", region: "Germany", group: "Premium UDP Europe"},
		{subdomain: "87-1-gr", region: "Greece", group: "Premium UDP Europe"},
		{subdomain: "97-1-gr", region: "Greece", group: "Premium TCP Europe"},
		{subdomain: "97-1-gl", region: "Greenland", group: "Premium TCP Europe"},
		{subdomain: "87-1-gl", region: "Greenland", group: "Premium UDP Europe"},
		{subdomain: "96-1-hk", region: "Hong Kong", group: "Premium TCP Asia"},
		{subdomain: "95-1-hk", region: "Hong Kong", group: "Premium UDP Asia"},
		{subdomain: "87-1-hu", region: "Hungary", group: "Premium UDP Europe"},
		{subdomain: "97-1-hu", region: "Hungary", group: "Premium TCP Europe"},
		{subdomain: "97-1-is", region: "Iceland", group: "Premium TCP Europe"},
		{subdomain: "87-1-is", region: "Iceland", group: "Premium UDP Europe"},
		{subdomain: "87-1-in", region: "India", group: "Premium UDP Europe"},
		{subdomain: "97-1-in", region: "India", group: "Premium TCP Europe"},
		{subdomain: "95-1-id", region: "Indonesia", group: "Premium UDP Asia"},
		{subdomain: "96-1-id", region: "Indonesia", group: "Premium TCP Asia"},
		{subdomain: "87-1-ir", region: "Iran", group: "Premium UDP Europe"},
		{subdomain: "97-1-ir", region: "Iran", group: "Premium TCP Europe"},
		{subdomain: "87-1-ie", region: "Ireland", group: "Premium UDP Europe"},
		{subdomain: "97-1-ie", region: "Ireland", group: "Premium TCP Europe"},
		{subdomain: "87-1-im", region: "Isle of Man", group: "Premium UDP Europe"},
		{subdomain: "97-1-im", region: "Isle of Man", group: "Premium TCP Europe"},
		{subdomain: "87-1-il", region: "Israel", group: "Premium UDP Europe"},
		{subdomain: "97-1-il", region: "Israel", group: "Premium TCP Europe"},
		{subdomain: "97-1-it", region: "Italy", group: "Premium TCP Europe"},
		{subdomain: "87-1-it", region: "Italy", group: "Premium UDP Europe"},
		{subdomain: "95-1-jp", region: "Japan", group: "Premium UDP Asia"},
		{subdomain: "96-1-jp", region: "Japan", group: "Premium TCP Asia"},
		{subdomain: "97-1-kz", region: "Kazakhstan", group: "Premium TCP Europe"},
		{subdomain: "87-1-kz", region: "Kazakhstan", group: "Premium UDP Europe"},
		{subdomain: "95-1-ke", region: "Kenya", group: "Premium UDP Asia"},
		{subdomain: "96-1-ke", region: "Kenya", group: "Premium TCP Asia"},
		{subdomain: "95-1-kr", region: "Korea", group: "Premium UDP Asia"},
		{subdomain: "96-1-kr", region: "Korea", group: "Premium TCP Asia"},
		{subdomain: "97-1-lv", region: "Latvia", group: "Premium TCP Europe"},
		{subdomain: "87-1-lv", region: "Latvia", group: "Premium UDP Europe"},
		{subdomain: "97-1-li", region: "Liechtenstein", group: "Premium TCP Europe"},
		{subdomain: "87-1-li", region: "Liechtenstein", group: "Premium UDP Europe"},
		{subdomain: "97-1-lt", region: "Lithuania", group: "Premium TCP Europe"},
		{subdomain: "87-1-lt", region: "Lithuania", group: "Premium UDP Europe"},
		{subdomain: "87-1-lu", region: "Luxembourg", group: "Premium UDP Europe"},
		{subdomain: "97-1-lu", region: "Luxembourg", group: "Premium TCP Europe"},
		{subdomain: "96-1-mo", region: "Macao", group: "Premium TCP Asia"},
		{subdomain: "95-1-mo", region: "Macao", group: "Premium UDP Asia"},
		{subdomain: "97-1-mk", region: "Macedonia", group: "Premium TCP Europe"},
		{subdomain: "87-1-mk", region: "Macedonia", group: "Premium UDP Europe"},
		{subdomain: "95-1-my", region: "Malaysia", group: "Premium UDP Asia"},
		{subdomain: "96-1-my", region: "Malaysia", group: "Premium TCP Asia"},
		{subdomain: "87-1-mt", region: "Malta", group: "Premium UDP Europe"},
		{subdomain: "97-1-mt", region: "Malta", group: "Premium TCP Europe"},
		{subdomain: "93-1-mx", region: "Mexico", group: "Premium TCP USA"},
		{subdomain: "94-1-mx", region: "Mexico", group: "Premium UDP USA"},
		{subdomain: "87-1-md", region: "Moldova", group: "Premium UDP Europe"},
		{subdomain: "97-1-md", region: "Moldova", group: "Premium TCP Europe"},
		{subdomain: "87-1-mc", region: "Monaco", group: "Premium UDP Europe"},
		{subdomain: "97-1-mc", region: "Monaco", group: "Premium TCP Europe"},
		{subdomain: "96-1-mn", region: "Mongolia", group: "Premium TCP Asia"},
		{subdomain: "95-1-mn", region: "Mongolia", group: "Premium UDP Asia"},
		{subdomain: "87-1-me", region: "Montenegro", group: "Premium UDP Europe"},
		{subdomain: "97-1-me", region: "Montenegro", group: "Premium TCP Europe"},
		{subdomain: "97-1-ma", region: "Morocco", group: "Premium TCP Europe"},
		{subdomain: "87-1-ma", region: "Morocco", group: "Premium UDP Europe"},
		{subdomain: "97-1-nl", region: "Netherlands", group: "Premium TCP Europe"},
		{subdomain: "87-1-nl", region: "Netherlands", group: "Premium UDP Europe"},
		{subdomain: "95-1-nz", region: "New Zealand", group: "Premium UDP Asia"},
		{subdomain: "96-1-nz", region: "New Zealand", group: "Premium TCP Asia"},
		{subdomain: "87-1-ng", region: "Nigeria", group: "Premium UDP Europe"},
		{subdomain: "97-1-ng", region: "Nigeria", group: "Premium TCP Europe"},
		{subdomain: "97-1-no", region: "Norway", group: "Premium TCP Europe"},
		{subdomain: "87-1-no", region: "Norway", group: "Premium UDP Europe"},
		{subdomain: "97-1-pk", region: "Pakistan", group: "Premium TCP Europe"},
		{subdomain: "87-1-pk", region: "Pakistan", group: "Premium UDP Europe"},
		{subdomain: "97-1-pa", region: "Panama", group: "Premium TCP Europe"},
		{subdomain: "87-1-pa", region: "Panama", group: "Premium UDP Europe"},
		{subdomain: "95-1-ph", region: "Philippines", group: "Premium UDP Asia"},
		{subdomain: "96-1-ph", region: "Philippines", group: "Premium TCP Asia"},
		{subdomain: "97-1-pl", region: "Poland", group: "Premium TCP Europe"},
		{subdomain: "87-1-pl", region: "Poland", group: "Premium UDP Europe"},
		{subdomain: "97-1-pt", region: "Portugal", group: "Premium TCP Europe"},
		{subdomain: "87-1-pt", region: "Portugal", group: "Premium UDP Europe"},
		{subdomain: "97-1-qa", region: "Qatar", group: "Premium TCP Europe"},
		{subdomain: "87-1-qa", region: "Qatar", group: "Premium UDP Europe"},
		{subdomain: "87-1-ro", region: "Romania", group: "Premium UDP Europe"},
		{subdomain: "87-8-ro", region: "Romania", group: "NoSpy UDP Europe"},
		{subdomain: "97-8-ro", region: "Romania", group: "NoSpy TCP Europe"},
		{subdomain: "97-1-ro", region: "Romania", group: "Premium TCP Europe"},
		{subdomain: "87-1-ru", region: "Russian Federation", group: "Premium UDP Europe"},
		{subdomain: "97-1-ru", region: "Russian Federation", group: "Premium TCP Europe"},
		{subdomain: "97-1-sa", region: "Saudi Arabia", group: "Premium TCP Europe"},
		{subdomain: "87-1-sa", region: "Saudi Arabia", group: "Premium UDP Europe"},
		{subdomain: "97-1-rs", region: "Serbia", group: "Premium TCP Europe"},
		{subdomain: "87-1-rs", region: "Serbia", group: "Premium UDP Europe"},
		{subdomain: "95-1-sg", region: "Singapore", group: "Premium UDP Asia"},
		{subdomain: "96-1-sg", region: "Singapore", group: "Premium TCP Asia"},
		{subdomain: "87-1-sk", region: "Slovakia", group: "Premium UDP Europe"},
		{subdomain: "97-1-sk", region: "Slovakia", group: "Premium TCP Europe"},
		{subdomain: "87-1-si", region: "Slovenia", group: "Premium UDP Europe"},
		{subdomain: "97-1-si", region: "Slovenia", group: "Premium TCP Europe"},
		{subdomain: "87-1-za", region: "South Africa", group: "Premium UDP Europe"},
		{subdomain: "95-1-za", region: "South Africa", group: "Premium UDP Asia"},
		{subdomain: "97-1-za", region: "South Africa", group: "Premium TCP Europe"},
		{subdomain: "96-1-za", region: "South Africa", group: "Premium TCP Asia"},
		{subdomain: "97-1-es", region: "Spain", group: "Premium TCP Europe"},
		{subdomain: "87-1-es", region: "Spain", group: "Premium UDP Europe"},
		{subdomain: "97-1-lk", region: "Sri Lanka", group: "Premium TCP Europe"},
		{subdomain: "87-1-lk", region: "Sri Lanka", group: "Premium UDP Europe"},
		{subdomain: "97-1-se", region: "Sweden", group: "Premium TCP Europe"},
		{subdomain: "87-1-se", region: "Sweden", group: "Premium UDP Europe"},
		{subdomain: "87-1-ch", region: "Switzerland", group: "Premium UDP Europe"},
		{subdomain: "97-1-ch", region: "Switzerland", group: "Premium TCP Europe"},
		{subdomain: "96-1-tw", region: "Taiwan", group: "Premium TCP Asia"},
		{subdomain: "95-1-tw", region: "Taiwan", group: "Premium UDP Asia"},
		{subdomain: "96-1-th", region: "Thailand", group: "Premium TCP Asia"},
		{subdomain: "95-1-th", region: "Thailand", group: "Premium UDP Asia"},
		{subdomain: "87-1-tr", region: "Turkey", group: "Premium UDP Europe"},
		{subdomain: "97-1-tr", region: "Turkey", group: "Premium TCP Europe"},
		{subdomain: "97-1-ua", region: "Ukraine", group: "Premium TCP Europe"},
		{subdomain: "87-1-ua", region: "Ukraine", group: "Premium UDP Europe"},
		{subdomain: "87-1-ae", region: "United Arab Emirates", group: "Premium UDP Europe"},
		{subdomain: "97-1-ae", region: "United Arab Emirates", group: "Premium TCP Europe"},
		{subdomain: "97-1-gb", region: "United Kingdom", group: "Premium TCP Europe"},
		{subdomain: "87-1-gb", region: "United Kingdom", group: "Premium UDP Europe"},
		{subdomain: "94-1-us", region: "United States", group: "Premium UDP USA"},
		{subdomain: "93-1-us", region: "United States", group: "Premium TCP USA"},
		{subdomain: "87-1-ve", region: "Venezuela", group: "Premium UDP Europe"},
		{subdomain: "97-1-ve", region: "Venezuela", group: "Premium TCP Europe"},
		{subdomain: "95-1-vn", region: "Vietnam", group: "Premium UDP Asia"},
		{subdomain: "96-1-vn", region: "Vietnam", group: "Premium TCP Asia"},
	}
}

func vyprvpnServers() []server {
	return []server{
		{subdomain: "ae1", region: "Dubai"},
		{subdomain: "ar1", region: "Argentina"},
		{subdomain: "at1", region: "Austria"},
		{subdomain: "au1", region: "Australia Sydney"},
		{subdomain: "au2", region: "Australia Melbourne"},
		{subdomain: "au3", region: "Australia Perth"},
		{subdomain: "be1", region: "Belgium"},
		{subdomain: "bg1", region: "Bulgaria"},
		{subdomain: "bh1", region: "Bahrain"},
		{subdomain: "br1", region: "Brazil"},
		{subdomain: "ca1", region: "Canada"},
		{subdomain: "ch1", region: "Switzerland"},
		{subdomain: "co1", region: "Columbia"},
		{subdomain: "cr1", region: "Costa Rica"},
		{subdomain: "cz1", region: "Czech Republic"},
		{subdomain: "de1", region: "Germany"},
		{subdomain: "dk1", region: "Denmark"},
		{subdomain: "dz1", region: "Algeria"},
		{subdomain: "eg1", region: "Egypt"},
		{subdomain: "es1", region: "Spain"},
		{subdomain: "eu1", region: "Netherlands"},
		{subdomain: "fi1", region: "Finland"},
		{subdomain: "fr1", region: "France"},
		{subdomain: "gr1", region: "Greece"},
		{subdomain: "hk1", region: "Hong Kong"},
		{subdomain: "id1", region: "Indonesia"},
		{subdomain: "ie1", region: "Ireland"},
		{subdomain: "il1", region: "Israel"},
		{subdomain: "in1", region: "India"},
		{subdomain: "is1", region: "Iceland"},
		{subdomain: "it1", region: "Italy"},
		{subdomain: "jp1", region: "Japan"},
		{subdomain: "kr1", region: "South Korea"},
		{subdomain: "li1", region: "Liechtenstein"},
		{subdomain: "lt1", region: "Lithuania"},
		{subdomain: "lu1", region: "Luxembourg"},
		{subdomain: "lv1", region: "Latvia"},
		{subdomain: "mh1", region: "Marshall Islands"},
		{subdomain: "mo1", region: "Macao"},
		{subdomain: "mv1", region: "Maldives"},
		{subdomain: "mx1", region: "Mexico"},
		{subdomain: "my1", region: "Malaysia"},
		{subdomain: "no1", region: "Norway"},
		{subdomain: "nz1", region: "New Zealand"},
		{subdomain: "pa1", region: "Panama"},
		{subdomain: "ph1", region: "Philippines"},
		{subdomain: "pk1", region: "Pakistan"},
		{subdomain: "pl1", region: "Poland"},
		{subdomain: "pt1", region: "Portugal"},
		{subdomain: "qa1", region: "Qatar"},
		{subdomain: "ro1", region: "Romania"},
		{subdomain: "ru1", region: "Russia"},
		{subdomain: "sa1", region: "Saudi Arabia"},
		{subdomain: "se1", region: "Sweden"},
		{subdomain: "sg1", region: "Singapore"},
		{subdomain: "si1", region: "Slovenia"},
		{subdomain: "sk1", region: "Slovakia"},
		{subdomain: "sv1", region: "El Salvador"},
		{subdomain: "th1", region: "Thailand"},
		{subdomain: "tr1", region: "Turkey"},
		{subdomain: "tw1", region: "Taiwan"},
		{subdomain: "ua1", region: "Ukraine"},
		{subdomain: "uk1", region: "United Kingdom"},
		{subdomain: "us1", region: "USA Los Angeles"},
		{subdomain: "us2", region: "USA Washington DC"},
		{subdomain: "us3", region: "USA Austin"},
		{subdomain: "us4", region: "USA Miami"},
		{subdomain: "us5", region: "USA New York"},
		{subdomain: "us6", region: "USA Chicago"},
		{subdomain: "us7", region: "USA San Francisco"},
		{subdomain: "us8", region: "USA Seattle"},
		{subdomain: "uy1", region: "Uruguay"},
		{subdomain: "vn1", region: "Vietnam"},
	}
}

func purevpnServers() []server {
	servers := []server{
		{subdomain: "vlus-dz1-ovpn", region: "Africa", country: "Algeria", city: "Algiers"},
		{subdomain: "vlus-ao1-ovpn", region: "Africa", country: "Angola", city: "Benguela"},
		{subdomain: "vleu-cv-ovpn", region: "Africa", country: "Cape Verde", city: "Praia"},
		{subdomain: "vlus-eg1-ovpn", region: "Africa", country: "Egypt", city: "Cairo"},
		{subdomain: "et1-ovpn", region: "Africa", country: "Ethiopia", city: "Addis Ababa"},
		{subdomain: "gh1-ovpn", region: "Africa", country: "Ghana", city: "Accra"},
		{subdomain: "ke1-ovpn", region: "Africa", country: "Kenya", city: "Mombasa"},
		{subdomain: "vlus-mg1-ovpn", region: "Africa", country: "Madagascar", city: "Antananarivo"},
		{subdomain: "vlus-mr1-ovpn", region: "Africa", country: "Mauritania", city: "Nouakchott"},
		{subdomain: "mu1-ovpn", region: "Africa", country: "Mauritius", city: "Port Louis"},
		{subdomain: "ma1-ovpn", region: "Africa", country: "Morocco", city: "Rabat"},
		{subdomain: "vlus-ne1-ovpn", region: "Africa", country: "Niger", city: "Niamey"},
		{subdomain: "ng1-ovpn", region: "Africa", country: "Nigeria", city: "Suleja"},
		{subdomain: "vlus-sn1-ovpn", region: "Africa", country: "Senegal", city: "Dakar"},
		{subdomain: "sc1-ovpn", region: "Africa", country: "Seychelles", city: "Victoria"},
		{subdomain: "za2-ovpn", region: "Africa", country: "South Africa", city: "Johannesburg"},
		{subdomain: "vlus-tz1-ovpn", region: "Africa", country: "Tanzania", city: "Dar Es Salaam"},
		{subdomain: "vlus-tn1-ovpn", region: "Africa", country: "Tunisia", city: "Tunis"},
		{subdomain: "vlus-af1-ovpn", region: "Asia", country: "Afghanistan", city: "Kabul"},
		{subdomain: "sg2-ovpn", region: "Asia", country: "Armenia", city: "Singapore"},
		{subdomain: "az1-ovpn", region: "Asia", country: "Azerbaijan", city: "Baku"},
		{subdomain: "vlus-bd1-ovpn", region: "Asia", country: "Bangladesh", city: "Dhaka"},
		{subdomain: "bn2-ovpn", region: "Asia", country: "Brunei Darussalam", city: "Bandar Seri Begawan"},
		{subdomain: "kh1-ovpn", region: "Asia", country: "Cambodia", city: "Phnom Penh"},
		{subdomain: "hk2-ovpn", region: "Asia", country: "Hong Kong (SAR)", city: "Hong Kong"},
		{subdomain: "in2-ovpn", region: "Asia", country: "India", city: "Chennai"},
		{subdomain: "idn1-ovpn", region: "Asia", country: "Indonesia", city: "Jakarta"},
		{subdomain: "jp-tk1-ovpn", region: "Asia", country: "Japan", city: "Tokyo"},
		{subdomain: "vlus-kz1-ovpn", region: "Asia", country: "Kazakhstan", city: "Almaty"},
		{subdomain: "kr2-ovpn", region: "Asia", country: "Korea, South", city: "Seoul"},
		{subdomain: "vlus-kg1-ovpn", region: "Asia", country: "Kyrgyzstan", city: "Bishkek"},
		{subdomain: "vlus-la1-ovpn", region: "Asia", country: "Laos", city: "Vientiane"},
		{subdomain: "mo1-ovpn", region: "Asia", country: "Macao", city: "Beyrouth"},
		{subdomain: "my2-ovpn", region: "Asia", country: "Malaysia", city: "Johor Baharu"},
		{subdomain: "my-kl2-ovpn", region: "Asia", country: "Malaysia", city: "Kuala Lumpur"},
		{subdomain: "vlus-mn1-ovpn", region: "Asia", country: "Mongolia", city: "Ulaanbaatar"},
		{subdomain: "pk1-ovpn", region: "Asia", country: "Pakistan", city: "Islamabad"},
		{subdomain: "vlus-pg1-ovpn", region: "Asia", country: "Papua New Guinea", city: "Port Moresby"},
		{subdomain: "vlap-ph2-ovpn", region: "Asia", country: "Philippines", city: "Manila"},
		{subdomain: "vlus-lk1-ovpn", region: "Asia", country: "Sri Lanka", city: "Colombo"},
		{subdomain: "tw2-ovpn", region: "Asia", country: "Taiwan", city: "Taipei"},
		{subdomain: "vlus-tj-ovpn", region: "Asia", country: "Tajikistan", city: "Dushanbe"},
		{subdomain: "vlap-th2-ovpn", region: "Asia", country: "Thailand", city: "Bangkok"},
		{subdomain: "tr2-ovpn", region: "Asia", country: "Turkey", city: "Istanbul"},
		{subdomain: "vlus-tm1-ovpn", region: "Asia", country: "Turkmenistan", city: "Ashgabat"},
		{subdomain: "vlus-uz-ovpn", region: "Asia", country: "Uzbekistan", city: "Tashkent"},
		{subdomain: "vlap-vn2-ovpn", region: "Asia", country: "Vietnam", city: "Hanoi"},
		{subdomain: "al1-ovpn", region: "Europe", country: "Albania", city: "Tirane"},
		{subdomain: "vleu-am1-ovpn", region: "Europe", country: "Armenia", city: "Yerevan"},
		{subdomain: "at2-ovpn", region: "Europe", country: "Austria", city: "Vienna"},
		{subdomain: "vleu-be2-ovpn", region: "Europe", country: "Belgium", city: "Brussels"},
		{subdomain: "ba1-ovpn", region: "Europe", country: "Bosnia and Herzegovina", city: "Sarajevo"},
		{subdomain: "bg2-ovpn", region: "Europe", country: "Bulgaria", city: "Sofia"},
		{subdomain: "vlus-hr1-ovpn", region: "Europe", country: "Croatia", city: "Zagreb"},
		{subdomain: "cy1-ovpn", region: "Europe", country: "Cyprus", city: "Nicosia"},
		{subdomain: "dk2-ovpn", region: "Europe", country: "Denmark", city: "Copenhagen"},
		{subdomain: "ee1-ovpn", region: "Europe", country: "Estonia", city: "Tallinn"},
		{subdomain: "fr2-ovpn", region: "Europe", country: "France", city: "Paris"},
		{subdomain: "vlus-ge1-ovpn", region: "Europe", country: "Georgia", city: "Tbilisi"},
		{subdomain: "de2-ovpn", region: "Europe", country: "Germany", city: "Frankfurt"},
		{subdomain: "de2-ovpn", region: "Europe", country: "Germany", city: "Munich"},
		{subdomain: "de-ao1-ovpn", region: "Europe", country: "Germany", city: "Nuremberg"},
		{subdomain: "gr2-ovpn", region: "Europe", country: "Greece", city: "Thessaloniki"},
		{subdomain: "hu2-ovpn", region: "Europe", country: "Hungary", city: "Budapest"},
		{subdomain: "is1-ovpn", region: "Europe", country: "Iceland", city: "Reykjavik"},
		{subdomain: "ie2-ovpn", region: "Europe", country: "Ireland", city: "Dublin"},
		{subdomain: "im1-ovpn", region: "Europe", country: "Isle of Man", city: "Onchan"},
		{subdomain: "vlus-it1-ovpn", region: "Europe", country: "Italy", city: "Milano"},
		{subdomain: "lv1-ovpn", region: "Europe", country: "Latvia", city: "RIGA"},
		{subdomain: "li1-ovpn", region: "Europe", country: "Liechtenstein", city: "Vaduz"},
		{subdomain: "lt1-ovpn", region: "Europe", country: "Lithuania", city: "Vilnius"},
		{subdomain: "lu2-ovpn", region: "Europe", country: "Luxembourg", city: "Luxembourg"},
		{subdomain: "mt1-ovpn", region: "Europe", country: "Malta", city: "Sliema"},
		{subdomain: "mn1-ovpn", region: "Europe", country: "Monaco", city: "Monaco"},
		{subdomain: "vleu-me1-ovpn", region: "Europe", country: "Montenegro", city: "Podgorica"},
		{subdomain: "nl2-ovpn", region: "Europe", country: "Netherlands", city: "Amsterdam"},
		{subdomain: "vleu-no2-ovpn", region: "Europe", country: "Norway", city: "Oslo"},
		{subdomain: "pl2-ovpn", region: "Europe", country: "Poland", city: "Warsaw"},
		{subdomain: "pt2-ovpn", region: "Europe", country: "Portugal", city: "Lisbon"},
		{subdomain: "ro2-ovpn", region: "Europe", country: "Romania", city: "Bucharest"},
		{subdomain: "rs2-ovpn", region: "Europe", country: "Serbia", city: "Niš"},
		{subdomain: "sk1-ovpn", region: "Europe", country: "Slovakia", city: "Bratislava"},
		{subdomain: "si1-ovpn", region: "Europe", country: "Slovenia", city: "Ljubljana"},
		{subdomain: "es-ovpn", region: "Europe", country: "Spain", city: "Barcelona"},
		{subdomain: "vlus-se1-ovpn", region: "Europe", country: "Sweden", city: "Stockholm"},
		{subdomain: "ch2-ovpn", region: "Europe", country: "Switzerland", city: "Zurich"},
		{subdomain: "ukg2-ovpn", region: "Europe", country: "United Kingdom", city: "Gosport"},
		{subdomain: "ukl2-ovpn", region: "Europe", country: "United Kingdom", city: "London"},
		{subdomain: "ukm2-ovpn", region: "Europe", country: "United Kingdom", city: "Maidenhead"},
		{subdomain: "vlus-uk-man1-ovpn", region: "Europe", country: "United Kingdom", city: "Manchester"},
		{subdomain: "bh-ovpn", region: "Middle East", country: "Bahrain", city: "Manama"},
		{subdomain: "vlus-jo1-ovpn", region: "Middle East", country: "Jordan", city: "Amman"},
		{subdomain: "vlus-kw1-ovpn", region: "Middle East", country: "Kuwait", city: "Kuwait"},
		{subdomain: "om1-ovpn", region: "Middle East", country: "Oman", city: "Salalah"},
		{subdomain: "qa1-ovpn", region: "Middle East", country: "Qatar", city: "Doha"},
		{subdomain: "sa1-ovpn", region: "Middle East", country: "Saudi Arabia", city: "Jeddah"},
		{subdomain: "ae2-ovpn", region: "Middle East", country: "United Arab Emirates", city: "Dubai"},
		{subdomain: "aw1-ovpn", region: "North America", country: "Aruba", city: "Oranjestad"},
		{subdomain: "vleu-bb-ovpn", region: "North America", country: "Barbados", city: "Bridgetown"},
		{subdomain: "bz1-ovpn", region: "North America", country: "Belize", city: "Belmopan"},
		{subdomain: "vleu-bm-ovpn", region: "North America", country: "Bermuda", city: "Hamilton"},
		{subdomain: "caq1-ovpn", region: "North America", country: "Canada", city: "Montreal"},
		{subdomain: "cato-ovpn", region: "North America", country: "Canada", city: "Toronto"},
		{subdomain: "cav2-ovpn", region: "North America", country: "Canada", city: "Vancouver"},
		{subdomain: "vleu-ky-ovpn", region: "North America", country: "Cayman Islands", city: "George Town"},
		{subdomain: "vlus-cr1-ovpn", region: "North America", country: "Costa Rica", city: "San Jose"},
		{subdomain: "vleu-dm-ovpn", region: "North America", country: "Dominica", city: "Roseau"},
		{subdomain: "vleu-do-ovpn", region: "North America", country: "Dominican Republic", city: "Santo Domingo"},
		{subdomain: "vleu-sv-ovpn", region: "North America", country: "El Salvador", city: "San Salvador"},
		{subdomain: "vleu-gd-ovpn", region: "North America", country: "Grenada", city: "St George's"},
		{subdomain: "vleu-gt-ovpn", region: "North America", country: "Guatemala", city: "Guatemala"},
		{subdomain: "vleu-ht1-ovpn", region: "North America", country: "Haiti", city: "PORT-AU-PRINCE"},
		{subdomain: "vleu-hn-ovpn", region: "North America", country: "Honduras", city: "TEGUCIGALPA"},
		{subdomain: "jm1-ovpn", region: "North America", country: "Jamaica", city: "Kingston"},
		{subdomain: "vlus-mx2-ovpn", region: "North America", country: "Mexico", city: "Mexico City"},
		{subdomain: "vleu-ms-ovpn", region: "North America", country: "Montserrat", city: "plymouth"},
		{subdomain: "pr1-ovpn", region: "North America", country: "Puerto Rico", city: "San Juan"},
		{subdomain: "vleu-lc-ovpn", region: "North America", country: "Saint Lucia", city: "Castries"},
		{subdomain: "bs1-ovpn", region: "North America", country: "The Bahamas", city: "Freeport"},
		{subdomain: "vleu-tt-ovpn", region: "North America", country: "Trinidad and Tobago", city: "Port of Spain"},
		{subdomain: "vleu-tc-ovpn", region: "North America", country: "Turks and Caicos Islands", city: "Balfour Town"},
		{subdomain: "usva-ovpn", region: "North America", country: "United States", city: "Ashburn"},
		{subdomain: "usil2-ovpn", region: "North America", country: "United States", city: "Chicago"},
		{subdomain: "usoh1-ovpn", region: "North America", country: "United States", city: "Columbus"},
		{subdomain: "usga2-ovpn", region: "North America", country: "United States", city: "Georgia"},
		{subdomain: "ustx2-ovpn", region: "North America", country: "United States", city: "Houston"},
		{subdomain: "usla2-ovpn", region: "North America", country: "United States", city: "Los Angeles"},
		{subdomain: "usfl2-ovpn", region: "North America", country: "United States", city: "Miami"},
		{subdomain: "usnj2-ovpn", region: "North America", country: "United States", city: "New Jersey"},
		{subdomain: "usny2-ovpn", region: "North America", country: "United States", city: "New York"},
		{subdomain: "usphx2-ovpn", region: "North America", country: "United States", city: "Phoenix"},
		{subdomain: "usut2-ovpn", region: "North America", country: "United States", city: "Salt Lake City"},
		{subdomain: "ussf2-ovpn", region: "North America", country: "United States", city: "San Francisco"},
		{subdomain: "ussa-ovpn", region: "North America", country: "United States", city: "Seattle"},
		{subdomain: "uswdc2-ovpn", region: "North America", country: "United States", city: "Washington, D.C."},
		{subdomain: "au-bn-ovpn", region: "Oceania", country: "Australia", city: "Brisbane"},
		{subdomain: "au-me1-ovpn", region: "Oceania", country: "Australia", city: "Melbourne"},
		{subdomain: "au2-pe-ovpn", region: "Oceania", country: "Australia", city: "Perth"},
		{subdomain: "au-sd2-ovpn", region: "Oceania", country: "Australia", city: "Sydney"},
		{subdomain: "nz2-ovpn", region: "Oceania", country: "New Zealand", city: "Auckland"},
		{subdomain: "vlus-ar1-ovpn", region: "South America", country: "Argentina", city: "Buenos Aires"},
		{subdomain: "vleu-bo-ovpn", region: "South America", country: "Bolivia", city: "Sucre"},
		{subdomain: "br2-ovpn", region: "South America", country: "Brazil", city: "Sao Paulo"},
		{subdomain: "vg1-ovpn", region: "South America", country: "British Virgin Island", city: "Road Town"},
		{subdomain: "vlbr-cl-ovpn", region: "South America", country: "Chile", city: "Santiago"},
		{subdomain: "co1-ovpn", region: "South America", country: "Colombia", city: "Bogota"},
		{subdomain: "ec1-ovpn", region: "South America", country: "Ecuador", city: "Quito"},
		{subdomain: "vleu-gy-ovpn", region: "South America", country: "Guyana", city: "Georgetown"},
		{subdomain: "pa2-ovpn", region: "South America", country: "Panama", city: "Panama City"},
		{subdomain: "vleu-py-ovpn", region: "South America", country: "Paraguay", city: "Asuncion"},
		{subdomain: "pe1-ovpn", region: "South America", country: "Peru", city: "Lima"},
		{subdomain: "vleu-sr-ovpn", region: "South America", country: "Suriname", city: "Paramaribo"},
	}
	for i := range servers {
		servers[i].subdomain += "-udp"
	}
	return servers
}
