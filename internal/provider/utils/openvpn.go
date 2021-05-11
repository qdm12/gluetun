package utils

func WrapOpenvpnCA(certificate string) (lines []string) {
	return []string{
		"<ca>",
		"-----BEGIN CERTIFICATE-----",
		certificate,
		"-----END CERTIFICATE-----",
		"</ca>",
	}
}

func WrapOpenvpnCert(clientCertificate string) (lines []string) {
	return []string{
		"<cert>",
		"-----BEGIN CERTIFICATE-----",
		clientCertificate,
		"-----END CERTIFICATE-----",
		"</cert>",
	}
}

func WrapOpenvpnCRLVerify(x509CRL string) (lines []string) {
	return []string{
		"<crl-verify>",
		"-----BEGIN X509 CRL-----",
		x509CRL,
		"-----END X509 CRL-----",
		"</crl-verify>",
	}
}

func WrapOpenvpnKey(clientKey string) (lines []string) {
	return []string{
		"<key>",
		"-----BEGIN PRIVATE KEY-----",
		clientKey,
		"-----END PRIVATE KEY-----",
		"</key>",
	}
}

func WrapOpenvpnRSAKey(rsaPrivateKey string) (lines []string) {
	return []string{
		"<key>",
		"-----BEGIN RSA PRIVATE KEY-----",
		rsaPrivateKey,
		"-----END RSA PRIVATE KEY-----",
		"</key>",
	}
}

func WrapOpenvpnTLSAuth(staticKeyV1 string) (lines []string) {
	return []string{
		"<tls-auth>",
		"-----BEGIN OpenVPN Static key V1-----",
		staticKeyV1,
		"-----END OpenVPN Static key V1-----",
		"</tls-auth>",
	}
}

func WrapOpenvpnTLSCrypt(staticKeyV1 string) (lines []string) {
	return []string{
		"<tls-crypt>",
		"-----BEGIN OpenVPN Static key V1-----",
		staticKeyV1,
		"-----END OpenVPN Static key V1-----",
		"</tls-crypt>",
	}
}
