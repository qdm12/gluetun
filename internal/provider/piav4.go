package provider

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
)

type piaV4 struct {
	servers []models.PIAServer
	timeNow func() time.Time
}

func newPrivateInternetAccessV4(servers []models.PIAServer) *piaV4 {
	return &piaV4{
		servers: servers,
		timeNow: time.Now,
	}
}

func (p *piaV4) GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error) {
	return getPIAOpenVPNConnections(p.servers, selection)
}

func (p *piaV4) BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string) {
	return buildPIAConf(connections, verbosity, root, cipher, auth, extras)
}

func (p *piaV4) PortForward(ctx context.Context, client *http.Client,
	fileManager files.FileManager, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath models.Filepath)) {
	if gateway == nil {
		pfLogger.Error("VPN gateway IP address was not found, cannot do anything")
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
			pfLogger.Warn("Forwarded port data expired on %s, getting another one", data.Expiration)
		} else {
			pfLogger.Info("Forwarded port data expires in %s", durationToExpiration)
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
	pfLogger.Info("Port forwarded is %d expiring in %s", data.Port, durationToExpiration)

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
			pfLogger.Warn("Forward port has expired on %s, getting another one", data.Expiration)
			oldPort := data.Port
			for {
				data, err = refreshPIAPortForwardData(client, gateway, fileManager)
				if err != nil {
					pfLogger.Error(err)
					continue
				}
				break
			}
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
			durationToExpiration := data.Expiration.Sub(p.timeNow())
			expiryTimer.Reset(durationToExpiration)
			if err := bindPIAPort(client, gateway, data); err != nil {
				pfLogger.Error(err)
			}
			keepAliveTicker.Reset(keepAlivePeriod)
		}
	}
}

func refreshPIAPortForwardData(client *http.Client, gateway net.IP, fileManager files.FileManager) (data piaPortForwardData, err error) {
	data.Token, err = fetchPIAToken(fileManager, client)
	if err != nil {
		return data, err
	}
	data.Port, data.Signature, data.Expiration, err = fetchPIAPortForwardData(client, gateway, data.Token)
	if err != nil {
		return data, err
	}
	if err := writePIAPortForwardData(fileManager, data); err != nil {
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
		return fmt.Errorf("cannot encode output data: %w", err)
	}
	err = fileManager.WriteToFile(string(constants.PIAPortForward), b)
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
