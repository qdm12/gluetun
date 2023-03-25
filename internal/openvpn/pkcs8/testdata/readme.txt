The key files in this directory are generated using OpenSSL.
Re-generating them is fine and should work with existing tests.

For DES encrypted RSA key files, openssl version 1.x.x is required, and the following commands in order generate the files:

openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:512 -des -pass pass:password -out rsa_pkcs8_aes128cbc_encrypted.pem
openssl pkcs8 -topk8 -in rsa_pkcs8_aes128cbc_encrypted.pem -passin pass:password -nocrypt -out rsa_pkcs8_aes128cbc_decrypted.pem

For AES encrypted RSA key files, the following commands in order generate the files:

openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:512 -aes-128-cbc -pass pass:password -out rsa_pkcs8_descbc_encrypted.pem
openssl pkcs8 -topk8 -in rsa_pkcs8_descbc_encrypted.pem -passin pass:password -nocrypt -out rsa_pkcs8_descbc_decrypted.pem
