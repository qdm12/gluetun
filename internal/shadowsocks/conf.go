package shadowsocks

import (
	"encoding/json"
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/files"
)

func (c *configurator) MakeConf(port uint16, password, method, nameserver string, uid, gid int) (err error) {
	c.logger.Info("generating configuration file")
	data := generateConf(port, password, method, nameserver)
	return c.fileManager.WriteToFile(
		string(constants.ShadowsocksConf),
		data,
		files.Ownership(uid, gid),
		files.Permissions(0400))
}

func generateConf(port uint16, password, method, nameserver string) (data []byte) {
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
		Nameserver   *string           `json:"nameserver,omitempty"`
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
		Workers:   2,
		Interface: "tun",
	}
	if len(nameserver) > 0 {
		conf.Nameserver = &nameserver
	}
	data, _ = json.Marshal(conf)
	return data
}
