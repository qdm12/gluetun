package updater

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
)

var ErrHTTPStatusCodeNotOK = errors.New("HTTP status code not OK")

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

// Session Structure for the sessions api endpoint

type protonSession struct {
	Code         int64         `json:"Code"`
	AccessToken  string        `json:"AccessToken"`
	RefreshToken string        `json:"RefreshToken"`
	TokenType    string        `json:"TokenType"`
	Scopes       []interface{} `json:"Scopes"` // This is likely to be []string, however cannot confirm
	UID          string        `json:"UID"`
	LocalID      int64         `json:"LocalID"`
}

func fetchAPI(ctx context.Context, client *http.Client) (
	data apiData, err error,
) {
	var pmSession protonSession

	const TokenType = "Bearer"
	const ProtonAppVer = "web-account@5.0.235.1" // Setting this here incase version needs updating
	const sessionsURL = "https://account.proton.me/api/auth/v4/sessions"

	// Old Logicals API endpoint: https://api.protonmail.ch/vpn/logicals
	// New Logicals API endpoint: https://account.proton.me/api/vpn/v1/logicals with SecureCoreFilter

	const url = "https://account.proton.me/api/vpn/v1/logicals?SecureCoreFilter=all"

	sessionRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, sessionsURL, nil)
	if err != nil {
		return data, err
	}

	// Setup API Request Headers, using app information as the web-account app
	// US locale and force an unauthed session, only the x-pm-appversion header
	// is required the other two headers are optional

	sessionRequest.Header.Set("x-pm-appversion", ProtonAppVer)
	sessionRequest.Header.Set("x-pm-locale", "en_US")
	sessionRequest.Header.Set("x-enforce-unauthsession", "true")

	sessionResponse, err := client.Do(sessionRequest)
	if err != nil {
		return data, err
	}
	defer sessionResponse.Body.Close()

	if sessionResponse.StatusCode != http.StatusOK {
		return data, fmt.Errorf("%w: %d %s", ErrHTTPStatusCodeNotOK,
			sessionResponse.StatusCode, sessionResponse.Status)
	}

	sessionDecoder := json.NewDecoder(sessionResponse.Body)
	if err := sessionDecoder.Decode(&pmSession); err != nil {
		return data, fmt.Errorf("decoding session response body: %w", err)
	}

	// Validate session response has required fields
	switch {
	case pmSession.AccessToken == "":
		return data, fmt.Errorf("session response has no value for the AccessToken field")
	case pmSession.UID == "":
		return data, fmt.Errorf("session response has no value for the UID field")
	}

	if err := sessionResponse.Body.Close(); err != nil {
		return data, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return data, err
	}

	// Setup the auth token from the newly obtained session, the logicals API end
	// point requires in addition to the auth token two custom header entries
	// one specifying the app that made the request and the proton uid attached
	// to the session. If either are missing a HTTP 401 is returned
	request.Header.Set("Authorization", TokenType+" "+pmSession.AccessToken)
	request.Header.Set("x-pm-uid", pmSession.UID)
	request.Header.Set("x-pm-appversion", ProtonAppVer)

	response, err := client.Do(request)
	if err != nil {
		return data, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return data, fmt.Errorf("%w: %d %s", ErrHTTPStatusCodeNotOK,
			response.StatusCode, response.Status)
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&data); err != nil {
		return data, fmt.Errorf("decoding response body: %w", err)
	}

	if err := response.Body.Close(); err != nil {
		return data, err
	}

	return data, nil
}
