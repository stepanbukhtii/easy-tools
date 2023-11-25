package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func DecodeRSAPublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, ErrParsePEMBlock
	}

	if block.Type != "RSA PUBLIC KEY" {
		return nil, ErrRSAPublicKeyWrongType
	}

	return x509.ParsePKCS1PublicKey(block.Bytes)
}

func DecodeRSAPrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, ErrParsePEMBlock
	}

	if block.Type != "RSA PRIVATE KEY" {
		return nil, ErrRSAPrivateKeyWrongType
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func ParseRSAPublicKey(publicKeyString string) (*rsa.PublicKey, error) {
	publicKeyPEM := fmt.Sprintf("-----BEGIN RSA PUBLIC KEY-----\n%s\n-----END RSA PUBLIC KEY-----", publicKeyString)
	return DecodeRSAPublicKey(publicKeyPEM)
}

func ParseRSAPrivateKey(privateKeyString string) (*rsa.PrivateKey, error) {
	privateKeyPEM := fmt.Sprintf("-----BEGIN RSA PRIVATE KEY-----\n%s\n-----END RSA PRIVATE KEY-----", privateKeyString)
	return DecodeRSAPrivateKey(privateKeyPEM)
}
