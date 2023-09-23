package privateinternetaccess

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func fetchToken(ctx context.Context, client *http.Client,
	tokenType string, authFilePath string) (token string, err error) {
	username, password, err := getOpenvpnCredentials(authFilePath)
	if err != nil {
		return "", fmt.Errorf("getting username and password: %w", err)
	}

	errSubstitutions := map[string]string{
		url.QueryEscape(username): "<username>",
		url.QueryEscape(password): "<password>",
	}

	var path string

	switch tokenType {
	case "client":
		path = "/api/client/v2/token"
	case "gtoken":
		path = "/gtoken/generateToken"
	default:
		return "", fmt.Errorf("token type %q is not supported", tokenType)
	}

	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)
	url := url.URL{
		Scheme: "https",
		Host:   "www.privateinternetaccess.com",
		Path:   path,
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return "", ReplaceInErr(err, errSubstitutions)
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		return "", ReplaceInErr(err, errSubstitutions)
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
