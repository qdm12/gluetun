package constants

import (
	"crypto/md5" //nolint:gosec
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
)

func digestServerModelVersion(t *testing.T, server interface{}, version uint16) string { //nolint:unparam
	bytes, err := json.Marshal(server)
	if err != nil {
		t.Fatal(err)
	}
	bytes = append(bytes, []byte(fmt.Sprintf("%d", version))...)
	arr := md5.Sum(bytes) //nolint:gosec
	return base64.RawStdEncoding.EncodeToString(arr[:])
}

func Test_versions(t *testing.T) {
	t.Parallel()
	allServers := GetAllServers()
	assert.Equal(t, "e8eLGRpb1sNX8mDNPOjA6g", digestServerModelVersion(t, models.CyberghostServer{}, allServers.Cyberghost.Version))
	assert.Equal(t, "Dk7lO6AbS46VdHJKQb5Wxg", digestServerModelVersion(t, models.MullvadServer{}, allServers.Mullvad.Version))
	assert.Equal(t, "fjzfUqJH0KvetGRdZYEtOg", digestServerModelVersion(t, models.NordvpnServer{}, allServers.Nordvpn.Version))
	assert.Equal(t, "gYO+bJZCtQvxVk2dTi5d5Q", digestServerModelVersion(t, models.PIAServer{}, allServers.Pia.Version))
	assert.Equal(t, "EZ/SBXQOCS/iJU7A9yc7vg", digestServerModelVersion(t, models.PurevpnServer{}, allServers.Purevpn.Version))
	assert.Equal(t, "7yfMpHwzRpEngA/6nYsNag", digestServerModelVersion(t, models.SurfsharkServer{}, allServers.Surfshark.Version))
	assert.Equal(t, "7yfMpHwzRpEngA/6nYsNag", digestServerModelVersion(t, models.VyprvpnServer{}, allServers.Vyprvpn.Version))
	assert.Equal(t, "7yfMpHwzRpEngA/6nYsNag", digestServerModelVersion(t, models.WindscribeServer{}, allServers.Windscribe.Version))
}

func digestServersTimestamp(t *testing.T, servers interface{}, timestamp int64) string { //nolint:unparam
	bytes, err := json.Marshal(servers)
	if err != nil {
		t.Fatal(err)
	}
	bytes = append(bytes, []byte(fmt.Sprintf("%d", timestamp))...)
	arr := md5.Sum(bytes) //nolint:gosec
	return base64.RawStdEncoding.EncodeToString(arr[:])
}

func Test_timestamps(t *testing.T) {
	t.Parallel()
	allServers := GetAllServers()
	assert.Equal(t, "lZa+3P5DGuo9VXlsXsW5Jw", digestServersTimestamp(t, allServers.Cyberghost.Servers, allServers.Cyberghost.Timestamp))
	assert.Equal(t, "cK5eeY2KU+doigSAonCfVQ", digestServersTimestamp(t, allServers.Mullvad.Servers, allServers.Mullvad.Timestamp))
	assert.Equal(t, "ZfMT6wXJJBAT0fOqx3TuOA", digestServersTimestamp(t, allServers.Nordvpn.Servers, allServers.Nordvpn.Timestamp))
	assert.Equal(t, "vNjr3aOwoql8lBAKDx8Ygw", digestServersTimestamp(t, allServers.Pia.Servers, allServers.Pia.Timestamp))
	assert.Equal(t, "l3vaGSaIT4YOVGn3IaQkBg", digestServersTimestamp(t, allServers.PiaOld.Servers, allServers.PiaOld.Timestamp))
	assert.Equal(t, "IW1gWNvYTSRDxpAv4kwmzg", digestServersTimestamp(t, allServers.Purevpn.Servers, allServers.Purevpn.Timestamp))
	assert.Equal(t, "f934tXGfEVeNGT3TUdnpxw", digestServersTimestamp(t, allServers.Surfshark.Servers, allServers.Surfshark.Timestamp))
	assert.Equal(t, "wwkmrCGEW06x7ze8+FO2hg", digestServersTimestamp(t, allServers.Vyprvpn.Servers, allServers.Vyprvpn.Timestamp))
	assert.Equal(t, "jT4WjRKNpYojILLJWzGRRw", digestServersTimestamp(t, allServers.Windscribe.Servers, allServers.Windscribe.Timestamp))
}
