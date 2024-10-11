package pkcs8

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/youmark/pkcs8"
)

func parsePEMFile(t *testing.T, pemFilepath string) (base64DER string) {
	t.Helper()

	bytes, err := os.ReadFile(pemFilepath)
	require.NoError(t, err)

	pemBlock, _ := pem.Decode(bytes)
	require.NotNil(t, pemBlock)

	derBytes := pemBlock.Bytes
	base64DER = base64.StdEncoding.EncodeToString(derBytes)
	return base64DER
}

func Test_UpgradeEncryptedKey(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		encryptedPKCS8base64DERKey string
		passphrase                 string
		decryptedPKCS8Base64DERKey string
		errMessage                 string
	}{
		"AES-128-CBC key": {
			encryptedPKCS8base64DERKey: parsePEMFile(t, "testdata/rsa_pkcs8_aes128cbc_encrypted.pem"),
			passphrase:                 "password",
			decryptedPKCS8Base64DERKey: parsePEMFile(t, "testdata/rsa_pkcs8_aes128cbc_decrypted.pem"),
		},
		"DES-CBC key": {
			encryptedPKCS8base64DERKey: parsePEMFile(t, "testdata/rsa_pkcs8_descbc_encrypted.pem"),
			passphrase:                 "password",
			decryptedPKCS8Base64DERKey: parsePEMFile(t, "testdata/rsa_pkcs8_descbc_decrypted.pem"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			securelyEncryptedPKCS8DERKey, err := UpgradeEncryptedKey(testCase.encryptedPKCS8base64DERKey, testCase.passphrase)

			if testCase.errMessage != "" {
				assert.EqualError(t, err, testCase.errMessage)
				return
			}
			assert.NoError(t, err)

			// Decrypt possible re-encrypted key to verify it matches the expected
			// corresponding decrypted key.
			der, err := base64.StdEncoding.DecodeString(securelyEncryptedPKCS8DERKey)
			require.NoError(t, err)
			privateKey, err := pkcs8.ParsePKCS8PrivateKey(der, []byte(testCase.passphrase))
			require.NoError(t, err)
			der, err = x509.MarshalPKCS8PrivateKey(privateKey)
			require.NoError(t, err)
			base64DER := base64.StdEncoding.EncodeToString(der)
			assert.Equal(t, testCase.decryptedPKCS8Base64DERKey, base64DER)
		})
	}
}
