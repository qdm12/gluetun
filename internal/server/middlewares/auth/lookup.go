package auth

import (
	"fmt"

	"golang.org/x/exp/maps"
)

type internalRole struct {
	name    string
	checker authorizationChecker
}

func settingsToLookupMap(settings Settings) (routeToRoles map[string][]internalRole, err error) {
	authNameToChecker := make(map[string]authorizationChecker, len(settings.Auths))
	for _, auth := range settings.Auths {
		switch auth.Method {
		case MethodNone:
			authNameToChecker[auth.Name] = newNoneMethod()
		default:
			return nil, fmt.Errorf("%w: %s", ErrMethodNotSupported, auth.Name)
		}
	}

	routeToRoles = make(map[string][]internalRole)
	for _, role := range settings.Roles {
		for _, authName := range role.Auths {
			checker, ok := authNameToChecker[authName]
			if !ok {
				return nil, fmt.Errorf("%w: %s is not one of %s", ErrAuthNameNotDefined,
					authName, orStrings(maps.Keys(authNameToChecker)))
			}

			iRole := internalRole{
				name:    role.Name,
				checker: checker,
			}
			for _, route := range role.Routes {
				checkerExists := false
				for _, role := range routeToRoles[route] {
					if role.checker.equal(iRole.checker) {
						checkerExists = true
						break
					}
				}
				if checkerExists {
					// even if the role name is different, if the checker is the same, skip it.
					continue
				}
				routeToRoles[route] = append(routeToRoles[route], iRole)
			}
		}
	}
	return routeToRoles, nil
}
