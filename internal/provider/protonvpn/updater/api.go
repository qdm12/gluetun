package updater

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Define static errors to avoid dynamic error creation
var (
	ErrNoAccessToken = errors.New("session response has no value for the AccessToken field")
	ErrNoUID         = errors.New("session response has no value for the UID field")
	ErrMissingVPNScope = errors.New("missing VPN scope or insufficient permissions")
)

type apiData struct {
	AuthData  authData
	UserAgent string
}

type authData struct {
	AccessToken string
	UID         string
}

func fetchAPI(ctx context.Context, client *http.Client) (
	data apiData, err error) {
	// ProtonVPN API requires a specific user agent
	data.UserAgent = "gluetun"

	// Create an unauthenticated session to get access token with VPN scope
	authData, err := createSession(ctx, client)
	if err != nil {
		return data, fmt.Errorf("creating session: %w", err)
	}

	data.AuthData = authData
	return data, nil
}

// protonSession represents the session response from ProtonVPN API.
type protonSession struct {
	Code        int    `json:"Code"`
	AccessToken string `json:"AccessToken"`
	UID         string `json:"UID"`
	Scope       string `json:"Scope"`
	LocalID     int64  `json:"LocalID"`
}

// sessionRequest structure for creating an unauthenticated session with scope.
type sessionRequest struct {
	ClientSecret string `json:"ClientSecret,omitempty"`
	Payload      string `json:"Payload,omitempty"`
	Scope        string `json:"Scope,omitempty"` // Add scope field to request VPN access.
}

func createSession(ctx context.Context, client *http.Client) (
	authData authData, err error) {
	// Create an unauthenticated session request with VPN scope
	// This is required to access the VPN logicals endpoint
	sessionReq := sessionRequest{
		Scope: "vpn", // Request VPN scope for accessing VPN resources
	}

	requestBody, err := json.Marshal(sessionReq)
	if err != nil {
		return authData, fmt.Errorf("marshaling session request: %w", err)
	}

	// Use the sessions endpoint to create an unauthenticated session with VPN scope
	const sessionURL = "https://api.proton.me/auth/v4/sessions"
	
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, sessionURL, bytes.NewReader(requestBody))
	if err != nil {
		return authData, fmt.Errorf("creating session request: %w", err)
	}

	// Set required headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-pm-appversion", "Other")
	request.Header.Set("x-pm-apiversion", "3")
	request.Header.Set("User-Agent", "gluetun")

	response, err := client.Do(request)
	if err != nil {
		return authData, fmt.Errorf("doing session request: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return authData, fmt.Errorf("reading session response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return authData, fmt.Errorf("HTTP status code %d: %s", response.StatusCode, string(responseBody))
	}

	var pmSession protonSession
	err = json.Unmarshal(responseBody, &pmSession)
	if err != nil {
		return authData, fmt.Errorf("decoding session response body: %w", err)
	}

	// Check for API error response
	if pmSession.Code != 0 {
		// Code 1000 typically means missing scope or permissions
		if pmSession.Code == 1000 {
			return authData, fmt.Errorf("ProtonVPN API error code %d: %w", pmSession.Code, ErrMissingVPNScope)
		}
		return authData, fmt.Errorf("ProtonVPN API error: code %d", pmSession.Code)
	}

	// Validate session response has required fields
	switch {
	case pmSession.AccessToken == "":
		return authData, ErrNoAccessToken
	case pmSession.UID == "":
		return authData, ErrNoUID
	}

	authData.AccessToken = pmSession.AccessToken
	authData.UID = pmSession.UID

	return authData, nil
}

// fetchServers fetches the VPN server list using the session.
func fetchServers(ctx context.Context, client *http.Client, authData authData) (
	servers []Server, err error) {
	// Use versioned API endpoint with proper parameters
	// Using api.proton.me domain to match the session endpoint
	apiURL := "https://api.proton.me/vpn/v1/logicals"
	
	// Add query parameters to match official client
	params := url.Values{}
	params.Set("SecureCoreFilter", "all")
	fullURL := apiURL + "?" + params.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating servers request: %w", err)
	}

	// Set required headers with authentication
	request.Header.Set("Authorization", "Bearer "+authData.AccessToken)
	request.Header.Set("x-pm-uid", authData.UID)
	request.Header.Set("x-pm-appversion", "Other")
	request.Header.Set("x-pm-apiversion", "3")
	request.Header.Set("User-Agent", "gluetun")

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("doing servers request: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading servers response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		// Check if it's a scope/permission error
		if response.StatusCode == http.StatusForbidden || response.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf("access denied (HTTP %d): %w - %s", 
				response.StatusCode, ErrMissingVPNScope, string(responseBody))
		}
		return nil, fmt.Errorf("HTTP status code %d: %s", response.StatusCode, string(responseBody))
	}

	var serverResponse struct {
		Code           int      `json:"Code"`
		LogicalServers []Server `json:"LogicalServers"`
	}

	err = json.Unmarshal(responseBody, &serverResponse)
	if err != nil {
		return nil, fmt.Errorf("decoding servers response: %w", err)
	}

	// Check for API error in response
	if serverResponse.Code != 0 {
		if serverResponse.Code == 1000 {
			return nil, fmt.Errorf("ProtonVPN API error code %d: %w", serverResponse.Code, ErrMissingVPNScope)
		}
		return nil, fmt.Errorf("ProtonVPN API error: code %d", serverResponse.Code)
	}

	return serverResponse.LogicalServers, nil
}

// Server represents a ProtonVPN server.
type Server struct {
	Name     string `json:"Name"`
	EntryIP  string `json:"EntryIP"`
	ExitIP   string `json:"ExitIP"`
	Domain   string `json:"Domain"`
	Tier     int    `json:"Tier"`
	Features int    `json:"Features"`
	Region   string `json:"Region"`
	City     string `json:"City"`
	Score    float64 `json:"Score"`
	Status   int    `json:"Status"`
	// Add other fields as needed
}
