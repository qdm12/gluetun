package provider

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	gluetunLog "github.com/qdm12/gluetun/internal/logging"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
)

type pia struct {
	servers      []models.PIAServer
	timeNow      timeNowFunc
	randSource   rand.Source
	activeServer models.PIAServer
}

func newPrivateInternetAccess(servers []models.PIAServer, timeNow timeNowFunc) *pia {
	return &pia{
		servers:    servers,
		timeNow:    timeNow,
		randSource: rand.NewSource(timeNow().UnixNano()),
	}
}

func (p *pia) GetOpenVPNConnection(selection models.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	var port uint16
	switch selection.Protocol {
	case constants.TCP:
		switch selection.EncryptionPreset {
		case constants.PIAEncryptionPresetNormal:
			port = 502
		case constants.PIAEncryptionPresetStrong:
			port = 501
		}
	case constants.UDP:
		switch selection.EncryptionPreset {
		case constants.PIAEncryptionPresetNormal:
			port = 1198
		case constants.PIAEncryptionPresetStrong:
			port = 1197
		}
	}
	if port == 0 {
		return connection, fmt.Errorf(
			"combination of protocol %q and encryption %q does not yield any port number",
			selection.Protocol, selection.EncryptionPreset)
	}

	servers := p.servers
	if selection.TargetIP != nil {
		connection = models.OpenVPNConnection{IP: selection.TargetIP, Port: port, Protocol: selection.Protocol}
	} else {
		servers := filterPIAServers(servers, selection.Regions, selection.Protocol)
		if len(servers) == 0 {
			return connection, fmt.Errorf("no server found for region %s and protocol %s",
				commaJoin(selection.Regions), selection.Protocol)
		}

		var connections []models.OpenVPNConnection
		for _, server := range servers {
			connection := models.OpenVPNConnection{IP: server.IP, Port: port, Protocol: selection.Protocol}
			connections = append(connections, connection)
		}

		connection = pickRandomConnection(connections, p.randSource)
	}

	// Reverse lookup server from picked connection
	for _, server := range servers {
		if connection.IP.Equal(server.IP) {
			p.activeServer = server
			break
		}
	}

	return connection, nil
}

func (p *pia) BuildConf(connection models.OpenVPNConnection,
	username string, settings settings.OpenVPN) (lines []string) {
	var X509CRL, certificate string
	var defaultCipher, defaultAuth string
	if settings.Provider.ExtraConfigOptions.EncryptionPreset == constants.PIAEncryptionPresetNormal {
		defaultCipher = "aes-128-cbc"
		defaultAuth = "sha1"
		X509CRL = constants.PiaX509CRLNormal
		certificate = constants.PIACertificateNormal
	} else { // strong encryption
		defaultCipher = aes256cbc
		defaultAuth = "sha256"
		X509CRL = constants.PiaX509CRLStrong
		certificate = constants.PIACertificateStrong
	}
	if len(settings.Cipher) == 0 {
		settings.Cipher = defaultCipher
	}
	if len(settings.Auth) == 0 {
		settings.Auth = defaultAuth
	}
	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"remote-cert-tls server",
		"ping 10",
		"ping-exit 60",
		"ping-timer-rem",
		"tls-exit",

		// PIA specific
		"reneg-sec 0",
		"compress", // allow PIA server to choose the compression to use

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"suppress-timestamps",

		// Modified variables
		fmt.Sprintf("verb %d", settings.Verbosity),
		fmt.Sprintf("auth-user-pass %s", constants.OpenVPNAuthConf),
		fmt.Sprintf("proto %s", connection.Protocol),
		fmt.Sprintf("remote %s %d", connection.IP, connection.Port),
		fmt.Sprintf("cipher %s", settings.Cipher),
		fmt.Sprintf("auth %s", settings.Auth),
	}
	if strings.HasSuffix(settings.Cipher, "-gcm") {
		lines = append(lines, "ncp-disable")
	}
	if !settings.Root {
		lines = append(lines, "user "+username)
	}
	if settings.MSSFix > 0 {
		lines = append(lines, "mssfix "+strconv.Itoa(int(settings.MSSFix)))
	}
	lines = append(lines, []string{
		"<crl-verify>",
		"-----BEGIN X509 CRL-----",
		X509CRL,
		"-----END X509 CRL-----",
		"</crl-verify>",
	}...)
	lines = append(lines, []string{
		"<ca>",
		"-----BEGIN CERTIFICATE-----",
		certificate,
		"-----END CERTIFICATE-----",
		"</ca>",
		"",
	}...)
	return lines
}

