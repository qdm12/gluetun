package openvpn

const (
	// AuthConf is the file path to the OpenVPN auth file.
	AuthConf = "/etc/openvpn/auth.conf"
	// AskPassPath is the file path to the decryption passphrase for
	// and encrypted private key, which is pointed by `askpass`.
	AskPassPath = "/etc/openvpn/askpass" //nolint:gosec
)
