package updater

type sessionsResponse struct {
	Code         uint     `json:"Code"`         // 1000 on success
	AccessToken  string   `json:"AccessToken"`  // 32-chars lowercase and digits
	RefreshToken string   `json:"RefreshToken"` // 32-chars lowercase and digits
	TokenType    string   `json:"TokenType"`    // "Bearer"
	Scopes       []string `json:"Scopes"`       // should be [] for our usage
	UID          string   `json:"UID"`          // 32-chars lowercase and digits
	LocalID      uint     `json:"LocalID"`      // 0 in my case
}

type cookiesRequest struct {
	GrantType    string `json:"GrantType"`    // "refresh_token"
	Persistent   uint   `json:"Persistent"`   // 0
	RedirectURI  string `json:"RedirectURI"`  // "https://protonmail.com"
	RefreshToken string `json:"RefreshToken"` // 32-chars lowercase and digits
	ResponseType string `json:"ResponseType"` // "token"
	State        string `json:"State"`        // 24-chars letters and digits
	UID          string `json:"UID"`          // 32-chars lowercase and digits
}

type cookiesResponse struct {
	Code           uint   `json:"Code"`           // 1000 on success
	UID            string `json:"UID"`            // should match request UID
	LocalID        uint   `json:"LocalID"`        // 0
	RefreshCounter uint   `json:"RefreshCounter"` // 1
}

type authInfoRequest struct {
	Intent   string `json:"Intent"`   // "Proton"
	Username string `json:"Username"` // user@protonmail.com
}

type authInfoResponse struct {
	Code            uint   `json:"Code"`            // 1000 on success
	Modulus         string `json:"Modulus"`         // PGP clearsigned modulus string
	ServerEphemeral string `json:"ServerEphemeral"` // base64
	Version         uint   `json:"Version"`         // 4 as of 2025-10-26
	Salt            string `json:"Salt"`            // base64
	SRPSession      string `json:"SRPSession"`      // hexadecimal
	Username        string `json:"Username"`        // user without @domain.com. Mine has its first letter capitalized.
}

type authRequest struct {
	ClientEphemeral string            `json:"ClientEphemeral"`   // base64(A)
	ClientProof     string            `json:"ClientProof"`       // base64(M1)
	Payload         map[string]string `json:"Payload,omitempty"` // not sure
	SRPSession      string            `json:"SRPSession"`        // hexadecimal
	Username        string            `json:"Username"`          // user@protonmail.com
}

type authResponse struct {
	Code              uint      `json:"Code"`         // 1000 on success
	LocalID           uint      `json:"LocalID"`      // 7 in my case
	Scopes            []string  `json:"Scopes"`       // this should contain "vpn". Same as `Scope` field value.
	UID               string    `json:"UID"`          // same as `Uid` field value
	UserID            string    `json:"UserID"`       // base64
	EventID           string    `json:"EventID"`      // base64
	PasswordMode      uint      `json:"PasswordMode"` // 1 in my case
	ServerProof       string    `json:"ServerProof"`  // base64(M2)
	TwoFactor         uint      `json:"TwoFactor"`    // 0 if 2FA not required
	TwoFA             twoFAInfo `json:"2FA"`
	TemporaryPassword uint      `json:"TemporaryPassword"` // 0 in my case
}

// twoFAInfo describes second factor state/options.
type twoFAInfo struct {
	Enabled twoFAStatus `json:"Enabled"`
	FIDO2   struct {
		AuthenticationOptions any   `json:"AuthenticationOptions"`
		RegisteredKeys        []any `json:"RegisteredKeys"`
	} `json:"FIDO2"`
	TOTP uint `json:"TOTP"`
}

// twoFAStatus represents 2FA state advertised by the API.
type twoFAStatus uint

const (
	twoFADisabled twoFAStatus = iota
	twoFAHasTOTP
	twoFAHasFIDO2
	twoFAHasFIDO2AndTOTP
)
