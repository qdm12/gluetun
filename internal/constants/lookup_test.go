package constants

import (
	"bytes"
	"fmt"
	"net"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_LookupPIAServers(t *testing.T) { //nolint:gocognit
	t.SkipNow()
	servers := []struct {
		subdomain string
		region    string
	}{
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
	for _, server := range servers {
		s, err := getIPsStringFromSubdomain(server.subdomain + ".privateinternetaccess.com")
		require.NoError(t, err)
		t.Logf("{Region: models.PIARegion(%q), IPs: []net.IP{%s}},", server.region, s)
	}
}

func Test_LookupWindscribeServers(t *testing.T) {
	t.SkipNow()
	servers := []struct {
		subdomain string
		region    string
	}{
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
	for _, server := range servers {
		s, err := getIPsStringFromSubdomain(server.subdomain + ".windscribe.com")
		require.NoError(t, err)
		t.Logf("{Region: models.WindscribeRegion(%q), IPs: []net.IP{%s}},", server.region, s)
	}
}

func getIPsStringFromSubdomain(subdomain string) (s string, err error) { //nolint:unused
	const tries = 3
	ipsChannel := make(chan []net.IP)
	errorsChannel := make(chan error)
	for i := 0; i < tries; i++ {
		go func() {
			ips, err := net.LookupIP(subdomain)
			if err != nil {
				errorsChannel <- err
			} else {
				ipsChannel <- ips
			}
		}()
	}
	var ips []net.IP
	for i := 0; i < tries; i++ {
		select {
		case err = <-errorsChannel:
		case newIPs := <-ipsChannel:
			ips = append(ips, newIPs...)
		}
	}
	close(errorsChannel)
	close(ipsChannel)
	if err != nil {
		return "", err
	}
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
	ipStrings := make([]string, len(ips))
	for i := range ips {
		ipStrings[i] = fmt.Sprintf("{%s}", strings.ReplaceAll(ips[i].String(), ".", ", "))
	}
	return strings.Join(ipStrings, ", "), nil
}
