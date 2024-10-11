package pkcs8

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pkcs8lib "github.com/youmark/pkcs8"
)

func Test_getEncryptionAlgorithmOid(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		makeDER                   func() (der []byte, err error)
		encryptionSchemeAlgorithm asn1.ObjectIdentifier
		errMessage                string
	}{
		"empty data": {
			makeDER: func() (der []byte, err error) { return nil, nil },
			errMessage: "decoding asn1 encrypted private key data: " +
				"asn1: syntax error: sequence truncated",
		},
		"algorithm not pbes2": {
			makeDER: func() (der []byte, err error) {
				data := encryptedPrivateKey{
					EncryptionAlgorithm: pkix.AlgorithmIdentifier{
						Algorithm: asn1.ObjectIdentifier{1, 2, 3, 4},
					},
				}
				return asn1.Marshal(data)
			},
			errMessage: "encryption algorithm is not PBES2: " +
				"1.2.3.4 instead of PBES2 1.2.840.113549.1.5.13",
		},
		"empty params full bytes": {
			makeDER: func() (der []byte, err error) {
				data := encryptedPrivateKey{
					EncryptionAlgorithm: pkix.AlgorithmIdentifier{
						Algorithm: asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 13},
						Parameters: asn1.RawValue{
							FullBytes: []byte{},
						},
					},
				}
				return asn1.Marshal(data)
			},
			errMessage: "decoding asn1 encryption algorithm parameters: " +
				"asn1: structure error: tags don't match " +
				"(16 vs {class:0 tag:0 length:0 isCompound:false}) {optional:false explicit:false application:false private:false defaultValue:<nil> tag:<nil> stringType:0 timeType:0 set:false omitEmpty:false} encryptedAlgorithmParams @2", //nolint:lll
		},
		"DES-CBC DER": {
			makeDER: func() (der []byte, err error) {
				DESCBCEncryptedPEM, err := os.ReadFile("testdata/rsa_pkcs8_descbc_encrypted.pem")
				if err != nil {
					return nil, fmt.Errorf("reading file: %w", err)
				}
				pemBlock, _ := pem.Decode(DESCBCEncryptedPEM)
				if pemBlock == nil {
					return nil, errors.New("failed to decode PEM")
				}
				return pemBlock.Bytes, nil
			},
			encryptionSchemeAlgorithm: oidDESCBC,
		},
		"AES-128-CBC DER": {
			makeDER: func() (der []byte, err error) {
				AES128CBCEncryptedPEM, err := os.ReadFile("testdata/rsa_pkcs8_aes128cbc_encrypted.pem")
				if err != nil {
					return nil, fmt.Errorf("reading file: %w", err)
				}
				pemBlock, _ := pem.Decode(AES128CBCEncryptedPEM)
				if pemBlock == nil {
					return nil, errors.New("failed to decode PEM")
				}
				return pemBlock.Bytes, nil
			},
			encryptionSchemeAlgorithm: pkcs8lib.AES128CBC.OID(),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			der, err := testCase.makeDER()
			require.NoError(t, err)

			encryptionSchemeAlgorithm, err := getEncryptionAlgorithmOid(der)

			if testCase.errMessage != "" {
				assert.EqualError(t, err, testCase.errMessage)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.encryptionSchemeAlgorithm, encryptionSchemeAlgorithm)
		})
	}
}
