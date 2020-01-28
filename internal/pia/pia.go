package pia

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/qdm12/golibs/network"
	"github.com/qdm12/golibs/verification"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// Configurator contains methods to download, read and modify the openvpn configuration to connect as a client
type Configurator interface {
	Get(client network.Client,
		encryption constants.PIAEncryption, protocol constants.NetworkProtocol,
		region constants.PIARegion) (lines []string, err error)
	Read(lines []string) (IPs []string, port, device string, err error)
	Modify(lines, IPs []string, port string) (modifiedLines []string, err error)
}

type configurator struct {
	verifier verification.Verifier
}

// NewConfigurator returns a new Configurator object
func NewConfigurator() Configurator {
	return &configurator{verification.NewVerifier()}
}

// Get downloads the PIA client openvpn configuration file for a certain encryption, protocol and region
func (c *configurator) Get(client network.Client, encryption constants.PIAEncryption,
	protocol constants.NetworkProtocol, region constants.PIARegion) (lines []string, err error) {
	URL := c.buildZipURL(encryption, protocol)
	content, status, err := client.GetContent(URL)
	if err != nil {
		return nil, err
	} else if status != 200 {
		return nil, fmt.Errorf("HTTP Get %s resulted in HTTP status code %d", URL, status)
	}
	filename := fmt.Sprintf("%s.ovpn", region)
	fileContent, err := getFileInZip(content, filename)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", URL, err)
	}
	lines = strings.Split(string(fileContent), "\n")
	return lines, nil
}

func (c *configurator) Modify(lines, IPs []string, port string) (modifiedLines []string, err error) {
	// Remove lines
	for _, line := range lines {
		if strings.Contains(line, "privateinternetaccess.com") ||
			strings.Contains(line, "resolve-retry") {
			continue
		}
		modifiedLines = append(modifiedLines, line)
	}
	// Add lines
	for _, IP := range IPs {
		modifiedLines = append(modifiedLines, fmt.Sprintf("remote %s %s", IP, port))
	}
	modifiedLines = append(modifiedLines, "auth-user-pass "+constants.OpenVPNAuthConf)
	modifiedLines = append(modifiedLines, "auth-retry nointeract")
	modifiedLines = append(modifiedLines, "pull-filter ignore \"auth-token\"") // prevent auth failed loops
	modifiedLines = append(modifiedLines, "user nonrootuser")
	modifiedLines = append(modifiedLines, "mute-replay-warnings")
	return modifiedLines, nil
}

func (c *configurator) buildZipURL(encryption constants.PIAEncryption, protocol constants.NetworkProtocol) (URL string) {
	URL = constants.PIAOpenVPNURL + "/openvpn"
	if encryption == constants.PIAEncryptionStrong {
		URL += "-strong"
	}
	if protocol == constants.TCP {
		URL += "-tcp"
	}
	return URL + ".zip"
}

func getFileInZip(zipContent []byte, filename string) (fileContent []byte, err error) {
	contentLength := int64(len(zipContent))
	r, err := zip.NewReader(bytes.NewReader(zipContent), contentLength)
	if err != nil {
		return nil, err
	}
	for _, f := range r.File {
		if f.Name == filename {
			readCloser, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer readCloser.Close()
			return ioutil.ReadAll(readCloser)
		}
	}
	return nil, fmt.Errorf("%s not found in zip archive file", filename)
}

func (c *configurator) Read(lines []string) (IPs []string, port, device string, err error) {
	remoteLineFound := false
	deviceLineFound := false
	for _, line := range lines {
		if strings.HasPrefix(line, "remote ") {
			remoteLineFound = true
			words := strings.Split(line, " ")
			if len(words) != 3 {
				return nil, "", "", fmt.Errorf("line %q misses information", line)
			}
			host := words[1]
			if err := c.verifier.VerifyPort(words[2]); err != nil {
				return nil, "", "", fmt.Errorf("line %q has an invalid port: %w", line, err)
			}
			port = words[2]
			netIPs, err := net.LookupIP(host) // TODO use Unbound
			if err != nil {
				return nil, "", "", err
			}
			for _, netIP := range netIPs {
				IPs = append(IPs, netIP.String())
			}
		} else if strings.HasPrefix(line, "dev ") {
			deviceLineFound = true
			words := strings.Split(line, " ")
			if len(words) != 2 {
				return nil, "", "", fmt.Errorf("line %q misses information", line)
			}
			device = words[1]
			if device != "tun" && device != "tap" {
				return nil, "", "", fmt.Errorf("device %q is not valid", device)
			}
		}
	}
	if remoteLineFound && deviceLineFound {
		return IPs, port, device, nil
	} else if !remoteLineFound {
		return nil, "", "", fmt.Errorf("remote line not found in Openvpn configuration")
	} else {
		return nil, "", "", fmt.Errorf("device line not found in Openvpn configuration")
	}
}
