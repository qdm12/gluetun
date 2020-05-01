package constants

import (
	"fmt"
	"net"
	"strings"
	"testing"
)

func Test_LookupPIAServers(t *testing.T) {
	t.SkipNow()
	subdomains := []string{
		"au-melbourne",
		"au-perth",
		"au-sydney",
		"austria",
		"belgium",
		"ca-montreal",
		"ca-toronto",
		"ca-vancouver",
		"czech",
		"de-berlin",
		"de-frankfurt",
		"denmark",
		"fi",
		"france",
		"hk",
		"hungary",
		"in",
		"ireland",
		"israel",
		"italy",
		"japan",
		"lu",
		"mexico",
		"nl",
		"nz",
		"no",
		"poland",
		"ro",
		"sg",
		"spain",
		"sweden",
		"swiss",
		"ae",
		"uk-london",
		"uk-manchester",
		"uk-southampton",
		"us-atlanta",
		"us-california",
		"us-chicago",
		"us-denver",
		"us-east",
		"us-florida",
		"us-houston",
		"us-lasvegas",
		"us-newyorkcity",
		"us-seattle",
		"us-siliconvalley",
		"us-texas",
		"us-washingtondc",
		"us-west",
	}
	for _, subdomain := range subdomains {
		ips, err := net.LookupIP(subdomain + ".privateinternetaccess.com")
		if err != nil {
			t.Log(err)
			continue
		}
		s := make([]string, len(ips))
		for i := range ips {
			s[i] = fmt.Sprintf("{%s}", strings.ReplaceAll(ips[i].String(), ".", ", "))
		}
		t.Logf("%s: %s", subdomain, strings.Join(s, ", "))
	}
}

func Test_LookupWindscribeServers(t *testing.T) {
	t.SkipNow()
	subdomains := []string{
		"al",
		"ar",
		"ar",
		"au",
		"at",
		"az",
		"be",
		"ba",
		"br",
		"bg",
		"ca",
		"ca-west",
		"co",
		"hr",
		"cy",
		"cz",
		"dk",
		"ee",
		"aq",
		"fi",
		"fr",
		"ge",
		"de",
		"gr",
		"hk",
		"hu",
		"is",
		"in",
		"id",
		"ie",
		"il",
		"it",
		"jp",
		"lv",
		"lt",
		"mk",
		"my",
		"mx",
		"md",
		"nl",
		"nz",
		"no",
		"ph",
		"pl",
		"pt",
		"ro",
		"ru",
		"rs",
		"sg",
		"sk",
		"si",
		"za",
		"kr",
		"es",
		"se",
		"ch",
		"th",
		"tn",
		"tr",
		"ua",
		"ae",
		"uk",
		"us-central",
		"us-east",
		"us-west",
		"vn",
		"wf-ca",
		"wf-jp",
		"wf-uk",
		"wf-us",
	}
	for _, subdomain := range subdomains {
		ips, err := net.LookupIP(subdomain + ".windscribe.com")
		if err != nil {
			t.Log(err)
			continue
		}
		s := make([]string, len(ips))
		for i := range ips {
			s[i] = fmt.Sprintf("{%s}", strings.ReplaceAll(ips[i].String(), ".", ", "))
		}
		t.Logf("%s: %s", subdomain, strings.Join(s, ", "))
	}
}
