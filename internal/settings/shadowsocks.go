package settings

// ShadowSocks contains settings to configure the Shadowsocks server
type ShadowSocks struct {
	Enabled  bool
	Password string
	Log      bool
	Port     int
}
