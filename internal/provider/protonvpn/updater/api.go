package updater

import (
	"bytes"
	"context"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/netip"
	"slices"
	"strings"

	srp "github.com/ProtonMail/go-srp"
)

// apiClient is a minimal Proton API client using the legacy api.protonvpn.ch
// endpoint which supports direct SRP auth without session/CAPTCHA requirements.
type apiClient struct {
	apiURLBase string
	httpClient *http.Client
	appVersion string
	userAgent  string
	generator  *rand.ChaCha8
}

// newAPIClient returns an [apiClient] with sane defaults.
func newAPIClient(ctx context.Context, httpClient *http.Client) (client *apiClient, err error) {
	var seed [32]byte
	_, _ = crand.Read(seed[:])
	generator := rand.NewChaCha8(seed)

	userAgents := [...]string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:143.0) Gecko/20100101 Firefox/143.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:143.0) Gecko/20100101 Firefox/143.0",
		"Mozilla/5.0 (X11; Linux x86_64; rv:143.0) Gecko/20100101 Firefox/143.0",
	}
	userAgent := userAgents[generator.Uint64()%uint64(len(userAgents))]

	appVersion, err := getMostRecentStableTag(ctx, httpClient)
	if err != nil {
		return nil, fmt.Errorf("getting most recent version for proton app: %w", err)
	}

	return &apiClient{
		apiURLBase: "https://api.protonvpn.ch",
		httpClient: httpClient,
		appVersion: appVersion,
		userAgent:  userAgent,
		generator:  generator,
	}, nil
}

var ErrCodeNotSuccess = errors.New("response code is not success")

// setHeaders sets the minimal necessary headers for Proton API requests.
func (c *apiClient) setHeaders(request *http.Request, authToken string) {
	request.Header.Set("User-Agent", c.userAgent)
	request.Header.Set("x-pm-appversion", c.appVersion)
	request.Header.Set("x-pm-locale", "en_US")
	if authToken != "" {
		request.Header.Set("Authorization", "Bearer "+authToken)
	}
}

// authenticate performs direct SRP authentication against the legacy Proton VPN API
// and returns a bearer token for subsequent requests.
func (c *apiClient) authenticate(ctx context.Context, email, password string,
) (authCookie cookie, err error) {
	// Step 1: Get SRP auth info (no session needed with legacy API)
	username, modulusPGPClearSigned, serverEphemeralBase64, saltBase64,
		srpSessionHex, version, err := c.authInfo(ctx, email)
	if err != nil {
		return cookie{}, fmt.Errorf("getting auth information: %w", err)
	}

	// Step 2: Prepare SRP proof
	srpAuth, err := srp.NewAuth(version, username, []byte(password),
		saltBase64, modulusPGPClearSigned, serverEphemeralBase64)
	if err != nil {
		return cookie{}, fmt.Errorf("initializing SRP auth: %w", err)
	}

	const modulusBits = 2048
	proofs, err := srpAuth.GenerateProofs(modulusBits)
	if err != nil {
		return cookie{}, fmt.Errorf("generating SRP proofs: %w", err)
	}

	// Step 3: Submit SRP proof and get access token
	authCookie, err = c.auth(ctx, email, srpSessionHex, proofs)
	if err != nil {
		return cookie{}, fmt.Errorf("authentifying: %w", err)
	}

	return authCookie, nil
}

var ErrUsernameDoesNotExist = errors.New("username does not exist")