//nolint:gocognit
func (p *pia) PortForward(ctx context.Context, client *http.Client,
	openFile os.OpenFileFunc, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath models.Filepath)) {
	commonName := p.activeServer.ServerName
	if !p.activeServer.PortForward {
		pfLogger.Error("The server %s (region %s) does not support port forwarding",
			commonName, p.activeServer.Region)
		return
	}
	if gateway == nil {
		pfLogger.Error("aborting because: VPN gateway IP address was not found")
		return
	}

	client, err := newPIAHTTPClient(commonName)
	if err != nil {
		pfLogger.Error("aborting because: %s", err)
		return
	}
	defer pfLogger.Warn("loop exited")
	data, err := readPIAPortForwardData(openFile)
	if err != nil {
		pfLogger.Error(err)
	}
	dataFound := data.Port > 0
	durationToExpiration := data.Expiration.Sub(p.timeNow())
	expired := durationToExpiration <= 0

	if dataFound {
		pfLogger.Info("Found persistent forwarded port data for port %d", data.Port)
		if expired {
			pfLogger.Warn("Forwarded port data expired on %s, getting another one", data.Expiration.Format(time.RFC1123))
		} else {
			pfLogger.Info("Forwarded port data expires in %s", gluetunLog.FormatDuration(durationToExpiration))
		}
	}

	if !dataFound || expired {
		tryUntilSuccessful(ctx, pfLogger, func() error {
			data, err = refreshPIAPortForwardData(ctx, client, gateway, openFile)
			return err
		})
		if ctx.Err() != nil {
			return
		}
		durationToExpiration = data.Expiration.Sub(p.timeNow())
	}
	pfLogger.Info("Port forwarded is %d expiring in %s", data.Port, gluetunLog.FormatDuration(durationToExpiration))

	// First time binding
	tryUntilSuccessful(ctx, pfLogger, func() error {
		return bindPIAPort(ctx, client, gateway, data)
	})
	if ctx.Err() != nil {
		return
	}

	filepath := string(syncState(data.Port))
	pfLogger.Info("Writing port to %s", filepath)
	if err := writePortForwardedToFile(openFile, filepath, data.Port); err != nil {
		pfLogger.Error(err)
	}

	if err := fw.SetAllowedPort(ctx, data.Port, string(constants.TUN)); err != nil {
		pfLogger.Error(err)
	}

	expiryTimer := time.NewTimer(durationToExpiration)
	const keepAlivePeriod = 15 * time.Minute
	// Timer behaving as a ticker
	keepAliveTimer := time.NewTimer(keepAlivePeriod)
	for {
		select {
		case <-ctx.Done():
			removeCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			if err := fw.RemoveAllowedPort(removeCtx, data.Port); err != nil {
				pfLogger.Error(err)
			}
			if !keepAliveTimer.Stop() {
				<-keepAliveTimer.C
			}
			if !expiryTimer.Stop() {
				<-expiryTimer.C
			}
			return
		case <-keepAliveTimer.C:
			if err := bindPIAPort(ctx, client, gateway, data); err != nil {
				pfLogger.Error(err)
			}
			keepAliveTimer.Reset(keepAlivePeriod)
		case <-expiryTimer.C:
			pfLogger.Warn("Forward port has expired on %s, getting another one", data.Expiration.Format(time.RFC1123))
			oldPort := data.Port
			for {
				data, err = refreshPIAPortForwardData(ctx, client, gateway, openFile)
				if err != nil {
					pfLogger.Error(err)
					continue
				}
				break
			}
			durationToExpiration := data.Expiration.Sub(p.timeNow())
			pfLogger.Info("Port forwarded is %d expiring in %s", data.Port, gluetunLog.FormatDuration(durationToExpiration))
			if err := fw.RemoveAllowedPort(ctx, oldPort); err != nil {
				pfLogger.Error(err)
			}
			if err := fw.SetAllowedPort(ctx, data.Port, string(constants.TUN)); err != nil {
				pfLogger.Error(err)
			}
			filepath := syncState(data.Port)
			pfLogger.Info("Writing port to %s", filepath)
			if err := writePortForwardedToFile(openFile, string(filepath), data.Port); err != nil {
				pfLogger.Error(err)
			}
			if err := bindPIAPort(ctx, client, gateway, data); err != nil {
				pfLogger.Error(err)
			}
			if !keepAliveTimer.Stop() {
				<-keepAliveTimer.C
			}
			keepAliveTimer.Reset(keepAlivePeriod)
			expiryTimer.Reset(durationToExpiration)
		}
	}
}

