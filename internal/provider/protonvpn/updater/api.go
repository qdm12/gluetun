package updater

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"os"

	"github.com/mort666/go-proton-api"
	common "rtlabs.cloud/protonsession"

	"github.com/qdm12/gluetun/internal/constants"
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

func fetchAPI(ctx context.Context, client *http.Client) (
	data apiData, err error,
) {
	var pmSession *common.Session
	var keypass []byte

	username := os.Getenv("PROTON_USERNAME")
	password := os.Getenv("PROTON_PASSWORD")

	const TokenType = "Bearer"
	const AppVersion = "other"
	const ProtonAppVer = "web-account@5.0.235.1" // Setting this here incase version needs updating

	sessionStore := common.NewFileStore(constants.ServersDataPath+"/proton-sessions.db", "default")
	sessionStore.CacheDir = false

	protonOptions := []proton.Option{
		proton.WithAppVersion(AppVersion),
	}

	sessionConfig, err := sessionStore.Load()
	if err != nil {
		if err == common.ErrKeyNotFound {
			pmSession, err = common.SessionFromLogin(ctx, protonOptions, username, password)
			if err != nil {
				return data, err
			}

			keypass, err = common.SaltKeyPass(ctx, pmSession.Client, []byte(password))
			if err != nil {
				return data, err
			}
		} else {
			return data, err
		}
	} else {
		sessionCreds := &common.SessionCredentials{
			UID:          sessionConfig.UID,
			AccessToken:  sessionConfig.AccessToken,
			RefreshToken: sessionConfig.RefreshToken,
		}

		pmSession, err = common.SessionFromRefresh(ctx, protonOptions, sessionCreds)
		if err != nil {
			return data, err
		}
	}

	// Old Logicals API endpoint: https://api.protonmail.ch/vpn/logicals
	// New Logicals API endpoint: https://account.proton.me/api/vpn/logicals

	const url = "https://account.proton.me/api/vpn/logicals"

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return data, err
	}

	// Setup the auth token from the newly obtained session, the logicals API end
	// point requires in addition to the auth token two custom header entries
	// one specifying the app that made the request and the proton uid attached
	// to the session. If either are missing a HTTP 401 is returned
	request.Header.Set("Authorization", TokenType+" "+pmSession.Auth.AccessToken)
	request.Header.Set("x-pm-uid", pmSession.Auth.UID)
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

	config := common.SessionConfig{
		UID:           pmSession.Auth.UID,
		RefreshToken:  pmSession.Auth.RefreshToken,
		AccessToken:   pmSession.Auth.AccessToken,
		SaltedKeyPass: common.Base64Encode(keypass),
	}

	if err := sessionStore.Save(&config); err != nil {
		return data, nil
	}

	return data, nil
}
