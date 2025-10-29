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
	"strings"

	srp "github.com/ProtonMail/go-srp"
)

// apiClient is a minimal Proton v4 API client which can handle all the
// oddities of Proton's authentication flow they want to keep hidden
// from the public.
type apiClient struct {
	apiURLBase string
	httpClient *http.Client
	appVersion string
	userAgent  string
	generator  *rand.ChaCha8
}

// newAPIClient returns an [apiClient] with sane defaults matching Proton's
// insane expectations.
func newAPIClient(ctx context.Context, httpClient *http.Client) (client *apiClient, err error) {
	var seed [32]byte
	_, _ = crand.Read(seed[:])
	generator := rand.NewChaCha8(seed)

	// Pick a random user agent from this list. Because I'm not going to tell
	// Proton shit on where all these funny requests are coming from, given their
	// unhelpfulness in figuring out their authentication flow.
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
		apiURLBase: "https://account.proton.me/api",
		httpClient: httpClient,
		appVersion: appVersion,
		userAgent:  userAgent,
		generator:  generator,
	}, nil
}

var (
	ErrCodeNotSuccess      = errors.New("response code is not success")
	ErrHTTPStatusCodeNotOK = errors.New("HTTP status code not OK")
)

// setHeaders sets the minimal necessary headers for Proton API requests
// to succeed without being blocked by their "security" measures.
// See for example [getMostRecentStableTag] on how the app version must
// be set to a recent version or they block your request. "SeCuRiTy"...
func (c *apiClient) setHeaders(request *http.Request, cookie cookie) {
	request.Header.Set("Cookie", cookie.String())
	request.Header.Set("User-Agent", c.userAgent)
	request.Header.Set("x-pm-appversion", c.appVersion)
	request.Header.Set("x-pm-locale", "en_US")
	request.Header.Set("x-pm-uid", cookie.uid)
}

// authenticate performs the full Proton authentication flow
// to obtain an authenticated cookie (uid, token and session ID).
func (c *apiClient) authenticate(ctx context.Context, username, password string,
) (authCookie cookie, err error) {
	sessionID, err := c.getSessionID(ctx)
	if err != nil {
		return cookie{}, fmt.Errorf("getting session ID: %w", err)
	}

	tokenType, accessToken, refreshToken, uid, err := c.getUnauthSession(ctx, sessionID)
	if err != nil {
		return cookie{}, fmt.Errorf("getting unauthenticated session data: %w", err)
	}

	cookieToken, err := c.cookieToken(ctx, sessionID, tokenType, accessToken, refreshToken, uid)
	if err != nil {
		return cookie{}, fmt.Errorf("getting cookie token: %w", err)
	}

	unauthCookie := cookie{
		uid:       uid,
		token:     cookieToken,
		sessionID: sessionID,
	}
	info, err := c.authInfo(ctx, username, unauthCookie)
	if err != nil {
		return cookie{}, fmt.Errorf("getting auth information: %w", err)
	}

	// Prepare SRP proof generator using Proton's official SRP parameters and hashing.
	version := int(info.Version) //nolint:gosec
	srpAuth, err := srp.NewAuth(version, username, []byte(password),
		info.Salt, info.Modulus, info.ServerEphemeral)
	if err != nil {
		return cookie{}, fmt.Errorf("initializing SRP auth: %w", err)
	}

	// Generate SRP proofs (A, M1) with the usual 2048-bit modulus.
	const modulusBits = 2048
	proofs, err := srpAuth.GenerateProofs(modulusBits)
	if err != nil {
		return cookie{}, fmt.Errorf("generating SRP proofs: %w", err)
	}

	authCookie, err = c.auth(ctx, unauthCookie, username, info.SRPSession, proofs)
	if err != nil {
		return cookie{}, fmt.Errorf("authentifying: %w", err)
	}

	return authCookie, nil
}

var ErrSessionIDNotFound = errors.New("session ID not found in cookies")

func (c *apiClient) getSessionID(ctx context.Context) (sessionID string, err error) {
	const url = "https://account.proton.me/vpn"
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	err = response.Body.Close()
	if err != nil {
		return "", fmt.Errorf("closing response body: %w", err)
	}

	for _, cookie := range response.Cookies() {
		if cookie.Name == "Session-Id" {
			return cookie.Value, nil
		}
	}

	return "", fmt.Errorf("%w", ErrSessionIDNotFound)
}