func filterPIAServers(servers []models.PIAServer, regions []string, protocol models.NetworkProtocol) (
	filtered []models.PIAServer) {
	for _, server := range servers {
		switch {
		case filterByPossibilities(server.Region, regions):
		case server.Protocol != protocol:
		default:
			filtered = append(filtered, server)
		}
	}
	return filtered
}

func newPIAHTTPClient(serverName string) (client *http.Client, err error) {
	certificateBytes, err := base64.StdEncoding.DecodeString(constants.PIACertificateStrong)
	if err != nil {
		return nil, fmt.Errorf("cannot decode PIA root certificate: %w", err)
	}
	certificate, err := x509.ParseCertificate(certificateBytes)
	if err != nil {
		return nil, fmt.Errorf("cannot parse PIA root certificate: %w", err)
	}
	// certificate.DNSNames = []string{serverName, "10.0.0.1"}
	rootCAs := x509.NewCertPool()
	rootCAs.AddCert(certificate)
	TLSClientConfig := &tls.Config{
		RootCAs:    rootCAs,
		MinVersion: tls.VersionTLS12,
		ServerName: serverName,
	}
	//nolint:gomnd
	transport := http.Transport{
		TLSClientConfig: TLSClientConfig,
		Proxy:           http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	const httpTimeout = 30 * time.Second
	client = &http.Client{Transport: &transport, Timeout: httpTimeout}
	return client, nil
}

func refreshPIAPortForwardData(ctx context.Context, client *http.Client,
	gateway net.IP, openFile os.OpenFileFunc) (data piaPortForwardData, err error) {
	data.Token, err = fetchPIAToken(ctx, openFile, client)
	if err != nil {
		return data, fmt.Errorf("cannot obtain token: %w", err)
	}
	data.Port, data.Signature, data.Expiration, err = fetchPIAPortForwardData(ctx, client, gateway, data.Token)
	if err != nil {
		return data, fmt.Errorf("cannot obtain port forwarding data: %w", err)
	}
	if err := writePIAPortForwardData(openFile, data); err != nil {
		return data, fmt.Errorf("cannot persist port forwarding information to file: %w", err)
	}
	return data, nil
}

type piaPayload struct {
	Token      string    `json:"token"`
	Port       uint16    `json:"port"`
	Expiration time.Time `json:"expires_at"`
}

type piaPortForwardData struct {
	Port       uint16    `json:"port"`
	Token      string    `json:"token"`
	Signature  string    `json:"signature"`
	Expiration time.Time `json:"expires_at"`
}

func readPIAPortForwardData(openFile os.OpenFileFunc) (data piaPortForwardData, err error) {
	const filepath = string(constants.PIAPortForward)
	file, err := openFile(filepath, os.O_RDONLY, 0)
	if os.IsNotExist(err) {
		return data, nil
	} else if err != nil {
		return data, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		_ = file.Close()
		return data, err
	}
	return data, file.Close()
}

func writePIAPortForwardData(openFile os.OpenFileFunc, data piaPortForwardData) (err error) {
	const filepath = string(constants.PIAPortForward)
	file, err := openFile(filepath,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0644)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}

func unpackPIAPayload(payload string) (port uint16, token string, expiration time.Time, err error) {
	b, err := base64.RawStdEncoding.DecodeString(payload)
	if err != nil {
		return 0, "", expiration, fmt.Errorf("cannot decode payload: %w", err)
	}
	var payloadData piaPayload
	if err := json.Unmarshal(b, &payloadData); err != nil {
		return 0, "", expiration, fmt.Errorf("cannot parse payload data: %w", err)
	}
	return payloadData.Port, payloadData.Token, payloadData.Expiration, nil
}

func packPIAPayload(port uint16, token string, expiration time.Time) (payload string, err error) {
	payloadData := piaPayload{
		Token:      token,
		Port:       port,
		Expiration: expiration,
	}
	b, err := json.Marshal(&payloadData)
	if err != nil {
		return "", fmt.Errorf("cannot serialize payload data: %w", err)
	}
	payload = base64.RawStdEncoding.EncodeToString(b)
	return payload, nil
}

func fetchPIAToken(ctx context.Context, openFile os.OpenFileFunc,
	client *http.Client) (token string, err error) {
	username, password, err := getOpenvpnCredentials(openFile)
	if err != nil {
		return "", fmt.Errorf("cannot get Openvpn credentials: %w", err)
	}
	url := url.URL{
		Scheme: "https",
		User:   url.UserPassword(username, password),
		Host:   "10.0.0.1",
		Path:   "/authv3/generateToken",
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return "", err
	}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(response.Body)
		shortenMessage := string(b)
		shortenMessage = strings.ReplaceAll(shortenMessage, "\n", "")
		shortenMessage = strings.ReplaceAll(shortenMessage, "  ", " ")
		return "", fmt.Errorf("%s: response received: %q", response.Status, shortenMessage)
	}
	decoder := json.NewDecoder(response.Body)
	var result struct {
		Token string `json:"token"`
	}
	if err := decoder.Decode(&result); err != nil {
		return "", err
	} else if len(result.Token) == 0 {
		return "", fmt.Errorf("token is empty")
	}
	return result.Token, nil
}

