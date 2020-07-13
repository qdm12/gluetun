package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/qdm12/golibs/network"
)

func main() {
	os.Exit(_main())
}

// Find subdomains from .ovpn files contained in a .zip file
func _main() int {
	provider := flag.String("provider", "surfshark", "VPN provider to parse openvpn files for, can be 'surfshark' or 'vyprvpn")
	flag.Parse()

	var urls []string
	var suffix string
	switch *provider {
	case "surfshark":
		urls = []string{
			"https://account.surfshark.com/api/v1/server/configurations",
			"https://v2uploads.zopim.io/p/2/L/p2LbwLkvfQoSdzOl6VEltzQA6StiZqrs/12500634259669c77012765139bcfe4f4c90db1e.zip",
		}
		suffix = ".prod.surfshark.com"
	case "vyprvpn":
		urls = []string{
			"https://support.vyprvpn.com/hc/article_attachments/360052617332/Vypr_OpenVPN_20200320.zip",
		}
		suffix = ".vyprvpn.com"
	default:
		fmt.Printf("Provider %q is not supported\n", *provider)
		return 1
	}
	contents, err := fetchAndExtractFiles(urls...)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	uniqueSubdomainsToFilename := make(map[string]string)
	for fileName, content := range contents {
		subdomain, err := extractInformation(content, suffix)
		if err != nil {
			fmt.Println(err)
			return 1
		} else if len(subdomain) > 0 {
			fileName = strings.TrimSuffix(fileName, ".ovpn")
			fileName = strings.ReplaceAll(fileName, " - ", " ")
			uniqueSubdomainsToFilename[subdomain] = fileName
		}
	}
	type subdomainFilename struct {
		subdomain string
		fileName  string
	}
	subdomains := make([]subdomainFilename, len(uniqueSubdomainsToFilename))
	i := 0
	for subdomain, fileName := range uniqueSubdomainsToFilename {
		subdomains[i] = subdomainFilename{
			subdomain: subdomain,
			fileName:  fileName,
		}
		i++
	}
	sort.Slice(subdomains, func(i, j int) bool {
		return subdomains[i].subdomain < subdomains[j].subdomain
	})
	fmt.Println("Subdomain Filename")
	for i := range subdomains {
		fmt.Printf("%s %s\n", subdomains[i].subdomain, subdomains[i].fileName)
	}
	return 0
}

func fetchAndExtractFiles(urls ...string) (contents map[string][]byte, err error) {
	client := network.NewClient(10 * time.Second)
	contents = make(map[string][]byte)
	for _, url := range urls {
		zipBytes, status, err := client.GetContent(url)
		if err != nil {
			return nil, err
		} else if status != http.StatusOK {
			return nil, fmt.Errorf("Getting %s results in HTTP status code %d", url, status)
		}
		newContents, err := zipExtractAll(zipBytes)
		if err != nil {
			return nil, err
		}
		for fileName, content := range newContents {
			contents[fileName] = content
		}
	}
	return contents, nil
}

func zipExtractAll(zipBytes []byte) (contents map[string][]byte, err error) {
	r, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return nil, err
	}
	contents = map[string][]byte{}
	for _, zf := range r.File {
		fileName := filepath.Base(zf.Name)
		if !strings.HasSuffix(fileName, ".ovpn") {
			continue
		}
		f, err := zf.Open()
		if err != nil {
			return nil, err
		}
		defer f.Close()
		contents[fileName], err = ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		if err := f.Close(); err != nil {
			return nil, err
		}
	}
	return contents, nil
}

func extractInformation(content []byte, suffix string) (subdomain string, err error) {
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "remote ") {
			words := strings.Fields(line)
			if len(words) < 2 {
				return "", fmt.Errorf("not enough words on line %q", line)
			}
			host := words[1]
			if net.ParseIP(host) != nil {
				return "", nil // ignore IP addresses
			}
			return strings.TrimSuffix(host, suffix), nil
		}
	}
	return "", fmt.Errorf("could not find remote line in: %s", string(content))
}
