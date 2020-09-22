package provider

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/crypto/random"
)

type piaV3 struct {
	random  random.Random
	servers []models.PIAServer
}

func newPrivateInternetAccessV3(servers []models.PIAServer) *piaV3 {
	return &piaV3{
		random:  random.NewRandom(),
		servers: servers,
	}
}

func (p *piaV3) GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error) {
	return getPIAOpenVPNConnections(p.servers, selection)
}

func (p *piaV3) BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string) {
	return buildPIAConf(connections, verbosity, root, cipher, auth, extras)
}

func (p *piaV3) GetPortForward(client *http.Client) (port uint16, err error) {
	b, err := p.random.GenerateRandomBytes(32)
	if err != nil {
		return 0, err
	}
	clientID := hex.EncodeToString(b)
	url := fmt.Sprintf("%s/?client_id=%s", constants.PIAPortForwardURL, clientID)
	response, err := client.Get(url) // TODO add ctx
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("%s for %s; does your PIA server support port forwarding?", response.Status, url)
	}
	b, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	} else if len(b) == 0 {
		return 0, fmt.Errorf("port forwarding is already activated on this connection, has expired, or you are not connected to a PIA region that supports port forwarding")
	}
	body := struct {
		Port uint16 `json:"port"`
	}{}
	if err := json.Unmarshal(b, &body); err != nil {
		return 0, fmt.Errorf("port forwarding response: %w", err)
	}
	return body.Port, nil
}
