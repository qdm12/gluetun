package privateinternetaccess

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/format"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

var ErrServerNameNotFound = errors.New("server name not found in servers")

// PortForward obtains a VPN server side port forwarded from PIA.
func (p *Provider) PortForward(ctx context.Context,
	objects utils.PortForwardObjects,
) (ports []uint16, err error) {
	switch {
	case objects.ServerName == "":
		panic("server name cannot be empty")
	case !objects.Gateway.IsValid():
		panic("gateway is not set")
	case objects.Username == "":
		panic("username is not set")
	case objects.Password == "":
		panic("password is not set")
	}

	serverName := objects.ServerName
	apiIP := buildAPIIPAddress(objects.Gateway)
	logger := objects.Logger

	if !objects.CanPortForward {
		return nil, fmt.Errorf("%w: for server %s", ErrServerNameNotFound, serverName)
	}

	privateIPClient, err := newHTTPClient(serverName)
	if err != nil {
		return nil, fmt.Errorf("creating custom HTTP client: %w", err)
	}

	data, err := readPIAPortForwardData(p.portForwardPath)
	if err != nil {
		return nil, fmt.Errorf("reading saved port forwarded data: %w", err)
	}

	dataFound := data.Port > 0
	durationToExpiration := data.Expiration.Sub(p.timeNow())
	expired := durationToExpiration <= 0

	if dataFound {
		logger.Info("Found saved forwarded port data for port " + strconv.Itoa(int(data.Port)))
		if expired {
			logger.Warn("Forwarded port data expired on " +
				data.Expiration.Format(time.RFC1123) + ", getting another one")
		}
	}

	if !dataFound || expired {
		client := objects.Client
		data, err = refreshPIAPortForwardData(ctx, client, privateIPClient, apiIP,
			p.portForwardPath, objects.Username, objects.Password)
		if err != nil {
			return nil, fmt.Errorf("refreshing port forward data: %w", err)
		}
		durationToExpiration = data.Expiration.Sub(p.timeNow())
	}
	logger.Info("Port forwarded data expires in " + format.FriendlyDuration(durationToExpiration))

	// First time binding
	if err := bindPort(ctx, privateIPClient, apiIP, data); err != nil {
		return nil, fmt.Errorf("binding port: %w", err)
	}

	return []uint16{data.Port}, nil
}

var ErrPortForwardedExpired = errors.New("port forwarded data expired")

func (p *Provider) KeepPortForward(ctx context.Context,
	objects utils.PortForwardObjects,
) (err error) {
	switch {
	case objects.ServerName == "":
		panic("server name cannot be empty")
	case !objects.Gateway.IsValid():
		panic("gateway is not set")
	}

	apiIP := buildAPIIPAddress(objects.Gateway)

	privateIPClient, err := newHTTPClient(objects.ServerName)
	if err != nil {
		return fmt.Errorf("creating custom HTTP client: %w", err)
	}

	data, err := readPIAPortForwardData(p.portForwardPath)
	if err != nil {
		return fmt.Errorf("reading saved port forwarded data: %w", err)
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
			err = bindPort(ctx, privateIPClient, apiIP, data)
			if err != nil {
				return fmt.Errorf("binding port: %w", err)
			}
			keepAliveTimer.Reset(keepAlivePeriod)
		case <-expiryTimer.C:
			return fmt.Errorf("%w: on %s", ErrPortForwardedExpired,
				data.Expiration.Format(time.RFC1123))
		}
	}
}

func buildAPIIPAddress(gateway netip.Addr) (api netip.Addr) {
	if gateway.Is6() {
		panic("IPv6 gateway not supported")
	}

	gatewayBytes := gateway.As4()
	gatewayBytes[2] = 128
	gatewayBytes[3] = 1
	return netip.AddrFrom4(gatewayBytes)
}

