package privateinternetaccess

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/provider/utils"
	"github.com/qdm12/golibs/format"
)

var (
	ErrGatewayIPIsNil           = errors.New("gateway IP address is nil")
	ErrServerNameEmpty          = errors.New("server name is empty")
	ErrCreateHTTPClient         = errors.New("cannot create custom HTTP client")
	ErrReadSavedPortForwardData = errors.New("cannot read saved port forwarded data")
	ErrRefreshPortForwardData   = errors.New("cannot refresh port forward data")
	ErrBindPort                 = errors.New("cannot bind port")
)

// PortForward obtains a VPN server side port forwarded from PIA.
func (p *PIA) PortForward(ctx context.Context, client *http.Client,
	logger utils.Logger, gateway net.IP, serverName string) (
	port uint16, err error) {
	server := constants.PIAServerWhereName(p.servers, serverName)
	if !server.PortForward {
		logger.Error("The server " + serverName +
			" (region " + server.Region + ") does not support port forwarding")
		return 0, nil
	}
	if gateway == nil {
		return 0, ErrGatewayIPIsNil
	} else if serverName == "" {
		return 0, ErrServerNameEmpty
	}

	privateIPClient, err := newHTTPClient(serverName)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrCreateHTTPClient, err)
	}

	data, err := readPIAPortForwardData(p.portForwardPath)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrReadSavedPortForwardData, err)
	}

	dataFound := data.Port > 0
	durationToExpiration := data.Expiration.Sub(p.timeNow())
	expired := durationToExpiration <= 0

	if dataFound {
		logger.Info("Found saved forwarded port data for port " + strconv.Itoa(int(data.Port)))
		if expired {
			logger.Warn("Forwarded port data expired on " +
				data.Expiration.Format(time.RFC1123) + ", getting another one")
		} else {
			logger.Info("Forwarded port data expires in " + format.FriendlyDuration(durationToExpiration))
		}
	}

	if !dataFound || expired {
		data, err = refreshPIAPortForwardData(ctx, client, privateIPClient, gateway,
			p.portForwardPath, p.authFilePath)
		if err != nil {
			return 0, fmt.Errorf("%w: %s", ErrRefreshPortForwardData, err)
		}
		durationToExpiration = data.Expiration.Sub(p.timeNow())
	}
	logger.Info("Port forwarded data expires in " + format.FriendlyDuration(durationToExpiration))

	// First time binding
	if err := bindPort(ctx, privateIPClient, gateway, data); err != nil {
		return 0, fmt.Errorf("%w: %s", ErrBindPort, err)
	}

	return data.Port, nil
}

var (
	ErrPortForwardedExpired = errors.New("port forwarded data expired")
)

func (p *PIA) KeepPortForward(ctx context.Context, client *http.Client,
	port uint16, gateway net.IP, serverName string) (err error) {
	privateIPClient, err := newHTTPClient(serverName)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCreateHTTPClient, err)
	}

	data, err := readPIAPortForwardData(p.portForwardPath)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrReadSavedPortForwardData, err)
	}

	durationToExpiration := data.Expiration.Sub(p.timeNow())
	expiryTimer := time.NewTimer(durationToExpiration)
	const keepAlivePeriod = 15 * time.Minute
	// Timer behaving as a ticker
	keepAliveTimer := time.NewTimer(keepAlivePeriod)

	for {
		select {
		case <-ctx.Done():
			if !keepAliveTimer.Stop() {
				<-keepAliveTimer.C
			}
			if !expiryTimer.Stop() {
				<-expiryTimer.C
			}
			return ctx.Err()
		case <-keepAliveTimer.C:
			err := bindPort(ctx, privateIPClient, gateway, data)
			if err != nil {
				return fmt.Errorf("%w: %s", ErrBindPort, err)
			}
			keepAliveTimer.Reset(keepAlivePeriod)
		case <-expiryTimer.C:
			return fmt.Errorf("%w: on %s", ErrPortForwardedExpired,
				data.Expiration.Format(time.RFC1123))
		}
	}
}

var (
	ErrFetchToken            = errors.New("cannot fetch token")
	ErrFetchPortForwarding   = errors.New("cannot fetch port forwarding data")
	ErrPersistPortForwarding = errors.New("cannot persist port forwarding data")
)

