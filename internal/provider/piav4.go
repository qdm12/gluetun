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
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	gluetunLog "github.com/qdm12/gluetun/internal/logging"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
)

type piaV4 struct {
	servers    []models.PIAServer
	timeNow    timeNowFunc
	randSource rand.Source
}

func newPrivateInternetAccessV4(servers []models.PIAServer, timeNow timeNowFunc) *piaV4 {
	return &piaV4{
		servers:    servers,
		timeNow:    timeNow,
		randSource: rand.NewSource(timeNow().UnixNano()),
	}
}

func (p *piaV4) GetOpenVPNConnection(selection models.ServerSelection) (connection models.OpenVPNConnection, err error) {
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
		return connection, fmt.Errorf("combination of protocol %q and encryption %q does not yield any port number", selection.Protocol, selection.EncryptionPreset)
	}

	if selection.TargetIP != nil {
		return models.OpenVPNConnection{IP: selection.TargetIP, Port: port, Protocol: selection.Protocol}, nil
	}

	servers := filterPIAServers(p.servers, selection.Region)
	if len(servers) == 0 {
		return connection, fmt.Errorf("no server found for region %q", selection.Region)
	}

	var connections []models.OpenVPNConnection
	for _, server := range servers {
		IPs := server.OpenvpnUDP.IPs
		if selection.Protocol == constants.TCP {
			IPs = server.OpenvpnTCP.IPs
		}
		for _, IP := range IPs {
			connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: selection.Protocol})
		}
	}

	return pickRandomConnection(connections, p.randSource), nil
}

func (p *piaV4) BuildConf(connection models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string) {
	return buildPIAConf(connection, verbosity, root, cipher, auth, extras)
}

//nolint:gocognit
func (p *piaV4) PortForward(ctx context.Context, client *http.Client,
	fileManager files.FileManager, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath models.Filepath)) {
	if gateway == nil {
		pfLogger.Error("aborting because: VPN gateway IP address was not found")
		return
	}
	client, err := newPIAv4HTTPClient()
	if err != nil {
		pfLogger.Error("aborting because: %s", err)
		return
	}
	defer pfLogger.Warn("loop exited")
	data, err := readPIAPortForwardData(fileManager)
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
			data, err = refreshPIAPortForwardData(client, gateway, fileManager)
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
		return bindPIAPort(client, gateway, data)
	})
	if ctx.Err() != nil {
		return
	}

	filepath := syncState(data.Port)
	pfLogger.Info("Writing port to %s", filepath)
	if err := fileManager.WriteToFile(
		string(filepath), []byte(fmt.Sprintf("%d", data.Port)),
		files.Permissions(0666),
	); err != nil {
		pfLogger.Error(err)
	}

	if err := fw.SetAllowedPort(ctx, data.Port, string(constants.TUN)); err != nil {
		pfLogger.Error(err)
	}

	expiryTimer := time.NewTimer(durationToExpiration)
	defer expiryTimer.Stop()
	const keepAlivePeriod = 15 * time.Minute
	keepAliveTicker := time.NewTicker(keepAlivePeriod)
	defer keepAliveTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			removeCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			if err := fw.RemoveAllowedPort(removeCtx, data.Port); err != nil {
				pfLogger.Error(err)
			}
			return
		case <-keepAliveTicker.C:
			if err := bindPIAPort(client, gateway, data); err != nil {
				pfLogger.Error(err)
			}
		case <-expiryTimer.C:
			pfLogger.Warn("Forward port has expired on %s, getting another one", data.Expiration.Format(time.RFC1123))
			oldPort := data.Port
			for {
				data, err = refreshPIAPortForwardData(client, gateway, fileManager)
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
			if err := fileManager.WriteToFile(
				string(filepath), []byte(fmt.Sprintf("%d", data.Port)),
				files.Permissions(0666),
			); err != nil {
				pfLogger.Error(err)
			}
			if err := bindPIAPort(client, gateway, data); err != nil {
				pfLogger.Error(err)
			}
			keepAliveTicker.Reset(keepAlivePeriod)
			expiryTimer.Reset(durationToExpiration)
		}
	}
}

func filterPIAServers(servers []models.PIAServer, region string) (filtered []models.PIAServer) {
	if len(region) == 0 {
		return servers
	}
	for _, server := range servers {
		if strings.EqualFold(server.Region, region) {
			return []models.PIAServer{server}
		}
	}
	return nil
}

func newPIAv4HTTPClient() (client *http.Client, err error) {
	certificateBytes, err := base64.StdEncoding.DecodeString(constants.PIACertificateStrong)
	if err != nil {
		return nil, fmt.Errorf("cannot decode PIA root certificate: %w", err)
	}
	certificate, err := x509.ParseCertificate(certificateBytes)
	if err != nil {
		return nil, fmt.Errorf("cannot parse PIA root certificate: %w", err)
	}
	rootCAs := x509.NewCertPool()
	rootCAs.AddCert(certificate)
	TLSClientConfig := &tls.Config{
		RootCAs:            rootCAs,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true, //nolint:gosec
	} // TODO fix and remove InsecureSkipVerify
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
	const httpTimeout = 5 * time.Second
	client = &http.Client{Transport: &transport, Timeout: httpTimeout}
	return client, nil
}

