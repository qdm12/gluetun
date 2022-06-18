package updater

import (
	"context"
	"net/http"
	"strings"

	"github.com/qdm12/gluetun/internal/updater/openvpn"
)

func getAllHostToURL(ctx context.Context, client *http.Client) (
	tcpHostToURL, udpHostToURL map[string]string, err error) {
	tcpHostToURL, err = getHostToURL(ctx, client, "TCP")
	if err != nil {
		return nil, nil, err
	}

	udpHostToURL, err = getHostToURL(ctx, client, "UDP")
	if err != nil {
		return nil, nil, err
	}

	return tcpHostToURL, udpHostToURL, nil
}

func getHostToURL(ctx context.Context, client *http.Client, protocol string) (
	hostToURL map[string]string, err error) {
	const baseURL = "https://vpn.hidemyass.com/vpn-config"
	indexURL := baseURL + "/" + strings.ToUpper(protocol) + "/"

	urls, err := fetchIndex(ctx, client, indexURL)
	if err != nil {
		return nil, err
	}

	const failEarly = true
	hostToURL, errors := openvpn.FetchMultiFiles(ctx, client, urls, failEarly)
	if len(errors) > 0 {
		return nil, errors[0]
	}

	return hostToURL, nil
}
