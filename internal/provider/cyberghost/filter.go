package cyberghost

import (
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

var ErrGroupMismatchesProtocol = errors.New("server group does not match protocol")

func (c *Cyberghost) filterServers(selection configuration.ServerSelection) (
	servers []models.CyberghostServer, err error) {
	if len(selection.Groups) == 0 {
		if selection.TCP {
			selection.Groups = tcpGroupChoices()
		} else {
			selection.Groups = udpGroupChoices()
		}
	}

	// Check each group match the protocol
	groupsCheckFn := groupsAreAllUDP
	if selection.TCP {
		groupsCheckFn = groupsAreAllTCP
	}
	if err := groupsCheckFn(selection.Groups); err != nil {
		return nil, err
	}

	for _, server := range c.servers {
		switch {
		case
			utils.FilterByPossibilities(server.Group, selection.Groups),
			utils.FilterByPossibilities(server.Region, selection.Regions),
			utils.FilterByPossibilities(server.Hostname, selection.Hostnames):
		default:
			servers = append(servers, server)
		}
	}

	if len(servers) == 0 {
		return nil, utils.NoServerFoundError(selection)
	}

	return servers, nil
}

func tcpGroupChoices() (choices []string) {
	const tcp = true
	return groupsForTCP(tcp)
}

func udpGroupChoices() (choices []string) {
	const tcp = false
	return groupsForTCP(tcp)
}

func groupsForTCP(tcp bool) (choices []string) {
	allGroups := constants.CyberghostGroupChoices()
	choices = make([]string, 0, len(allGroups))
	for _, group := range allGroups {
		switch {
		case tcp && groupIsTCP(group):
			choices = append(choices, group)
		case !tcp && !groupIsTCP(group):
			choices = append(choices, group)
		}
	}
	return choices
}

func groupIsTCP(group string) bool {
	return strings.Contains(strings.ToLower(group), "tcp")
}

func groupsAreAllTCP(groups []string) error {
	for _, group := range groups {
		if !groupIsTCP(group) {
			return fmt.Errorf("%w: group %s for protocol TCP",
				ErrGroupMismatchesProtocol, group)
		}
	}
	return nil
}

func groupsAreAllUDP(groups []string) error {
	for _, group := range groups {
		if groupIsTCP(group) {
			return fmt.Errorf("%w: group %s for protocol UDP",
				ErrGroupMismatchesProtocol, group)
		}
	}
	return nil
}
