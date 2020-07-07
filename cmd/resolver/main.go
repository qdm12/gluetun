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
	provider := flag.String("provider", "pia", "VPN provider to resolve for, 'pia', 'windscribe' or 'cyberghost'")
	region := flag.String("region", "all", "Comma separated list of VPN provider region names to resolve for, use 'all' to resolve all")
	flag.Parse()

	resolver := newResolver(*resolverAddress)
	lookupIP := newLookupIP(resolver)

	var domain string
	var servers []server
	switch *provider {
	case "pia":
		domain = "privateinternetaccess.com"
		servers = piaServers()
	case "windscribe":
		domain = "windscribe.com"
		servers = windscribeServers()
	case "surfshark":
		domain = "prod.surfshark.com"
		servers = surfsharkServers()
	case "cyberghost":
		domain = "cg-dialup.net"
		servers = cyberghostServers()
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
			"{Region: models.PIARegion(%q), IPs: []net.IP{%s}},",
			s.region, ipString,
		)
	case "windscribe":
		return fmt.Sprintf(
			"{Region: models.WindscribeRegion(%q), IPs: []net.IP{%s}},",
			s.region, ipString,
		)
	case "surfshark":
		return fmt.Sprintf(
			"{Region: models.SurfsharkRegion(%q), IPs: []net.IP{%s}},",
			s.region, ipString,
		)
	case "cyberghost":
		return fmt.Sprintf(
			"{Region: models.CyberghostRegion(%q), Group: models.CyberghostGroup(%q), IPs: []net.IP{%s}},",
			s.region, s.group, ipString,
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
}

func piaServers() []server {
	return []server{
		{subdomain: "au-melbourne", region: "AU Melbourne"},
		{subdomain: "au-perth", region: "AU Perth"},
		{subdomain: "au-sydney", region: "AU Sydney"},
		{subdomain: "austria", region: "Austria"},
		{subdomain: "belgium", region: "Belgium"},
		{subdomain: "ca-montreal", region: "CA Montreal"},
		{subdomain: "ca-toronto", region: "CA Toronto"},
		{subdomain: "ca-vancouver", region: "CA Vancouver"},
		{subdomain: "czech", region: "Czech Republic"},
		{subdomain: "de-berlin", region: "DE Berlin"},
		{subdomain: "de-frankfurt", region: "DE Frankfurt"},
		{subdomain: "denmark", region: "Denmark"},
		{subdomain: "fi", region: "Finlan"},
		{subdomain: "france", region: "France"},
		{subdomain: "hk", region: "Hong Kong"},
		{subdomain: "hungary", region: "Hungary"},
		{subdomain: "in", region: "India"},
		{subdomain: "ireland", region: "Ireland"},
		{subdomain: "israel", region: "Israel"},
		{subdomain: "italy", region: "Italy"},
		{subdomain: "japan", region: "Japan"},
		{subdomain: "lu", region: "Luxembourg"},
		{subdomain: "mexico", region: "Mexico"},
		{subdomain: "nl", region: "Netherlands"},
		{subdomain: "nz", region: "New Zealand"},
		{subdomain: "no", region: "Norway"},
		{subdomain: "poland", region: "Poland"},
		{subdomain: "ro", region: "Romania"},
		{subdomain: "sg", region: "Singapore"},
		{subdomain: "spain", region: "Spain"},
		{subdomain: "sweden", region: "Sweden"},
		{subdomain: "swiss", region: "Switzerland"},
		{subdomain: "ae", region: "UAE"},
		{subdomain: "uk-london", region: "UK London"},
		{subdomain: "uk-manchester", region: "UK Manchester"},
		{subdomain: "uk-southampton", region: "UK Southampton"},
		{subdomain: "us-atlanta", region: "US Atlanta"},
		{subdomain: "us-california", region: "US California"},
		{subdomain: "us-chicago", region: "US Chicago"},
		{subdomain: "us-denver", region: "US Denver"},
		{subdomain: "us-east", region: "US East"},
		{subdomain: "us-florida", region: "US Florida"},
		{subdomain: "us-houston", region: "US Houston"},
		{subdomain: "us-lasvegas", region: "US Las Vegas"},
		{subdomain: "us-newyorkcity", region: "US New York City"},
		{subdomain: "us-seattle", region: "US Seattle"},
		{subdomain: "us-siliconvalley", region: "US Silicon Valley"},
		{subdomain: "us-texas", region: "US Texas"},
		{subdomain: "us-washingtondc", region: "US Washington DC"},
		{subdomain: "us-west", region: "US West"},
	}
}