func (c *apiClient) getUnauthSession(ctx context.Context, sessionID string) (
	tokenType, accessToken, refreshToken, uid string, err error,
) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURLBase+"/auth/v4/sessions", nil)
	if err != nil {
		return "", "", "", "", fmt.Errorf("creating request: %w", err)
	}
	unauthCookie := cookie{
		sessionID: sessionID,
	}
	c.setHeaders(request, unauthCookie)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", "", "", "", err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", "", "", fmt.Errorf("reading response body: %w", err)
	} else if response.StatusCode != http.StatusOK {
		// TODO parse JSON and fallback to plain
		return "", "", "", "", fmt.Errorf("%w: %s: %s",
			ErrHTTPStatusCodeNotOK, response.Status, string(responseBody))
	}

	var data sessionsResponse
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return "", "", "", "", fmt.Errorf("decoding response body: %w", err)
	}

	const successCode = 1000
	switch {
	case data.Code != successCode:
		return "", "", "", "", fmt.Errorf("%w: expected %d got %d",
			ErrCodeNotSuccess, successCode, data.Code)
	default: // TODO add more validation
	}

	return data.TokenType, data.AccessToken, data.RefreshToken, data.UID, nil
}

func (c *apiClient) cookieToken(ctx context.Context, sessionID, tokenType, accessToken,
	refreshToken, uid string,
) (cookieToken string, err error) {
	requestBody := cookiesRequest{
		GrantType:    "refresh_token",
		Persistent:   0,
		RedirectURI:  "https://protonmail.com",
		RefreshToken: refreshToken,
		ResponseType: "token",
		State:        generateLettersDigits(c.generator, 24), //nolint:mnd
		UID:          uid,
	}

	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(requestBody); err != nil {
		return "", fmt.Errorf("encoding request body: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURLBase+"/core/v4/auth/cookies", buffer)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	unauthCookie := cookie{
		uid:       uid,
		sessionID: sessionID,
	}
	c.setHeaders(request, unauthCookie)
	request.Header.Set("Authorization", tokenType+" "+accessToken)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	} else if response.StatusCode != http.StatusOK {
		// TODO parse JSON and fallback to plain
		return "", fmt.Errorf("%w: %s: %s",
			ErrHTTPStatusCodeNotOK, response.Status, string(responseBody))
	}

	var cookies cookiesResponse
	err = json.Unmarshal(responseBody, &cookies)
	if err != nil {
		return "", fmt.Errorf("decoding response body: %w", err)
	}

	const successCode = 1000
	switch {
	case cookies.Code != successCode:
		return "", fmt.Errorf("%w: expected %d got %d",
			ErrCodeNotSuccess, successCode, cookies.Code)
	case cookies.UID != requestBody.UID:
		return "", fmt.Errorf("mismatched UID: expected %s got %s",
			requestBody.UID, cookies.UID)
	}

	for _, cookie := range response.Cookies() {
		if cookie.Name == "AUTH-"+uid {
			return cookie.Value, nil
		}
	}

	return "", fmt.Errorf("%w", ErrAuthCookieNotFound)
}

