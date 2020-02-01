package pia

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

func (c *configurator) DownloadOvpnConfig(encryption models.PIAEncryption,
	protocol models.NetworkProtocol, region models.PIARegion) (lines []string, err error) {
	URL := buildZipURL(encryption, protocol)
	content, status, err := c.client.GetContent(URL)
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

func buildZipURL(encryption models.PIAEncryption, protocol models.NetworkProtocol) (URL string) {
	URL = string(constants.PIAOpenVPNURL) + "/openvpn"
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