func windscribeServers() []server {
	return []server{
		{subdomain: "al", region: "albania"},
		{subdomain: "ar", region: "argentina"},
		{subdomain: "au", region: "australia"},
		{subdomain: "at", region: "austria"},
		{subdomain: "az", region: "azerbaijan"},
		{subdomain: "be", region: "belgium"},
		{subdomain: "ba", region: "bosnia"},
		{subdomain: "br", region: "brazil"},
		{subdomain: "bg", region: "bulgaria"},
		{subdomain: "ca", region: "canada east"},
		{subdomain: "ca-west", region: "canada west"},
		{subdomain: "co", region: "colombia"},
		{subdomain: "hr", region: "croatia"},
		{subdomain: "cy", region: "cyprus"},
		{subdomain: "cz", region: "czech republic"},
		{subdomain: "dk", region: "denmark"},
		{subdomain: "ee", region: "estonia"},
		{subdomain: "aq", region: "fake antarctica"},
		{subdomain: "fi", region: "finland"},
		{subdomain: "fr", region: "france"},
		{subdomain: "ge", region: "georgia"},
		{subdomain: "de", region: "germany"},
		{subdomain: "gr", region: "greece"},
		{subdomain: "hk", region: "hong kong"},
		{subdomain: "hu", region: "hungary"},
		{subdomain: "is", region: "iceland"},
		{subdomain: "in", region: "india"},
		{subdomain: "id", region: "indonesia"},
		{subdomain: "ie", region: "ireland"},
		{subdomain: "il", region: "israel"},
		{subdomain: "it", region: "italy"},
		{subdomain: "jp", region: "japan"},
		{subdomain: "lv", region: "latvia"},
		{subdomain: "lt", region: "lithuania"},
		{subdomain: "mk", region: "macedonia"},
		{subdomain: "my", region: "malaysia"},
		{subdomain: "mx", region: "mexico"},
		{subdomain: "md", region: "moldova"},
		{subdomain: "nl", region: "netherlands"},
		{subdomain: "nz", region: "new zealand"},
		{subdomain: "no", region: "norway"},
		{subdomain: "ph", region: "philippines"},
		{subdomain: "pl", region: "poland"},
		{subdomain: "pt", region: "portugal"},
		{subdomain: "ro", region: "romania"},
		{subdomain: "ru", region: "russia"},
		{subdomain: "rs", region: "serbia"},
		{subdomain: "sg", region: "singapore"},
		{subdomain: "sk", region: "slovakia"},
		{subdomain: "si", region: "slovenia"},
		{subdomain: "za", region: "south africa"},
		{subdomain: "kr", region: "south korea"},
		{subdomain: "es", region: "spain"},
		{subdomain: "se", region: "sweden"},
		{subdomain: "ch", region: "switzerland"},
		{subdomain: "th", region: "thailand"},
		{subdomain: "tn", region: "tunisia"},
		{subdomain: "tr", region: "turkey"},
		{subdomain: "ua", region: "ukraine"},
		{subdomain: "ae", region: "united arab emirates"},
		{subdomain: "uk", region: "united kingdom"},
		{subdomain: "us-central", region: "us central"},
		{subdomain: "us-east", region: "us east"},
		{subdomain: "us-west", region: "us west"},
		{subdomain: "vn", region: "vietnam"},
		{subdomain: "wf-ca", region: "windflix ca"},
		{subdomain: "wf-jp", region: "windflix jp"},
		{subdomain: "wf-uk", region: "windflix uk"},
		{subdomain: "wf-us", region: "windflix us"},
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
		{subdomain: "it-rom", region: "italy Rome"},
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
