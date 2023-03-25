package pkcs8

// DESEncryptedKey is a PKCS#8 encrypted key using DES-CBC.
// It was generated using the command with openssl 1.x.x:
// openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:512 -des -out encrypted_des.pem
// and inputting the passphrase "password".
const DESEncryptedKey = `-----BEGIN ENCRYPTED PRIVATE KEY-----
MIIBsTBLBgkqhkiG9w0BBQ0wPjApBgkqhkiG9w0BBQwwHAQIZOGelpnEEnACAggA
MAwGCCqGSIb3DQIJBQAwEQYFKw4DAgcECCvn5W5X5GJIBIIBYNdWJwLVUtG52oB1
yNX7Lvh2aN6C0jsvC9z2TGjfxXHu+4JNJPhRsRwNUrIaZg2zJCH6OzHUYl8EUaYP
TXPo4nRcE/CDJ7KBz4zDPZEt1pVdzDf2tozNNcONqQpQoevlsMV4AtiqnKille/U
xp3Ite89LJbSGVd72poq1YR9FrVx+0xQbw/Kgak08v6ABjfyUbP7fYh0iI4dh27O
ONBsN0L4wDkNsXm4c7d+5xSSClWetjK6C856rlI6Dd74DRrm/a186Wg6+62h2qWU
ROTlT1jVbMUXNRQH9DlH22vimV39IyEchp+Ker2qJu9zfkwSDWbMYgFJnRXfx7Oe
ff2BMu0QAUvbUtx8+cdWufBWHsDPZYyfk+nGHJayPklFOgfa+DGFlmZEmueam8bH
VZLrNgG/KFwgPgUeLl244/9hbHYnll0ZDdv/reYciNWO7WyJq+Qd++3rZOpU6+vc
jRtpvvk=
-----END ENCRYPTED PRIVATE KEY-----`

// DESDecryptedKey is the PKCS#8 RSA key corresponding to `DESEncryptedKey`.
// It was derived using the command with openssl 1.x.x:
// openssl rsa -in encrypted_des.pem -out decrypted_des.pem
// using the passphrase "password".
const DESDecryptedKey = `-----BEGIN RSA PRIVATE KEY-----
MIIBOgIBAAJBAL9IXnkehzQgrdJmfwKaJzlPN7IewBYbf+Dm/RU0R6l9aFbsXYFa
G4bz3cQc2bBpP1d4oT76dm5bBJu6orEP67kCAwEAAQJAUmsOOcXLn8xM2RFMvIRL
TkgxyU+ymFP0/6THe3FxRzeA15/pSnPVGenykm0pZm6dF2F5JxUiQkkpfmyfrgcK
WQIhAOnLtrTv+sg8VKdo0joN13MfG/ATk4gpo/1HlyWTAxfXAiEA0XMIiewTFgbx
/y3D9jeTKXPv+F4F3DJMarYxd54bZu8CIQC7q2GzJju5hewyIcs2/KtoZp1nfl9b
2okfo9rpN3QxKwIgMvfXQBjenCGcighNA4GKoi/AWaQnsOnchqtHZmBnMqkCICWe
Cdx6hlxZd52LO2bChdB5nFNgKMBmdiDBw+53S+BX
-----END RSA PRIVATE KEY-----`

// AESEncryptedKey is a PKCS#8 encrypted key using AES-128-CBC.
// It was generated using:
// openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:512 -aes-128-cbc -out encrypted_aes.pem
// and inputting the passphrase "password".
const AESEncryptedKey = `-----BEGIN ENCRYPTED PRIVATE KEY-----
MIIBvTBXBgkqhkiG9w0BBQ0wSjApBgkqhkiG9w0BBQwwHAQI7M1eZVluel4CAggA
MAwGCCqGSIb3DQIJBQAwHQYJYIZIAWUDBAECBBD/R9Aia+oSqEUHB695qIy+BIIB
YIm2CaFc1lgpTjLtdHfEpdxms8bcluIDg2YkeUWKLBD0FaNFIEfypDYDVicKogTK
ElVIkizoDYcmz0hhLJhtI8sT9h8xXjhSEjCIted9kiswDfVtEBZ3JAptHYercv3P
4kupfA+yhmosmXSSLHVO1kMmekn220Ekjol22sMtsBZLahXeu6WDfRYmcs812JR6
VfXJdNEKZ1bNXkgdlh0TWqejhrV+c4JshGH733N6kjpOXCI12t2Qk1KhO+NQUDsu
1loeuJ/w0awr2ruJeJw9RPrm18U48bIEKWHCpE8URfAcUgb4uKsqSrGoEN/Hh5du
jt7x4MhanycrFyayCDvD/bRe+V0YGOH1TEvvp08Ldg1NpGUrw8gpcQKabh4ydxD6
pbQ8aCMZ2nd+tsgeXcIORh8Pk9PB/fz6iNBVR+SJVxZ1TNsH6ntfCFNpSIo/koGR
VrALj8WGlJdsHJBSe4zPXhw=
-----END ENCRYPTED PRIVATE KEY-----`

// AESDecryptedKey is the PKCS#8 RSA key corresponding to `AESEncryptedKey`.
// It was derived using the command:
// openssl rsa -in encrypted_aes.pem -out decrypted_aes.pem
// using the passphrase "password".
const AESDecryptedKey = `-----BEGIN RSA PRIVATE KEY-----
MIIBOgIBAAJBAPSTw5ILYK+sydv+1bpAEg6K2dKGLlhrG91mK19b4OsLftvqZatS
yMi/oaMzaW6y+7HH0Wvsj/mifQ+5NZjyXTECAwEAAQJBAMCwXgeE4ULm3g5uIpLf
gZpleKFtR7wvfr+ajBdP+s8SMCPLFsYBkcrXcOz0wZBoFHjN4fQd2tAUfOqD9Vy3
e7UCIQD+qCxv62Euw/b59vccnSZnuHP58kyXXwPPGnpeUNX9swIhAPXd+y3smrZk
RSX2eCa8HLl+AjaXOHY7W+Tz7+n+u2+LAiADsv2yQoEO5NnZl7TPPZkpOIy2vMZQ
DJlJkODmLdZt8QIgQiuZ/EQfZ1MZIRxyPcqG2I1HPzX3pipXkwjr2sgJ3f0CIHlU
oVy2eTtFw44aLRkr3EjYBRs8UQGFYaQnhBLEspnK
-----END RSA PRIVATE KEY-----`
