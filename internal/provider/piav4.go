package provider

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
)

type getVPNGatewayFunc func() (gateway net.IP, err error)

type piaV4 struct {
	servers       []models.PIAServer
	getVPNGateway getVPNGatewayFunc
	fileManager   files.FileManager
	logger        logging.Logger
	timeNow       func() time.Time
}

func newPrivateInternetAccessV4(servers []models.PIAServer, getVPNGateway getVPNGatewayFunc, fileManager files.FileManager, logger logging.Logger) *piaV4 {
	return &piaV4{
		servers:       servers,
		getVPNGateway: getVPNGateway,
		fileManager:   fileManager,
		logger:        logger,
		timeNow:       time.Now,
	}
}

func (p *piaV4) GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error) {
	return getPIAOpenVPNConnections(p.servers, selection)
}

func (p *piaV4) BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string) {
	return buildPIAConf(connections, verbosity, root, cipher, auth, extras)
}

func (p *piaV4) GetPortForward(ctx context.Context, wg *sync.WaitGroup, client *http.Client) (port uint16, err error) {
	defer wg.Done()
	gateway, err := p.getVPNGateway()
	if err != nil {
		return 0, fmt.Errorf("cannot obtain VPN gateway: %w", err)
	}

	data, readErr := readPIAPortForwardData(p.fileManager)
	if errors.Is(readErr, errPIAPortForwardFileNotExists) {
		data, err = p.refreshPortForwardData(client, gateway)
		if err != nil {
			return 0, err
		}
	} else if readErr != nil {
		return 0, readErr
	}
	now := p.timeNow()
	durationToExpiration := data.Expiration.Sub(now)
	if durationToExpiration <= 0 {
		p.logger.Warn("Forward port has expired on %s, getting another one", data.Expiration)
		data, err = p.refreshPortForwardData(client, gateway)
		if err != nil {
			return 0, err
		}
	}
	if err := bindPIAPort(client, gateway, data); err != nil {
		p.logger.Error(err)
	}

	wg.Add(1)
	go p.maintainPortForwarding(ctx, wg, durationToExpiration, client, gateway, data)

	return data.Port, nil
}

func (p *piaV4) maintainPortForwarding(ctx context.Context, wg *sync.WaitGroup, durationToExpiration time.Duration, client *http.Client, gateway net.IP, data piaPortForwardData) {
	defer wg.Done()

	expiryTimer := time.NewTimer(durationToExpiration)
	defer expiryTimer.Stop()
	const keepAlivePeriod = 15 * time.Minute
	keepAliveTicker := time.NewTicker(keepAlivePeriod)
	defer keepAliveTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-keepAliveTicker.C:
			if err := bindPIAPort(client, gateway, data); err != nil {
				p.logger.Error(err)
			}
		case <-expiryTimer.C:
			p.logger.Warn("Forward port has expired on %s, getting another one", data.Expiration)
			data, err := p.refreshPortForwardData(client, gateway)
			if err != nil {
				p.logger.Error(err)
			}
			now := p.timeNow()
			durationToExpiration := data.Expiration.Sub(now)
			expiryTimer.Reset(durationToExpiration)
			// TODO send port in channel for firewall or move firewall changes here (better?)
		}
	}
}

func (p *piaV4) refreshPortForwardData(client *http.Client, gateway net.IP) (data piaPortForwardData, err error) {
	data.Token, err = fetchPIAToken(p.fileManager, client)
	if err != nil {
		return data, err
	}
	data.Port, data.Signature, data.Expiration, err = fetchPIAPortForwardData(client, gateway, data.Token)
	if err != nil {
		return data, err
	}
	if err := p.writePortForwardData(data); err != nil {
		return data, err
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

var errPIAPortForwardFileNotExists = errors.New("PIA port forward data file does not exist")

func readPIAPortForwardData(fileManager files.FileManager) (data piaPortForwardData, err error) {
	const filepath = string(constants.PIAPortForward)
	exists, err := fileManager.FileExists(filepath)
	if err != nil {
		return data, err
	} else if !exists {
		return data, errPIAPortForwardFileNotExists
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

func (p *piaV4) writePortForwardData(data piaPortForwardData) (err error) {
	b, err := json.Marshal(&data)
	if err != nil {
		return fmt.Errorf("cannot encode output data: %w", err)
	}
	err = p.fileManager.WriteToFile(string(constants.PIAPortForward), b)
	if err != nil {
		return fmt.Errorf("cannot persist port forwarding information to file: %w", err)
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
	var payloadData piaPayload
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
	const url = "https://www.privateinternetaccess.com/api/client/v2/token"
	data := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{username, password}
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf(response.Status)
	}
	b, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var result struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(b, &result); err != nil {
		return "", err
	}
	if len(result.Token) == 0 {
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
		return 0, "", expiration, fmt.Errorf("cannot obtain signature: %w", err)
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
		return fmt.Errorf("response received from PIA has status %s", responseData.Status)
	}
	return nil
}
