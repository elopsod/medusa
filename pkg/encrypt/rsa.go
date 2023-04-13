package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
)

// Create key-pair
// https://www.scottbrady91.com/OpenSSL/Creating-RSA-Keys-using-OpenSSL

// ReadRsaPrivateKey sets the private key
func ReadRsaPrivateKey(key string) interface{} {
	keyData, err := ioutil.ReadFile(key)
	if err != nil {
		log.Printf("ERROR: fail get idrsa, %s", err.Error())
		os.Exit(1)
	}

	keyBlock, _ := pem.Decode(keyData)
	if keyBlock == nil {
		log.Printf("ERROR: fail get idrsa, invalid key")
		os.Exit(1)
	}

	const (
		// ECPrivateKeyBlockType is a possible value for pem.Block.Type.
		ECPrivateKeyBlockType = "EC PRIVATE KEY"
		// RSAPrivateKeyBlockType is a possible value for pem.Block.Type.
		RSAPrivateKeyBlockType = "RSA PRIVATE KEY"
		// PrivateKeyBlockType is a possible value for pem.Block.Type.
		PrivateKeyBlockType = "PRIVATE KEY"
	)

	switch keyBlock.Type {
	case ECPrivateKeyBlockType:
		// ECDSA Private Key in ASN.1 format
		key, err := x509.ParseECPrivateKey(keyBlock.Bytes)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			os.Exit(1)
		}
		return key
	case RSAPrivateKeyBlockType:
		// RSA Private Key in PKCS#1 format
		key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			os.Exit(1)
		}
		return key
	case PrivateKeyBlockType:
		// RSA or ECDSA Private Key in unencrypted PKCS#8 format
		key, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			os.Exit(1)
		}
		return key
	}
}

// ReadRsaPublicKey sets the public key
func ReadRsaPublicKey(key string) *rsa.PublicKey {
	keyData, err := ioutil.ReadFile(key)
	if err != nil {
		log.Printf("ERROR: fail get idrsapub, %s", err.Error())
		return nil
	}

	keyBlock, _ := pem.Decode(keyData)
	if keyBlock == nil {
		log.Printf("ERROR: fail get idrsapub, invalid key")
		return nil
	}

	publicKey, err := x509.ParsePKIXPublicKey(keyBlock.Bytes)
	if err != nil {
		log.Printf("ERROR: fail get idrsapub, %s", err.Error())
		return nil
	}
	switch publicKey := publicKey.(type) {
	case *rsa.PublicKey:
		return publicKey
	default:
		return nil
	}
}

// RsaEncrypt encrypt data
func RsaEncrypt(payload *[]byte, key *rsa.PublicKey) (string, error) {
	// params
	rnd := rand.Reader
	hash := sha256.New()

	// encrypt with OAEP
	ciperText, err := rsa.EncryptOAEP(hash, rnd, key, *payload, nil)
	if err != nil {
		log.Printf("ERROR: fail to encrypt, %s", err.Error())
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciperText), nil
}

// RsaDecrypt decrypts the data
func RsaDecrypt(payload string, key *rsa.PrivateKey) (*[]byte, error) {
	// decode base64 encoded signature
	msg, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		log.Printf("ERROR: fail to base64 decode, %s", err.Error())
		return nil, err
	}

	// params
	rnd := rand.Reader
	hash := sha256.New()

	// decrypt with OAEP
	plainText, err := rsa.DecryptOAEP(hash, rnd, key, msg, nil)
	if err != nil {
		log.Printf("ERROR: fail to decrypt, %s", err.Error())
		return nil, err
	}

	return &plainText, nil
}
