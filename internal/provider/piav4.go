package provider

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/network"
)

type piaV4 struct {
	servers []models.PIAServer
}

func newPrivateInternetAccessV4(servers []models.PIAServer) *piaV4 {
	return &piaV4{
		servers: servers,
	}
}

func (p *piaV4) GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error) {
	return getPIAOpenVPNConnections(p.servers, selection)
}

func (p *piaV4) BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string) {
	return buildPIAConf(connections, verbosity, root, cipher, auth, extras)
}

func (p *piaV4) GetPortForward(client network.Client) (port uint16, err error) {
	return 0, fmt.Errorf("not implemented")
}
