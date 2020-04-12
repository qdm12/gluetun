package shadowsocks

import (
	"encoding/json"
	"fmt"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func (c *configurator) MakeConf(port uint16, password, method string, uid, gid int) (err error) {
	c.logger.Info("generating configuration file")
	data := generateConf(port, password, method)
	return c.fileManager.WriteToFile(
		string(constants.ShadowsocksConf),
		data,
		files.Ownership(uid, gid),
		files.Permissions(0400))
}

func generateConf(port uint16, password, method string) (data []byte) {
	conf := struct {
		Server       string            `json:"server"`
		User         string            `json:"user"`
		Method       string            `json:"method"`
		Timeout      uint              `json:"timeout"`
		FastOpen     bool              `json:"fast_open"`
		Mode         string            `json:"mode"`
		PortPassword map[string]string `json:"port_password"`
		Workers      uint              `json:"workers"`
		Interface    string            `json:"interface"`
		Nameserver   string            `json:"nameserver"`
	}{
		Server:   "0.0.0.0",
		User:     "nonrootuser",
		Method:   method,
		Timeout:  30,
		FastOpen: false,
		Mode:     "tcp_and_udp",
		PortPassword: map[string]string{
			fmt.Sprintf("%d", port): password,
		},
		Workers:    2,
		Interface:  "tun",
		Nameserver: "127.0.0.1",
	}
	data, _ = json.Marshal(conf)
	return data
}
