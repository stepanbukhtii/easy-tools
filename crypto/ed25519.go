package crypto

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func DecodeED25519PublicKey(publicKeyPEM string) (ed25519.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, ErrParsePEMBlock
	}

	if block.Type != "PUBLIC KEY" {
		return nil, ErrRSAPublicKeyWrongType
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	edKey, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, ErrED25519PublicKeyWrongType
	}

	return edKey, err
}

func DecodeED25519PrivateKey(privateKeyPEM string) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, ErrParsePEMBlock
	}

	if block.Type != "PRIVATE KEY" {
		return nil, ErrRSAPrivateKeyWrongType
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ed25519Key, ok := privateKey.(ed25519.PrivateKey)
	if !ok {
		return nil, ErrRSAPrivateKeyWrongType
	}

	return ed25519Key, nil
}

func ParseED25519PublicKey(publicKeyString string) (ed25519.PublicKey, error) {
	publicKeyPEM := fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----", publicKeyString)
	return DecodeED25519PublicKey(publicKeyPEM)
}

func ParseED25519PrivateKey(privateKeyString string) (ed25519.PrivateKey, error) {
	privateKeyPEM := fmt.Sprintf("-----BEGIN PRIVATE KEY-----\n%s\n-----END PRIVATE KEY-----", privateKeyString)
	return DecodeED25519PrivateKey(privateKeyPEM)
}
