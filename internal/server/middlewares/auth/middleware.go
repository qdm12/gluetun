package auth

import (
	"fmt"
	"net/http"
)

func New(settings Settings, debugLogger DebugLogger) (
	middleware func(http.Handler) http.Handler,
	err error,
) {
	routeToRoles, err := settingsToLookupMap(settings)
	if err != nil {
		return nil, fmt.Errorf("converting settings to lookup maps: %w", err)
	}

	return func(handler http.Handler) http.Handler {
		return &authHandler{
			childHandler: handler,
			routeToRoles: routeToRoles,
			unprotectedRoutes: map[string]struct{}{
				http.MethodGet + " /openvpn/actions/restart": {},
				http.MethodGet + " /unbound/actions/restart": {},
				http.MethodGet + " /updater/restart":         {},
				http.MethodGet + " /v1/version":              {},
				http.MethodGet + " /v1/vpn/status":           {},
				http.MethodPut + " /v1/vpn/status":           {},
				// GET /v1/vpn/settings is protected by default
				// PUT /v1/vpn/settings is protected by default
				http.MethodGet + " /v1/openvpn/status":        {},
				http.MethodPut + " /v1/openvpn/status":        {},
				http.MethodGet + " /v1/openvpn/portforwarded": {},
				// GET /v1/openvpn/settings is protected by default
				http.MethodGet + " /v1/dns/status":     {},
				http.MethodPut + " /v1/dns/status":     {},
				http.MethodGet + " /v1/updater/status": {},
				http.MethodPut + " /v1/updater/status": {},
				http.MethodGet + " /v1/publicip/ip":    {},
			},
			logger: debugLogger,
		}
	}, nil
}

type authHandler struct {
	childHandler      http.Handler
	routeToRoles      map[string][]internalRole
	unprotectedRoutes map[string]struct{} // TODO v3.41.0 remove
	logger            DebugLogger
}

func (h *authHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	route := request.Method + " " + request.URL.Path
	roles := h.routeToRoles[route]
	if len(roles) == 0 {
		h.logger.Debugf("no authentication role defined for route %s", route)
		http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	responseHeader := make(http.Header, 0)
	for _, role := range roles {
		if !role.checker.isAuthorized(responseHeader, request) {
			continue
		}

		h.warnIfUnprotectedByDefault(role, route) // TODO v3.41.0 remove

		h.logger.Debugf("access to route %s authorized for role %s", route, role.name)
		h.childHandler.ServeHTTP(writer, request)
		return
	}

	// Flush out response headers if all roles failed to authenticate
	for headerKey, headerValues := range responseHeader {
		for _, headerValue := range headerValues {
			writer.Header().Add(headerKey, headerValue)
		}
	}

	allRoleNames := make([]string, len(roles))
	for i, role := range roles {
		allRoleNames[i] = role.name
	}
	h.logger.Debugf("access to route %s unauthorized after checking for roles %s",
		route, andStrings(allRoleNames))
	http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func (h *authHandler) warnIfUnprotectedByDefault(role internalRole, route string) {
	// TODO v3.41.0 remove
	if role.name != "public" {
		// custom role name, allow none authentication to be specified
		return
	}
	_, isNoneChecker := role.checker.(*noneMethod)
	if !isNoneChecker {
		// not the none authentication method
		return
	}
	_, isUnprotectedByDefault := h.unprotectedRoutes[route]
	if !isUnprotectedByDefault {
		// route is not unprotected by default, so this is a user decision
		return
	}
	h.logger.Warnf("route %s is unprotected by default, "+
		"please set up authentication following the documentation at "+
		"https://github.com/qdm12/gluetun-wiki/blob/main/setup/advanced/control-server.md#authentication "+
		"since this will become no longer publicly accessible after release v3.40.",
		route)
}
