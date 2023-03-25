package pkcs8

import (
	"encoding/pem"
	"errors"
	"fmt"

	pkcs8lib "github.com/youmark/pkcs8"
)

var (
	ErrPEMDecodingFailed      = errors.New("pem decoding failed")
	ErrNotEncryptedPrivateKey = errors.New("not an encrypted private key")
	ErrUnsupportedKeyType     = errors.New("unsupported key type")
)

// UpgradeEncryptedKey eventually upgrades an encrypted key to a newer encryption
// if its encryption is too weak for Openvpn/Openssl.
// If the key is encrypted using DES-CBC, it is decrypted and re-encrypted using AES-256-CBC.
// Otherwise, the key is returned unmodified.
// Note this function only supports:
// - PKCS8 encrypted keys
// - RSA and ECDSA keys
// - DES-CBC, 3DES, AES-128-CBC, AES-192-CBC, AES-256-CBC, AES-128-GCM, AES-192-GCM
// and AES-256-GCM encryption algorithms.
func UpgradeEncryptedKey(encryptedPKCS8PEMKey, passphrase string) (securelyEncryptedPKCS8PEMKey string, err error) {
	pemBlock, _ := pem.Decode([]byte(encryptedPKCS8PEMKey))
	if pemBlock == nil {
		return "", fmt.Errorf("%w", ErrPEMDecodingFailed)
	}

	if pemBlock.Type != "ENCRYPTED PRIVATE KEY" {
		return "", fmt.Errorf("%w: %s", ErrNotEncryptedPrivateKey, pemBlock.Type)
	}

	der := pemBlock.Bytes

	oidEncryptionAlgorithm, err := getEncryptionAlgorithmOid(der)
	if err != nil {
		return "", fmt.Errorf("finding encryption algorithm oid: %w", err)
	}

	if !oidEncryptionAlgorithm.Equal(oidDESCBC) {
		return encryptedPKCS8PEMKey, nil
	}

	// Convert DES-CBC encrypted key to an AES256CBC encrypted key
	privateKey, err := pkcs8lib.ParsePKCS8PrivateKey(der, []byte(passphrase))
	if err != nil {
		return "", fmt.Errorf("parsing pkcs8 encrypted private key: %w", err)
	}

	der, err = pkcs8lib.MarshalPrivateKey(privateKey, []byte(passphrase), pkcs8lib.DefaultOpts)
	if err != nil {
		return "", fmt.Errorf("encrypting and encoding private key: %w", err)
	}

	pemBlock = &pem.Block{
		Type:  "ENCRYPTED PRIVATE KEY",
		Bytes: der,
	}
	encryptedPEMKeyBytes := pem.EncodeToMemory(pemBlock)
	securelyEncryptedPKCS8PEMKey = string(encryptedPEMKeyBytes)
	return securelyEncryptedPKCS8PEMKey, nil
}
