package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

var secretKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

var secretArr []byte
var publicArr []byte

func NewCrypt() (secret []byte, public []byte) {

	// initialize keys
	if secretKey == nil || publicKey == nil {
		secretKey, publicKey = generate()

		x509Secret := x509.MarshalPKCS1PrivateKey(secretKey)
		x509Public, _ := x509.MarshalPKIXPublicKey(publicKey)

		secretArr = x509Secret
		publicArr = x509Public
	}

	return secretArr, publicArr
}

func Encrypt(data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
}

func Decrypt(data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, secretKey, data)
}

func generate() (secretKey *rsa.PrivateKey, publicKey *rsa.PublicKey) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}

	secretKey = key
	publicKey = &key.PublicKey

	secretKey.Precompute()
	if err := secretKey.Validate(); err != nil {
		panic(err)
	}

	return
}
