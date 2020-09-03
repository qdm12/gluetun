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
	case "windscribe":
		domain = "windscribe.com"
		servers = windscribeServers()
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
	case "windscribe":
		return fmt.Sprintf(
			"{Region: %q, IPs: []net.IP{%s}},",
			s.region, ipString,
		)
	case "cyberghost":
		return fmt.Sprintf(
			"{Region: %q, Group: %q, IPs: []net.IP{%s}},",
			s.region, s.group, ipString,
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