func refreshPIAPortForwardData(ctx context.Context, client, privateIPClient *http.Client,
	gateway net.IP, portForwardPath, authFilePath string) (data piaPortForwardData, err error) {
	data.Token, err = fetchToken(ctx, client, authFilePath)
	if err != nil {
		return data, fmt.Errorf("%w: %s", ErrFetchToken, err)
	}

	data.Port, data.Signature, data.Expiration, err = fetchPortForwardData(ctx, privateIPClient, gateway, data.Token)
	if err != nil {
		return data, fmt.Errorf("%w: %s", ErrFetchPortForwarding, err)
	}

	if err := writePIAPortForwardData(portForwardPath, data); err != nil {
		return data, fmt.Errorf("%w: %s", ErrPersistPortForwarding, err)
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

func readPIAPortForwardData(portForwardPath string) (data piaPortForwardData, err error) {
	file, err := os.Open(portForwardPath)
	if os.IsNotExist(err) {
		return data, nil
	} else if err != nil {
		return data, err
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		_ = file.Close()
		return data, err
	}

	return data, file.Close()
}

func writePIAPortForwardData(portForwardPath string, data piaPortForwardData) (err error) {
	file, err := os.OpenFile(portForwardPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(data); err != nil {
		_ = file.Close()
		return err
	}

	return file.Close()
}

func unpackPayload(payload string) (port uint16, token string, expiration time.Time, err error) {
	b, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return 0, "", expiration,
			fmt.Errorf("%w: for payload: %s", err, payload)
	}

	var payloadData piaPayload
	if err := json.Unmarshal(b, &payloadData); err != nil {
		return 0, "", expiration,
			fmt.Errorf("%w: for data: %s", err, string(b))
	}

	return payloadData.Port, payloadData.Token, payloadData.Expiration, nil
}

func packPayload(port uint16, token string, expiration time.Time) (payload string, err error) {
	payloadData := piaPayload{
		Token:      token,
		Port:       port,
		Expiration: expiration,
	}

	b, err := json.Marshal(&payloadData)
	if err != nil {
		return "", err
	}

	payload = base64.StdEncoding.EncodeToString(b)
	return payload, nil
}

var (
	errGetCredentials = errors.New("cannot get username and password")
	errEmptyToken     = errors.New("token received is empty")
)

func fetchToken(ctx context.Context, client *http.Client,
	authFilePath string) (token string, err error) {
	username, password, err := getOpenvpnCredentials(authFilePath)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errGetCredentials, err)
	}

	errSubstitutions := map[string]string{
		username: "<username>",
		password: "<password>",
	}

	url := url.URL{
		Scheme: "https",
		User:   url.UserPassword(username, password),
		Host:   "privateinternetaccess.com",
		Path:   "/gtoken/generateToken",
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return "", replaceInErr(err, errSubstitutions)
	}

	response, err := client.Do(request)
	if err != nil {
		return "", replaceInErr(err, errSubstitutions)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", makeNOKStatusError(response, nil)
	}

	decoder := json.NewDecoder(response.Body)
	var result struct {
		Token string `json:"token"`
	}
	if err := decoder.Decode(&result); err != nil {
		return "", fmt.Errorf("%w: %s", ErrUnmarshalResponse, err)
	}

	if result.Token == "" {
		return "", errEmptyToken
	}
	return result.Token, nil
}

var (
	errAuthFileRead      = errors.New("cannot read OpenVPN authentication file")
	errAuthFileMalformed = errors.New("authentication file is malformed")
)

func getOpenvpnCredentials(authFilePath string) (
	username, password string, err error) {
	file, err := os.Open(authFilePath)
	if err != nil {
		return "", "", fmt.Errorf("%w: %s", errAuthFileRead, err)
	}

	authData, err := io.ReadAll(file)
	if err != nil {
		_ = file.Close()
		return "", "", fmt.Errorf("%w: %s", errAuthFileRead, err)
	}

	if err := file.Close(); err != nil {
		return "", "", err
	}

	lines := strings.Split(string(authData), "\n")
	const minLines = 2
	if len(lines) < minLines {
		return "", "", fmt.Errorf("%w: only %d lines exist", errAuthFileMalformed, len(lines))
	}

	username, password = lines[0], lines[1]
	return username, password, nil
}

var (
	errGetSignaturePayload = errors.New("cannot obtain signature payload")
	errUnpackPayload       = errors.New("cannot unpack payload data")
)