// authInfo fetches SRP parameters for the account (salt, modulus, B, version, session).
func (c *apiClient) authInfo(ctx context.Context, username string, unauthCookie cookie) (
	data authInfoResponse, err error,
) {
	requestBody := authInfoRequest{
		Intent:   "Proton",
		Username: username,
	}

	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(requestBody); err != nil {
		return authInfoResponse{}, fmt.Errorf("encoding request body: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURLBase+"/core/v4/auth/info", buffer)
	if err != nil {
		return authInfoResponse{}, fmt.Errorf("creating request: %w", err)
	}
	c.setHeaders(request, unauthCookie)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return authInfoResponse{}, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return authInfoResponse{}, fmt.Errorf("reading response body: %w", err)
	} else if response.StatusCode != http.StatusOK {
		// TODO parse JSON and fallback to plain
		return authInfoResponse{},
			fmt.Errorf("%w: %s: %s", ErrHTTPStatusCodeNotOK, response.Status, string(responseBody))
	}

	var info authInfoResponse
	err = json.Unmarshal(responseBody, &info)
	if err != nil {
		return authInfoResponse{}, fmt.Errorf("decoding response body: %w", err)
	}

	return info, nil
}

type cookie struct {
	uid       string
	token     string
	sessionID string
}

func (c *cookie) String() string {
	s := ""
	if c.token != "" {
		s += fmt.Sprintf("AUTH-%s=%s; ", c.uid, c.token)
	}
	if c.sessionID != "" {
		s += fmt.Sprintf("Session-Id=%s; ", c.sessionID)
	}
	if c.token != "" {
		s += "Tag=default; iaas=W10; Domain=proton.me; Feature=VPNDashboard:A"
	}
	return s
}

var (
	// ErrServerProofNotValid indicates the M2 from the server didn't match the expected proof.
	ErrServerProofNotValid = errors.New("server proof from server is not valid")
	ErrAuthCookieNotFound  = errors.New("auth cookie not found")
)

// auth performs the SRP proof submission (and optionally TOTP) to obtain tokens.
func (c *apiClient) auth(ctx context.Context, unauthCookie cookie,
	username, srpSession string, proofs *srp.Proofs,
) (authCookie cookie, err error) {
	clientEphemeral := base64.StdEncoding.EncodeToString(proofs.ClientEphemeral)
	clientProof := base64.StdEncoding.EncodeToString(proofs.ClientProof)

	requestBody := authRequest{
		ClientEphemeral: clientEphemeral,
		ClientProof:     clientProof,
		// TODO add payload?
		SRPSession: srpSession,
		Username:   username,
	}

	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(requestBody); err != nil {
		return cookie{}, fmt.Errorf("encoding request body: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURLBase+"/core/v4/auth", buffer)
	if err != nil {
		return cookie{}, fmt.Errorf("creating request: %w", err)
	}
	c.setHeaders(request, unauthCookie)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return cookie{}, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return cookie{}, fmt.Errorf("reading response body: %w", err)
	} else if response.StatusCode != http.StatusOK {
		// TODO parse JSON and fallback to plain
		return cookie{}, fmt.Errorf("%w: %s: %s",
			ErrHTTPStatusCodeNotOK, response.Status, string(responseBody))
	}

	var auth authResponse
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

	const headerKey = "set-cookie"
	setCookieHeaders := response.Header[headerKey]
	if len(setCookieHeaders) == 0 {
		setCookieHeaders = response.Header["Set-Cookie"]
	}
	for _, setCookieHeader := range setCookieHeaders {
		parts := strings.Split(setCookieHeader, ";")
		for _, part := range parts {
			if strings.HasPrefix(part, "AUTH-"+unauthCookie.uid+"=") {
				authCookie = unauthCookie
				authCookie.token = strings.TrimPrefix(part, "AUTH-"+unauthCookie.uid+"=")
				return authCookie, nil
			}
		}
	}

	return cookie{}, fmt.Errorf("%w: in HTTP headers %s",
		ErrAuthCookieNotFound, httpHeadersToString(response.Header))
}

// generateLettersDigits mimicing Proton's own random string generator:
// https://github.com/ProtonMail/WebClients/blob/e4d7e4ab9babe15b79a131960185f9f8275512cd/packages/utils/generateLettersDigits.ts
func generateLettersDigits(rng *rand.ChaCha8, length uint) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	return generateFromCharset(rng, length, charset)
}

func generateLowercaseDigits(rng *rand.ChaCha8, length uint) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	return generateFromCharset(rng, length, charset)
}

func generateFromCharset(rng *rand.ChaCha8, length uint, charset string) string {
	result := make([]byte, length)
	randomBytes := make([]byte, length)
	_, _ = rng.Read(randomBytes)
	for i := range length {
		result[i] = charset[int(randomBytes[i])%len(charset)]
	}
	return string(result)
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

func (c *apiClient) fetchServers(ctx context.Context, cookie cookie) (
	data apiData, err error,
) {
	const url = "https://account.proton.me/api/vpn/logicals"
	// Old Logicals API endpoint: https://api.protonmail.ch/vpn/logicals
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return data, err
	}
	c.setHeaders(request, cookie)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return data, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(response.Body)
		return data, fmt.Errorf("%w: %d %s (%s)", ErrHTTPStatusCodeNotOK,
			response.StatusCode, response.Status, b)
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&data); err != nil {
		return data, fmt.Errorf("decoding response body: %w", err)
	}

	return data, nil
}
