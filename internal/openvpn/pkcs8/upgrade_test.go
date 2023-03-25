package pkcs8

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/youmark/pkcs8"
)

func Test_UpgradeEncryptedKey(t *testing.T) {
	t.Parallel()

	aes128cbcEncryptedKeyPEM, err := os.ReadFile("testdata/rsa_pkcs8_aes128cbc_encrypted.pem")
	require.NoError(t, err)

	aes128cbcDecryptedKeyPEM, err := os.ReadFile("testdata/rsa_pkcs8_aes128cbc_decrypted.pem")
	require.NoError(t, err)

	descbcEncryptedKeyPEM, err := os.ReadFile("testdata/rsa_pkcs8_descbc_encrypted.pem")
	require.NoError(t, err)

	descbcDecryptedKeyPEM, err := os.ReadFile("testdata/rsa_pkcs8_descbc_decrypted.pem")
	require.NoError(t, err)

	testCases := map[string]struct {
		encryptedPKCS8PEMKey string
		passphrase           string
		decryptedPKCS8PEMKey string
		errMessage           string
	}{
		"AES-128-CBC key": {
			encryptedPKCS8PEMKey: string(aes128cbcEncryptedKeyPEM),
			passphrase:           "password",
			decryptedPKCS8PEMKey: string(aes128cbcDecryptedKeyPEM),
		},
		"DES-CBC key": {
			encryptedPKCS8PEMKey: string(descbcEncryptedKeyPEM),
			passphrase:           "password",
			decryptedPKCS8PEMKey: string(descbcDecryptedKeyPEM),
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			securelyEncryptedPKCS8PEMKey, err := UpgradeEncryptedKey(testCase.encryptedPKCS8PEMKey, testCase.passphrase)

			if testCase.errMessage != "" {
				assert.EqualError(t, err, testCase.errMessage)
				return
			}
			assert.NoError(t, err)

			pemBlock, _ := pem.Decode([]byte(securelyEncryptedPKCS8PEMKey))
			require.NotNil(t, pemBlock)
			privateKey, err := pkcs8.ParsePKCS8PrivateKey(pemBlock.Bytes, []byte(testCase.passphrase))
			require.NoError(t, err)
			der, err := x509.MarshalPKCS8PrivateKey(privateKey)
			require.NoError(t, err)
			pemBlock = &pem.Block{
				Type:  "PRIVATE KEY",
				Bytes: der,
			}
			pemBytes := pem.EncodeToMemory(pemBlock)

			assert.Equal(t, testCase.decryptedPKCS8PEMKey, string(pemBytes))
		})
	}
}