func refreshPIAPortForwardData(client *http.Client, gateway net.IP, fileManager files.FileManager) (data piaPortForwardData, err error) {
	data.Token, err = fetchPIAToken(fileManager, client)
	if err != nil {
		return data, fmt.Errorf("cannot obtain token: %w", err)
	}
	data.Port, data.Signature, data.Expiration, err = fetchPIAPortForwardData(client, gateway, data.Token)
	if err != nil {
		if strings.HasSuffix(err.Error(), "connection refused") {
			return data, fmt.Errorf("cannot obtain port forwarding data: connection was refused, are you sure the region you are using supports port forwarding ;)")
		}
		return data, fmt.Errorf("cannot obtain port forwarding data: %w", err)
	}
	if err := writePIAPortForwardData(fileManager, data); err != nil {
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

func readPIAPortForwardData(fileManager files.FileManager) (data piaPortForwardData, err error) {
	const filepath = string(constants.PIAPortForward)
	exists, err := fileManager.FileExists(filepath)
	if err != nil {
		return data, err
	} else if !exists {
		return data, nil
	}
	b, err := fileManager.ReadFile(filepath)
	if err != nil {
		return data, err
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return data, err
	}
	return data, nil
}

func writePIAPortForwardData(fileManager files.FileManager, data piaPortForwardData) (err error) {
	b, err := json.Marshal(&data)
	if err != nil {
		return fmt.Errorf("cannot encode data: %w", err)
	}
	err = fileManager.WriteToFile(string(constants.PIAPortForward), b)
	if err != nil {
		return err
	}
	return nil
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

func fetchPIAToken(fileManager files.FileManager, client *http.Client) (token string, err error) {
	username, password, err := getOpenvpnCredentials(fileManager)
	if err != nil {
		return "", fmt.Errorf("cannot get Openvpn credentials: %w", err)
	}
	url := url.URL{
		Scheme: "https",
		User:   url.UserPassword(username, password),
		Host:   "10.0.0.1",
		Path:   "/authv3/generateToken",
	}
	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return "", err
	}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {
		shortenMessage := string(b)
		shortenMessage = strings.ReplaceAll(shortenMessage, "\n", "")
		shortenMessage = strings.ReplaceAll(shortenMessage, "  ", " ")
		return "", fmt.Errorf("%s: response received: %q", response.Status, shortenMessage)
	} else if err != nil {
		return "", err
	}
	var result struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(b, &result); err != nil {
		return "", err
	} else if len(result.Token) == 0 {
		return "", fmt.Errorf("token is empty")
	}
	return result.Token, nil
}

func getOpenvpnCredentials(fileManager files.FileManager) (username, password string, err error) {
	authData, err := fileManager.ReadFile(string(constants.OpenVPNAuthConf))
	if err != nil {
		return "", "", fmt.Errorf("cannot read openvpn auth file: %w", err)
	}
	lines := strings.Split(string(authData), "\n")
	if len(lines) < 2 {
		return "", "", fmt.Errorf("not enough lines (%d) in openvpn auth file", len(lines))
	}
	username, password = lines[0], lines[1]
	return username, password, nil
}

func fetchPIAPortForwardData(client *http.Client, gateway net.IP, token string) (port uint16, signature string, expiration time.Time, err error) {
	queryParams := url.Values{}
	queryParams.Add("token", token)
	url := url.URL{
		Scheme:   "https",
		Host:     net.JoinHostPort(gateway.String(), "19999"),
		Path:     "/getSignature",
		RawQuery: queryParams.Encode(),
	}
	response, err := client.Get(url.String())
	if err != nil {
		return 0, "", expiration, fmt.Errorf("cannot obtain signature: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return 0, "", expiration, fmt.Errorf("cannot obtain signature: %s", response.Status)
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, "", expiration, fmt.Errorf("cannot obtain signature: %w", err)
	}
	var data struct {
		Status    string `json:"status"`
		Payload   string `json:"payload"`
		Signature string `json:"signature"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return 0, "", expiration, fmt.Errorf("cannot decode received data: %w", err)
	} else if data.Status != "OK" {
		return 0, "", expiration, fmt.Errorf("response received from PIA has status %s", data.Status)
	}

	port, _, expiration, err = unpackPIAPayload(data.Payload)
	return port, data.Signature, expiration, err
}

func bindPIAPort(client *http.Client, gateway net.IP, data piaPortForwardData) (err error) {
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

	response, err := client.Get(url.String())
	if err != nil {
		return fmt.Errorf("cannot bind port: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot bind port: %s", response.Status)
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("cannot bind port: %w", err)
	}
	var responseData struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(b, &responseData); err != nil {
		return fmt.Errorf("cannot bind port: %w", err)
	} else if responseData.Status != "OK" {
		return fmt.Errorf("response received from PIA: %s (%s)", responseData.Status, responseData.Message)
	}
	return nil
}
