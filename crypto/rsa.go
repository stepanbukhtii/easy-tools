package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var (
	ErrParsePEMBlock          = errors.New("failed to parse PEM block")
	ErrRSAPublicKeyWrongType  = errors.New("RSA public key wrong type")
	ErrRSAPrivateKeyWrongType = errors.New("RSA private key wrong type")
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