func refreshPIAPortForwardData(ctx context.Context, client, privateIPClient *http.Client,
	apiIP netip.Addr, portForwardPath, username, password string,
) (data piaPortForwardData, err error) {
	data.Token, err = fetchToken(ctx, client, username, password)
	if err != nil {
		return data, fmt.Errorf("fetching token: %w", err)
	}

	data.Port, data.Signature, data.Expiration, err = fetchPortForwardData(ctx, privateIPClient, apiIP, data.Token)
	if err != nil {
		return data, fmt.Errorf("fetching port forwarding data: %w", err)
	}

	if err := writePIAPortForwardData(portForwardPath, data); err != nil {
		return data, fmt.Errorf("persisting port forwarding data: %w", err)
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
	const permission = fs.FileMode(0o644)
	file, err := os.OpenFile(portForwardPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, permission)
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

var errEmptyToken = errors.New("token received is empty")

func fetchToken(ctx context.Context, client *http.Client,
	username, password string,
) (token string, err error) {
	errSubstitutions := map[string]string{
		url.QueryEscape(username): "<username>",
		url.QueryEscape(password): "<password>",
	}

	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)
	url := url.URL{
		Scheme: "https",
		Host:   "www.privateinternetaccess.com",
		Path:   "/api/client/v2/token",
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return "", replaceInErr(err, errSubstitutions)
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		return "", replaceInErr(err, errSubstitutions)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", makeNOKStatusError(response, errSubstitutions)
	}

	decoder := json.NewDecoder(response.Body)
	var result struct {
		Token string `json:"token"`
	}
	if err := decoder.Decode(&result); err != nil {
		return "", fmt.Errorf("decoding response: %w", err)
	}

	if result.Token == "" {
		return "", errEmptyToken
	}
	return result.Token, nil
}

func fetchPortForwardData(ctx context.Context, client *http.Client, apiIP netip.Addr, token string) (
	port uint16, signature string, expiration time.Time, err error,
) {
	errSubstitutions := map[string]string{url.QueryEscape(token): "<token>"}

	queryParams := make(url.Values)
	queryParams.Add("token", token)
	url := url.URL{
		Scheme:   "https",
		Host:     net.JoinHostPort(apiIP.String(), "19999"),
		Path:     "/getSignature",
		RawQuery: queryParams.Encode(),
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		err = replaceInErr(err, errSubstitutions)
		return 0, "", expiration, fmt.Errorf("obtaining signature payload: %w", err)
	}

	response, err := client.Do(request)
	if err != nil {
		err = replaceInErr(err, errSubstitutions)
		return 0, "", expiration, fmt.Errorf("obtaining signature payload: %w", err)
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
		return 0, "", expiration, fmt.Errorf("decoding response: %w", err)
	}

	if data.Status != "OK" {
		return 0, "", expiration, fmt.Errorf("%w: status is: %s", ErrBadResponse, data.Status)
	}

	port, _, expiration, err = unpackPayload(data.Payload)
	if err != nil {
		return 0, "", expiration, fmt.Errorf("unpacking payload data: %w", err)
	}
	return port, data.Signature, expiration, err
}

var ErrBadResponse = errors.New("bad response received")

func bindPort(ctx context.Context, client *http.Client, apiIPAddress netip.Addr, data piaPortForwardData) (err error) {
	payload, err := packPayload(data.Port, data.Token, data.Expiration)
	if err != nil {
		return fmt.Errorf("serializing payload: %w", err)
	}

	queryParams := make(url.Values)
	queryParams.Add("payload", payload)
	queryParams.Add("signature", data.Signature)
	bindPortURL := url.URL{
		Scheme:   "https",
		Host:     net.JoinHostPort(apiIPAddress.String(), "19999"),
		Path:     "/bindPort",
		RawQuery: queryParams.Encode(),
	}

	errSubstitutions := map[string]string{
		url.QueryEscape(payload):        "<payload>",
		url.QueryEscape(data.Signature): "<signature>",
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, bindPortURL.String(), nil)
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
		return fmt.Errorf("decoding response: from %s: %w", bindPortURL.String(), err)
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

	return fmt.Errorf("%w: %s: %d %s: response received: %s",
		ErrHTTPStatusCodeNotOK, url, response.StatusCode,
		response.Status, shortenMessage)
}