func getOpenvpnCredentials(openFile os.OpenFileFunc) (username, password string, err error) {
	const filepath = string(constants.OpenVPNAuthConf)
	file, err := openFile(filepath, os.O_RDONLY, 0)
	if err != nil {
		return "", "", fmt.Errorf("cannot read openvpn auth file: %s", err)
	}
	authData, err := ioutil.ReadAll(file)
	if err != nil {
		_ = file.Close()
		return "", "", fmt.Errorf("cannot read openvpn auth file: %s", err)
	}
	if err := file.Close(); err != nil {
		return "", "", err
	}
	lines := strings.Split(string(authData), "\n")
	const minLines = 2
	if len(lines) < minLines {
		return "", "", fmt.Errorf("not enough lines (%d) in openvpn auth file", len(lines))
	}
	username, password = lines[0], lines[1]
	return username, password, nil
}

func fetchPIAPortForwardData(ctx context.Context, client *http.Client, gateway net.IP, token string) (
	port uint16, signature string, expiration time.Time, err error) {
	queryParams := url.Values{}
	queryParams.Add("token", token)
	url := url.URL{
		Scheme:   "https",
		Host:     net.JoinHostPort(gateway.String(), "19999"),
		Path:     "/getSignature",
		RawQuery: queryParams.Encode(),
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return 0, "", expiration, fmt.Errorf("cannot obtain signature: %w", err)
	}
	response, err := client.Do(request)
	if err != nil {
		return 0, "", expiration, fmt.Errorf("cannot obtain signature: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return 0, "", expiration, fmt.Errorf("cannot obtain signature: %s", response.Status)
	}
	decoder := json.NewDecoder(response.Body)
	var data struct {
		Status    string `json:"status"`
		Payload   string `json:"payload"`
		Signature string `json:"signature"`
	}
	if err := decoder.Decode(&data); err != nil {
		return 0, "", expiration, fmt.Errorf("cannot decode received data: %w", err)
	} else if data.Status != "OK" {
		return 0, "", expiration, fmt.Errorf("response received from PIA has status %s", data.Status)
	}

	port, _, expiration, err = unpackPIAPayload(data.Payload)
	return port, data.Signature, expiration, err
}

func bindPIAPort(ctx context.Context, client *http.Client, gateway net.IP, data piaPortForwardData) (err error) {
	payload, err := packPIAPayload(data.Port, data.Token, data.Expiration)
	if err != nil {
		return err
	}
	queryParams := url.Values{}
	queryParams.Add("payload", payload)
	queryParams.Add("signature", data.Signature)
	url := url.URL{
		Scheme:   "https",
		Host:     net.JoinHostPort(gateway.String(), "19999"),
		Path:     "/bindPort",
		RawQuery: queryParams.Encode(),
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return fmt.Errorf("cannot bind port: %w", err)
	}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("cannot bind port: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot bind port: %s", response.Status)
	}

	decoder := json.NewDecoder(response.Body)
	var responseData struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := decoder.Decode(&responseData); err != nil {
		return fmt.Errorf("cannot bind port: %w", err)
	} else if responseData.Status != "OK" {
		return fmt.Errorf("response received from PIA: %s (%s)", responseData.Status, responseData.Message)
	}
	return nil
}

func writePortForwardedToFile(openFile os.OpenFileFunc,
	filepath string, port uint16) (err error) {
	file, err := openFile(filepath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(fmt.Sprintf("%d", port)))
	if err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}
