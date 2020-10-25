package provider

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
)

type piaV3 struct {
	servers    []models.PIAOldServer
	randSource rand.Source
}

func newPrivateInternetAccessV3(servers []models.PIAOldServer, timeNow timeNowFunc) *piaV3 {
	return &piaV3{
		servers:    servers,
		randSource: rand.NewSource(timeNow().UnixNano()),
	}
}

func (p *piaV3) GetOpenVPNConnection(selection models.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
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
		return connection, fmt.Errorf(
			"combination of protocol %q and encryption %q does not yield any port number",
			selection.Protocol, selection.EncryptionPreset)
	}

	if selection.TargetIP != nil {
		return models.OpenVPNConnection{IP: selection.TargetIP, Port: port, Protocol: selection.Protocol}, nil
	}

	servers := filterPIAOldServers(p.servers, selection.Regions)
	if len(servers) == 0 {
		return connection, fmt.Errorf("no server found for regions %s", commaJoin(selection.Regions))
	}

	var connections []models.OpenVPNConnection
	for _, server := range servers {
		for _, IP := range server.IPs {
			connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: selection.Protocol})
		}
	}

	return pickRandomConnection(connections, p.randSource), nil
}

func (p *piaV3) BuildConf(connection models.OpenVPNConnection, verbosity, uid, gid int,
	root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string) {
	return buildPIAConf(connection, verbosity, root, cipher, auth, extras)
}

func (p *piaV3) PortForward(ctx context.Context, client *http.Client,
	fileManager files.FileManager, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath models.Filepath)) {
	const uuidLength = 32
	b := make([]byte, uuidLength)
	n, err := rand.New(p.randSource).Read(b) //nolint:gosec
	if err != nil {
		pfLogger.Error(err)
		return
	} else if n != uuidLength {
		pfLogger.Error("only read %d bytes instead of %d", n, uuidLength)
		return
	}
	clientID := hex.EncodeToString(b)
	url := fmt.Sprintf("%s/?client_id=%s", constants.PIAPortForwardURL, clientID)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		pfLogger.Error(err)
		return
	}
	response, err := client.Do(request)
	if err != nil {
		pfLogger.Error(err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		pfLogger.Error("%s for %s; does your PIA server support port forwarding?", response.Status, url)
		return
	}
	b, err = ioutil.ReadAll(response.Body)
	if err != nil {
		pfLogger.Error(err)
		return
	} else if len(b) == 0 {
		pfLogger.Error("port forwarding is already activated on this connection, has expired, or you are not connected to a PIA region that supports port forwarding") //nolint:lll
		return
	}
	body := struct {
		Port uint16 `json:"port"`
	}{}
	if err := json.Unmarshal(b, &body); err != nil {
		pfLogger.Error("port forwarding response: %s", err)
		return
	}
	port := body.Port

	filepath := syncState(port)
	pfLogger.Info("Writing port to %s", filepath)
	if err := fileManager.WriteToFile(
		string(filepath), []byte(fmt.Sprintf("%d", port)),
		files.Permissions(constants.AllReadWritePermissions),
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

func filterPIAOldServers(servers []models.PIAOldServer, regions []string) (filtered []models.PIAOldServer) {
	for _, server := range servers {
		switch {
		case filterByPossibilities(server.Region, regions):
		default:
			filtered = append(filtered, server)
		}
	}
	return filtered
}
