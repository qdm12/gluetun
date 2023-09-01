package privateinternetaccess

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
	"github.com/qdm12/gluetun/internal/wireguard"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strconv"
)

func (p *Provider) GetWireguardConnection(ctx context.Context, connection models.Connection, wireguardSettings settings.Wireguard, ipv6Supported bool) (settings wireguard.Settings, err error) {
	// create http client to make requests to the server's API
	client, err := newHTTPClient(connection.Hostname)

	// generate a new private key
	privateKey, err := wgtypes.GeneratePrivateKey()

	// generate public key from private key
	publicKey := privateKey.PublicKey()

	// fetch token from PIA's API
	token, err := fetchToken(ctx, client, "gtoken", p.authFilePath)

	gateway := connection.IP.String()

	// error substitutions
	errSubstitutions := map[string]string{url.QueryEscape(token): "<token>"}

	// let's register the public key with the server
	// this will also give us the server's information
	queryParams := make(url.Values)
	queryParams.Add("pt", token)
	queryParams.Add("pubkey", publicKey.String())
	addKeyUrl := url.URL{
		Scheme:   "https",
		Host:     net.JoinHostPort(gateway, "1337"),
		Path:     "/addKey",
		RawQuery: queryParams.Encode(),
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, addKeyUrl.String(), nil)
	if err != nil {
		return settings, ReplaceInErr(err, errSubstitutions)
	}

	response, err := client.Do(request)
	if err != nil {
		return settings, ReplaceInErr(err, errSubstitutions)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return settings, makeNOKStatusError(response, errSubstitutions)
	}

	// decode the response
	decoder := json.NewDecoder(response.Body)
	var responseData struct {
		Status     string `json:"status"`
		Message    string `json:"message"`
		PeerIp     string `json:"peer_ip"`
		ServerKey  string `json:"server_key"`
		ServerPort string `json:"server_port"`
	}

	if err := decoder.Decode(&responseData); err != nil {
		return settings, err
	}

	if responseData.Status != "OK" {
		return settings, fmt.Errorf("%w: %s: %s", ErrBadResponse, responseData.Status, responseData.Message)
	}

	connection.IP, err = netip.ParseAddr(responseData.PeerIp)
	if err != nil {
		return settings, err
	}

	port, err := strconv.ParseUint(responseData.ServerPort, 10, 16)
	if err != nil {
		return settings, err
	}
	connection.Port = uint16(port)

	psk := publicKey.String()
	wireguardSettings.PreSharedKey = &psk
	pk := privateKey.String()
	wireguardSettings.PrivateKey = &pk

	// create the wireguard settings
	settings = utils.BuildWireguardSettings(connection, wireguardSettings, ipv6Supported)

	return settings, nil
}