func fetchPortForwardData(ctx context.Context, client *http.Client, gateway net.IP, token string) (
	port uint16, signature string, expiration time.Time, err error) {
	errSubstitutions := map[string]string{token: "<token>"}

	queryParams := make(url.Values)
	queryParams.Add("token", token)
	url := url.URL{
		Scheme:   "https",
		Host:     net.JoinHostPort(gateway.String(), "19999"),
		Path:     "/getSignature",
		RawQuery: queryParams.Encode(),
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		err = replaceInErr(err, errSubstitutions)
		return 0, "", expiration, fmt.Errorf("%w: %s", errGetSignaturePayload, err)
	}

	response, err := client.Do(request)
	if err != nil {
		err = replaceInErr(err, errSubstitutions)
		return 0, "", expiration, fmt.Errorf("%w: %s", errGetSignaturePayload, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, "", expiration, makeNOKStatusError(response, errSubstitutions)
	}

	decoder := json.NewDecoder(response.Body)
	var data struct {
		Status    string `json:"status"`
		Payload   string `json:"payload"`
		Signature string `json:"signature"`
	}
	if err := decoder.Decode(&data); err != nil {
		return 0, "", expiration, fmt.Errorf("%w: %s", ErrUnmarshalResponse, err)
	}

	if data.Status != "OK" {
		return 0, "", expiration, fmt.Errorf("%w: status is: %s", ErrBadResponse, data.Status)
	}

	port, _, expiration, err = unpackPayload(data.Payload)
	if err != nil {
		return 0, "", expiration, fmt.Errorf("%w: %s", errUnpackPayload, err)
	}
	return port, data.Signature, expiration, err
}

var (
	ErrSerializePayload  = errors.New("cannot serialize payload")
	ErrUnmarshalResponse = errors.New("cannot unmarshal response")
	ErrBadResponse       = errors.New("bad response received")
)

func bindPort(ctx context.Context, client *http.Client, gateway net.IP, data piaPortForwardData) (err error) {
	payload, err := packPayload(data.Port, data.Token, data.Expiration)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrSerializePayload, err)
	}

	queryParams := make(url.Values)
	queryParams.Add("payload", payload)
	queryParams.Add("signature", data.Signature)
	url := url.URL{
		Scheme:   "https",
		Host:     net.JoinHostPort(gateway.String(), "19999"),
		Path:     "/bindPort",
		RawQuery: queryParams.Encode(),
	}

	errSubstitutions := map[string]string{
		payload:        "<payload>",
		data.Signature: "<signature>",
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return replaceInErr(err, errSubstitutions)
	}

	response, err := client.Do(request)
	if err != nil {
		return replaceInErr(err, errSubstitutions)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return makeNOKStatusError(response, errSubstitutions)
	}

	decoder := json.NewDecoder(response.Body)
	var responseData struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := decoder.Decode(&responseData); err != nil {
		return fmt.Errorf("%w: from %s: %s", ErrUnmarshalResponse, url.String(), err)
	}

	if responseData.Status != "OK" {
		return fmt.Errorf("%w: %s: %s", ErrBadResponse, responseData.Status, responseData.Message)
	}

	return nil
}

// replaceInErr is used to remove sensitive information from errors.
func replaceInErr(err error, substitutions map[string]string) error {
	s := replaceInString(err.Error(), substitutions)
	return errors.New(s) //nolint:goerr113
}

// replaceInString is used to remove sensitive information.
func replaceInString(s string, substitutions map[string]string) string {
	for old, new := range substitutions {
		s = strings.ReplaceAll(s, old, new)
	}
	return s
}

var ErrHTTPStatusCodeNotOK = errors.New("HTTP status code is not OK")

func makeNOKStatusError(response *http.Response, substitutions map[string]string) (err error) {
	url := response.Request.URL.String()
	url = replaceInString(url, substitutions)

	b, _ := io.ReadAll(response.Body)
	shortenMessage := string(b)
	shortenMessage = strings.ReplaceAll(shortenMessage, "\n", "")
	shortenMessage = strings.ReplaceAll(shortenMessage, "  ", " ")
	shortenMessage = replaceInString(shortenMessage, substitutions)

	return fmt.Errorf("%w: %s: %s: response received: %s",
		ErrHTTPStatusCodeNotOK, url, response.Status, shortenMessage)
}
