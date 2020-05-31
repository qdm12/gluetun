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
	provider := flag.String("provider", "pia", "VPN provider to resolve for, 'pia' or 'windscribe'")
	region := flag.String("region", "all", "Comma separated list of VPN provider region names to resolve for, use 'all' to resolve all")
	flag.Parse()

	resolver := newResolver(*resolverAddress)
	lookupIP := newLookupIP(resolver)

	var domain, template string
	var servers []server
	switch *provider {
	case "pia":
		domain = "privateinternetaccess.com"
		template = "{Region: models.PIARegion(%q), IPs: []net.IP{%s}},"
		servers = piaServers()
	case "windscribe":
		domain = "windscribe.com"
		template = "{Region: models.WindscribeRegion(%q), IPs: []net.IP{%s}},"
		servers = windscribeServers()
	case "surfshark":
		domain = "prod.surfshark.com"
		template = "{Region: models.SurfsharkRegion(%q), IPs: []net.IP{%s}},"
		servers = surfsharkServers()
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
	for _, s := range servers {
		s := s
		go func() {
			ips, err := resolveRepeat(ctx, lookupIP, s.subdomain+"."+domain, 3)
			if err != nil {
				errorChannel <- err
				return
			}
			ipsString := formatIPs(ips)
			stringChannel <- fmt.Sprintf(template, s.region, ipsString)
		}()
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
	ipsChannel := make(chan []net.IP)
	errorsChannel := make(chan error)
	for i := 0; i < n; i++ {
		go func() {
			ips, err := lookupIP(ctx, host)
			if err != nil {
				errorsChannel <- err
			} else {
				ipsChannel <- ips
			}
		}()
	}
	for i := 0; i < n; i++ {
		select {
		case err = <-errorsChannel:
		case newIPs := <-ipsChannel:
			ips = append(ips, newIPs...)
		}
	}
	close(errorsChannel)
	close(ipsChannel)
	if err != nil {
		return nil, err
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

func formatIPs(ips []net.IP) (s string) {
	ipStrings := make([]string, len(ips))
	for i := range ips {
		ipStrings[i] = fmt.Sprintf("{%s}", strings.ReplaceAll(ips[i].String(), ".", ", "))
	}
	return strings.Join(ipStrings, ", ")
}

type server struct {
	subdomain string
	region    string
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
		{subdomain: "ba-sjj", region: "Bosnia and Herzegovina "},
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
