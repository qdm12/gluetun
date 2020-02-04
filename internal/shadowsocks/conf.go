package shadowsocks

import (
	"encoding/json"
	"fmt"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func (c *configurator) MakeConf(port uint16, password string) (err error) {
	c.logger.Info("%s: generating configuration file", logPrefix)
	data, err := generateConf(port, password)
	if err != nil {
		return err
	}
	return c.fileManager.WriteToFile(string(constants.ShadowsocksConf), data)
}

func generateConf(port uint16, password string) (data []byte, err error) {
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
		Method:   "chacha20-ietf-poly1305",
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
	return json.Marshal(conf)
}
