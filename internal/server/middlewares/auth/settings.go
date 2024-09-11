package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/validate"
)

type Settings struct {
	// Roles is a list of roles with their associated authentication
	// and routes.
	Roles []Role
}

func (s *Settings) SetDefaults() {
	s.Roles = gosettings.DefaultSlice(s.Roles, []Role{{ // TODO v3.41.0 leave empty
		Name: "public",
		Auth: "none",
		Routes: []string{
			http.MethodGet + " /openvpn/actions/restart",
			http.MethodGet + " /unbound/actions/restart",
			http.MethodGet + " /updater/restart",
			http.MethodGet + " /v1/version",
			http.MethodGet + " /v1/vpn/status",
			http.MethodPut + " /v1/vpn/status",
			http.MethodGet + " /v1/openvpn/status",
			http.MethodPut + " /v1/openvpn/status",
			http.MethodGet + " /v1/openvpn/portforwarded",
			http.MethodGet + " /v1/dns/status",
			http.MethodPut + " /v1/dns/status",
			http.MethodGet + " /v1/updater/status",
			http.MethodPut + " /v1/updater/status",
			http.MethodGet + " /v1/publicip/ip",
		},
	}})
}

func (s Settings) Validate() (err error) {
	for i, role := range s.Roles {
		err = role.validate()
		if err != nil {
			return fmt.Errorf("role %s (%d of %d): %w",
				role.Name, i+1, len(s.Roles), err)
		}
	}

	return nil
}

const (
	AuthNone   = "none"
	AuthAPIKey = "apikey"
	AuthBasic  = "basic"
)

// Role contains the role name, authentication method name and
// routes that the role can access.
type Role struct {
	// Name is the role name and is only used for documentation
	// and in the authentication middleware debug logs.
	Name string
	// Auth is the authentication method to use, which can be 'none' or 'apikey'.
	Auth string
	// APIKey is the API key to use when using the 'apikey' authentication.
	APIKey string
	// Username for HTTP Basic authentication method.
	Username string
	// Password for HTTP Basic authentication method.
	Password string
	// Routes is a list of routes that the role can access in the format
	// "HTTP_METHOD PATH", for example "GET /v1/vpn/status"
	Routes []string
}

var (
	ErrMethodNotSupported = errors.New("authentication method not supported")
	ErrAPIKeyEmpty        = errors.New("api key is empty")
	ErrBasicUsernameEmpty = errors.New("username is empty")
	ErrBasicPasswordEmpty = errors.New("password is empty")
	ErrRouteNotSupported  = errors.New("route not supported by the control server")
)

func (r Role) validate() (err error) {
	err = validate.IsOneOf(r.Auth, AuthNone, AuthAPIKey, AuthBasic)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrMethodNotSupported, r.Auth)
	}

	switch {
	case r.Auth == AuthAPIKey && r.APIKey == "":
		return fmt.Errorf("for role %s: %w", r.Name, ErrAPIKeyEmpty)
	case r.Auth == AuthBasic && r.Username == "":
		return fmt.Errorf("for role %s: %w", r.Name, ErrBasicUsernameEmpty)
	case r.Auth == AuthBasic && r.Password == "":
		return fmt.Errorf("for role %s: %w", r.Name, ErrBasicPasswordEmpty)
	}

	for i, route := range r.Routes {
		_, ok := validRoutes[route]
		if !ok {
			return fmt.Errorf("route %d of %d: %w: %s",
				i+1, len(r.Routes), ErrRouteNotSupported, route)
		}
	}

	return nil
}

// WARNING: do not mutate programmatically.
var validRoutes = map[string]struct{}{ //nolint:gochecknoglobals
	http.MethodGet + " /openvpn/actions/restart":  {},
	http.MethodGet + " /unbound/actions/restart":  {},
	http.MethodGet + " /updater/restart":          {},
	http.MethodGet + " /v1/version":               {},
	http.MethodGet + " /v1/vpn/status":            {},
	http.MethodPut + " /v1/vpn/status":            {},
	http.MethodGet + " /v1/vpn/settings":          {},
	http.MethodPut + " /v1/vpn/settings":          {},
	http.MethodGet + " /v1/openvpn/status":        {},
	http.MethodPut + " /v1/openvpn/status":        {},
	http.MethodGet + " /v1/openvpn/portforwarded": {},
	http.MethodGet + " /v1/openvpn/settings":      {},
	http.MethodGet + " /v1/dns/status":            {},
	http.MethodPut + " /v1/dns/status":            {},
	http.MethodGet + " /v1/updater/status":        {},
	http.MethodPut + " /v1/updater/status":        {},
	http.MethodGet + " /v1/publicip/ip":           {},
}
