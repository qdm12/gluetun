package pkcs8

import (
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/youmark/pkcs8"
)

func Test_UpgradeEncryptedKey(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		encryptedPKCS8PEMKey string
		passphrase           string
		decryptedPKCS8PEMKey string
		errMessage           string
	}{
		"DES-CBC key": {
			encryptedPKCS8PEMKey: DESEncryptedKey,
			passphrase:           "password",
			decryptedPKCS8PEMKey: DESDecryptedKey,
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
			der, err := pkcs8.ConvertPrivateKeyToPKCS8(privateKey)
			require.NoError(t, err)
			pemBlock = &pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: der,
			}
			pemBytes := pem.EncodeToMemory(pemBlock)

			assert.Equal(t, testCase.decryptedPKCS8PEMKey, string(pemBytes))
		})
	}
}
