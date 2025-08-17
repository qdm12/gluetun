package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Settings_String(t *testing.T) {
	t.Parallel()

	withDefaults := func(s Settings) Settings {
		s.SetDefaults()
		return s
	}

	testCases := map[string]struct {
		settings Settings
		s        string
	}{
		"default settings": {
			settings: withDefaults(Settings{}),
			s: `Settings summary:
├── VPN settings:
|   ├── VPN provider settings:
|   |   ├── Name: private internet access
|   |   └── Server selection settings:
|   |       ├── VPN type: openvpn
|   |       └── OpenVPN server selection settings:
|   |           ├── Protocol: UDP
|   |           └── Private Internet Access encryption preset: strong
|   └── OpenVPN settings:
|       ├── OpenVPN version: 2.6
|       ├── User: [not set]
|       ├── Password: [not set]
|       ├── Private Internet Access encryption preset: strong
|       ├── Network interface: tun0
|       ├── Run OpenVPN as: root
|       └── Verbosity level: 1
├── DNS settings:
|   ├── Keep existing nameserver(s): no
|   ├── DNS server address to use: 127.0.0.1
|   └── DNS over TLS settings:
|       ├── Enabled: yes
|       ├── Update period: every 24h0m0s
|       ├── Upstream resolvers:
|       |   └── Cloudflare
|       ├── Caching: yes
|       ├── IPv6: no
|       └── DNS filtering settings:
|           ├── Block malicious: yes
|           ├── Block ads: no
|           └── Block surveillance: yes
├── Firewall settings:
|   └── Enabled: yes
├── Log settings:
|   └── Log level: INFO
├── IPv6 settings:
|   └── Check address: [2606:4700::6810:84e5]:443
├── Health settings:
|   ├── Server listening address: 127.0.0.1:9999
|   ├── Target address: cloudflare.com:443
|   ├── Duration to wait after success: 5s
|   ├── Read header timeout: 100ms
|   ├── Read timeout: 500ms
|   └── VPN wait durations:
|       ├── Initial duration: 6s
|       └── Additional duration: 5s
├── Shadowsocks server settings:
|   └── Enabled: no
├── HTTP proxy settings:
|   └── Enabled: no
├── Control server settings:
|   ├── Listening address: :8000
|   ├── Logging: yes
|   └── Authentication file path: /gluetun/auth/config.toml
├── Storage settings:
|   └── Filepath: /gluetun/servers.json
├── OS Alpine settings:
|   ├── Process UID: 1000
|   └── Process GID: 1000
├── Public IP settings:
|   ├── IP file path: /tmp/gluetun/ip
|   ├── Public IP data base API: ipinfo
|   └── Public IP data backup APIs:
|       ├── cloudflare
|       ├── ifconfigco
|       └── ip2location
└── Version settings:
    └── Enabled: yes`,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := testCase.settings.String()

			assert.Equal(t, testCase.s, s)
		})
	}
}
