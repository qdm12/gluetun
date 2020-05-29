package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/qdm12/golibs/network"
)

func main() {
	os.Exit(_main())
}

func _main() int {
	provider := flag.String("provider", "surfshark", "VPN provider to parse openvpn files for, can be 'surfshark'")
	flag.Parse()

	var url string
	switch *provider {
	case "surfshark":
		url = "https://account.surfshark.com/api/v1/server/configurations"
	default:
		fmt.Printf("Provider %q is not supported\n", *provider)
		return 1
	}
	client := network.NewClient(time.Second)
	zipBytes, status, err := client.GetContent(url)
	if err != nil {
		fmt.Println(err)
		return 1
	} else if status != http.StatusOK {
		fmt.Printf("Getting %s results in HTTP status code %d\n", url, status)
		return 1
	}
	contents, err := zipExtractAll(zipBytes)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	uniqueSubdomains := make(map[string]struct{})
	for _, content := range contents {
		subdomain, err := extractInformation(content)
		if err != nil {
			fmt.Println(err)
			return 1
		}
		uniqueSubdomains[subdomain] = struct{}{}
	}
	subdomains := make([]string, len(uniqueSubdomains))
	i := 0
	for subdomain := range uniqueSubdomains {
		subdomains[i] = subdomain
		i++
	}
	sort.Slice(subdomains, func(i, j int) bool {
		return subdomains[i] < subdomains[j]
	})
	fmt.Println("Subdomains found are: ", strings.Join(subdomains, ","))
	return 0
}

func zipExtractAll(zipBytes []byte) (contents [][]byte, err error) {
	r, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return nil, err
	}
	contents = make([][]byte, len(r.File))
	for i, zf := range r.File {
		f, err := zf.Open()
		if err != nil {
			return nil, err
		}
		defer f.Close()
		contents[i], err = ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		if err := f.Close(); err != nil {
			return nil, err
		}
	}
	return contents, nil
}

func extractInformation(content []byte) (subdomain string, err error) {
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "remote ") {
			words := strings.Fields(line)
			if len(words) < 2 {
				return "", fmt.Errorf("not enough words on line %q", line)
			}
			host := words[1]
			return strings.TrimSuffix(host, ".prod.surfshark.com"), nil
		}
	}
	return "", fmt.Errorf("could not find remote line")
}
