package crypto

import "errors"

var (
	ErrParsePEMBlock              = errors.New("failed to parse PEM block")
	ErrRSAPublicKeyWrongType      = errors.New("RSA public key wrong type")
	ErrRSAPrivateKeyWrongType     = errors.New("RSA private key wrong type")
	ErrED25519PublicKeyWrongType  = errors.New("ed25519 public key wrong type")
	ErrED25519PrivateKeyWrongType = errors.New("ed25519 private key wrong type")
)
