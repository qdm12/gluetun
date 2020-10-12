package provider

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/crypto/random"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
)

type piaV3 struct {
	random  random.Random
	servers []models.PIAOldServer
}

func newPrivateInternetAccessV3(servers []models.PIAOldServer) *piaV3 {
	return &piaV3{
		random:  random.NewRandom(),
		servers: servers,
	}
}

func (p *piaV3) GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error) {
	return getPIAOldOpenVPNConnections(p.servers, selection)
}

func (p *piaV3) BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string) {
	return buildPIAConf(connections, verbosity, root, cipher, auth, extras)
}

func (p *piaV3) PortForward(ctx context.Context, client *http.Client,
	fileManager files.FileManager, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath models.Filepath)) {
	b, err := p.random.GenerateRandomBytes(32)
	if err != nil {
		pfLogger.Error(err)
		return
	}
	clientID := hex.EncodeToString(b)
	url := fmt.Sprintf("%s/?client_id=%s", constants.PIAPortForwardURL, clientID)
	response, err := client.Get(url) // TODO add ctx
	if err != nil {
		pfLogger.Error(err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		pfLogger.Error(fmt.Errorf("%s for %s; does your PIA server support port forwarding?", response.Status, url))
		return
	}
	b, err = ioutil.ReadAll(response.Body)
	if err != nil {
		pfLogger.Error(err)
		return
	} else if len(b) == 0 {
		pfLogger.Error(fmt.Errorf("port forwarding is already activated on this connection, has expired, or you are not connected to a PIA region that supports port forwarding"))
		return
	}
	body := struct {
		Port uint16 `json:"port"`
	}{}
	if err := json.Unmarshal(b, &body); err != nil {
		pfLogger.Error(fmt.Errorf("port forwarding response: %w", err))
		return
	}
	port := body.Port

	filepath := syncState(port)
	pfLogger.Info("Writing port to %s", filepath)
	if err := fileManager.WriteToFile(
		string(filepath), []byte(fmt.Sprintf("%d", port)),
		files.Permissions(0666),
	); err != nil {
		pfLogger.Error(err)
	}

	if err := fw.SetAllowedPort(ctx, port, string(constants.TUN)); err != nil {
		pfLogger.Error(err)
	}

	<-ctx.Done()
	if err := fw.RemoveAllowedPort(ctx, port); err != nil {
		pfLogger.Error(err)
	}
}

func filterPIAOldServers(servers []models.PIAOldServer, region string) (filtered []models.PIAOldServer) {
	if len(region) == 0 {
		return servers
	}
	for _, server := range servers {
		if strings.EqualFold(server.Region, region) {
			return []models.PIAOldServer{server}
		}
	}
	return nil
}

func getPIAOldOpenVPNConnections(allServers []models.PIAOldServer, selection models.ServerSelection) (connections []models.OpenVPNConnection, err error) {
	servers := filterPIAOldServers(allServers, selection.Region)
	if len(servers) == 0 {
		return nil, fmt.Errorf("no server found for region %q", selection.Region)
	}

	var port uint16
	switch selection.Protocol {
	case constants.TCP:
		switch selection.EncryptionPreset {
		case constants.PIAEncryptionPresetNormal:
			port = 502
		case constants.PIAEncryptionPresetStrong:
			port = 501
		}
	case constants.UDP:
		switch selection.EncryptionPreset {
		case constants.PIAEncryptionPresetNormal:
			port = 1198
		case constants.PIAEncryptionPresetStrong:
			port = 1197
		}
	}
	if port == 0 {
		return nil, fmt.Errorf("combination of protocol %q and encryption %q does not yield any port number", selection.Protocol, selection.EncryptionPreset)
	}

	for _, server := range servers {
		for _, IP := range server.IPs {
			if selection.TargetIP != nil {
				if selection.TargetIP.Equal(IP) {
					return []models.OpenVPNConnection{{IP: IP, Port: port, Protocol: selection.Protocol}}, nil
				}
			} else {
				connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: selection.Protocol})
			}
		}
	}

	if selection.TargetIP != nil {
		return nil, fmt.Errorf("target IP %s not found in IP addresses", selection.TargetIP)
	}

	if len(connections) > 64 {
		connections = connections[:64]
	}

	return connections, nil
}