// authInfo fetches SRP parameters for the account using the legacy API.
func (c *apiClient) authInfo(ctx context.Context, email string) (
	username, modulusPGPClearSigned, serverEphemeralBase64, saltBase64, srpSessionHex string,
	version int, err error,
) {
	type requestBodySchema struct {
		Username string `json:"Username"`
	}
	requestBody := requestBodySchema{
		Username: email,
	}

	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(requestBody); err != nil {
		return "", "", "", "", "", 0, fmt.Errorf("encoding request body: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURLBase+"/auth/info", buffer)
	if err != nil {
		return "", "", "", "", "", 0, fmt.Errorf("creating request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	c.setHeaders(request, "")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", "", "", "", "", 0, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", "", "", "", 0, fmt.Errorf("reading response body: %w", err)
	} else if response.StatusCode != http.StatusOK {
		return "", "", "", "", "", 0, buildError(response.StatusCode, responseBody)
	}

	var info struct {
		Code            uint   `json:"Code"`
		Modulus         string `json:"Modulus"`
		ServerEphemeral string `json:"ServerEphemeral"`
		Version         *uint  `json:"Version,omitempty"`
		Salt            string `json:"Salt"`
		SRPSession      string `json:"SRPSession"`
		Username        string `json:"Username"`
	}
	err = json.Unmarshal(responseBody, &info)
	if err != nil {
		return "", "", "", "", "", 0, fmt.Errorf("decoding response body: %w", err)
	}

	const successCode = 1000
	switch {
	case info.Code != successCode:
		return "", "", "", "", "", 0, fmt.Errorf("%w: expected %d got %d",
			ErrCodeNotSuccess, successCode, info.Code)
	case info.Modulus == "":
		return "", "", "", "", "", 0, fmt.Errorf("%w: modulus is empty", ErrDataFieldMissing)
	case info.ServerEphemeral == "":
		return "", "", "", "", "", 0, fmt.Errorf("%w: server ephemeral is empty", ErrDataFieldMissing)
	case info.Salt == "":
		return "", "", "", "", "", 0, fmt.Errorf("%w (salt data field is empty)", ErrUsernameDoesNotExist)
	case info.SRPSession == "":
		return "", "", "", "", "", 0, fmt.Errorf("%w: SRP session is empty", ErrDataFieldMissing)
	case info.Version == nil:
		return "", "", "", "", "", 0, fmt.Errorf("%w: version is missing", ErrDataFieldMissing)
	}

	// Username may be nil/empty in legacy API response, use email as fallback
	username = info.Username
	if username == "" {
		username = email
	}

	version = int(*info.Version) //nolint:gosec
	return username, info.Modulus, info.ServerEphemeral, info.Salt,
		info.SRPSession, version, nil
}

var ErrDataFieldMissing = errors.New("data field missing in response")

type cookie struct {
	uid   string
	token string
}

var (
	ErrServerProofNotValid = errors.New("server proof from server is not valid")
	ErrVPNScopeNotFound    = errors.New("VPN scope not found in scopes")
	ErrTwoFANotSupported   = errors.New("two factor authentication not supported in this client")
)

// auth performs the SRP proof submission to obtain an access token.
func (c *apiClient) auth(ctx context.Context,
	username, srpSession string, proofs *srp.Proofs,
) (authCookie cookie, err error) {
	clientEphemeral := base64.StdEncoding.EncodeToString(proofs.ClientEphemeral)
	clientProof := base64.StdEncoding.EncodeToString(proofs.ClientProof)

	type requestBodySchema struct {
		ClientEphemeral string `json:"ClientEphemeral"`
		ClientProof     string `json:"ClientProof"`
		SRPSession      string `json:"SRPSession"`
		Username        string `json:"Username"`
	}
	requestBody := requestBodySchema{
		ClientEphemeral: clientEphemeral,
		ClientProof:     clientProof,
		SRPSession:      srpSession,
		Username:        username,
	}

	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(requestBody); err != nil {
		return cookie{}, fmt.Errorf("encoding request body: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURLBase+"/auth", buffer)
	if err != nil {
		return cookie{}, fmt.Errorf("creating request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	c.setHeaders(request, "")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return cookie{}, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return cookie{}, fmt.Errorf("reading response body: %w", err)
	} else if response.StatusCode != http.StatusOK {
		return cookie{}, buildError(response.StatusCode, responseBody)
	}

	var auth struct {
		Code         uint     `json:"Code"`
		UID          string   `json:"UID"`
		AccessToken  string   `json:"AccessToken"`
		TokenType    string   `json:"TokenType"`
		Scopes       []string `json:"Scopes"`
		ServerProof  string   `json:"ServerProof"`
		TwoFactor    uint     `json:"TwoFactor"`
	}

	err = json.Unmarshal(responseBody, &auth)
	if err != nil {
		return cookie{}, fmt.Errorf("decoding response body: %w", err)
	}

	m2, err := base64.StdEncoding.DecodeString(auth.ServerProof)
	if err != nil {
		return cookie{}, fmt.Errorf("decoding server proof: %w", err)
	}
	if !bytes.Equal(m2, proofs.ExpectedServerProof) {
		return cookie{}, fmt.Errorf("%w: expected %x got %x",
			ErrServerProofNotValid, proofs.ExpectedServerProof, m2)
	}

	const successCode = 1000
	switch {
	case auth.Code != successCode:
		return cookie{}, fmt.Errorf("%w: expected %d got %d",
			ErrCodeNotSuccess, successCode, auth.Code)
	case auth.TwoFactor != 0:
		return cookie{}, fmt.Errorf("%w", ErrTwoFANotSupported)
	case !slices.Contains(auth.Scopes, "vpn"):
		return cookie{}, fmt.Errorf("%w: in %v", ErrVPNScopeNotFound, auth.Scopes)
	}

	return cookie{
		uid:   auth.UID,
		token: auth.AccessToken,
	}, nil
}

func httpHeadersToString(headers http.Header) string {
	var builder strings.Builder
	first := true
	for key, values := range headers {
		for _, value := range values {
			if !first {
				builder.WriteString(", ")
			}
			builder.WriteString(fmt.Sprintf("%s: %s", key, value))
			first = false
		}
	}
	return builder.String()
}

type apiData struct {
	LogicalServers []logicalServer `json:"LogicalServers"`
}

type logicalServer struct {
	Name        string           `json:"Name"`
	ExitCountry string           `json:"ExitCountry"`
	Region      *string          `json:"Region"`
	City        *string          `json:"City"`
	Servers     []physicalServer `json:"Servers"`
	Features    uint16           `json:"Features"`
	Tier        *uint8           `json:"Tier,omitempty"`
}

type physicalServer struct {
	EntryIP         netip.Addr `json:"EntryIP"`
	ExitIP          netip.Addr `json:"ExitIP"`
	Domain          string     `json:"Domain"`
	Status          uint8      `json:"Status"`
	X25519PublicKey string     `json:"X25519PublicKey"`
}

func (c *apiClient) fetchServers(ctx context.Context, authCookie cookie) (
	data apiData, err error,
) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiURLBase+"/vpn/logicals", nil)
	if err != nil {
		return data, err
	}
	c.setHeaders(request, authCookie.token)
	request.Header.Set("x-pm-uid", authCookie.uid)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return data, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(response.Body)
		return data, buildError(response.StatusCode, b)
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&data); err != nil {
		return data, fmt.Errorf("decoding response body: %w", err)
	}

	return data, nil
}

var ErrHTTPStatusCodeNotOK = errors.New("HTTP status code not OK")

func buildError(httpCode int, body []byte) error {
	prettyCode := http.StatusText(httpCode)
	var protonError struct {
		Code    *int              `json:"Code,omitempty"`
		Error   *string           `json:"Error,omitempty"`
		Details map[string]string `json:"Details"`
	}
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&protonError)
	if err != nil || protonError.Error == nil || protonError.Code == nil {
		return fmt.Errorf("%w: %s: %s",
			ErrHTTPStatusCodeNotOK, prettyCode, body)
	}

	details := make([]string, 0, len(protonError.Details))
	for key, value := range protonError.Details {
		details = append(details, fmt.Sprintf("%s: %s", key, value))
	}

	return fmt.Errorf("%w: %s: %s (code %d with details: %s)",
		ErrHTTPStatusCodeNotOK, prettyCode, *protonError.Error, *protonError.Code, strings.Join(details, ", "))
}
