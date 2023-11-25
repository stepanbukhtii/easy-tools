package api

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/stepanbukhtii/easy-tools/crypto"
)

type ClaimsWithRoles struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

type JWTGenerator struct {
	issuer     string
	privateKey *rsa.PrivateKey
	ttl        time.Duration
}

func NewJWTGenerator(issuer, privateKeyString string, ttl time.Duration) (*JWTGenerator, error) {
	privateKeyPEM := fmt.Sprintf("-----BEGIN RSA PRIVATE KEY-----\n%s\n-----END RSA PRIVATE KEY-----", privateKeyString)
	privateKey, err := crypto.DecodeRSAPrivateKey(privateKeyPEM)
	if err != nil {
		return nil, err
	}

	return &JWTGenerator{
		issuer:     issuer,
		privateKey: privateKey,
		ttl:        ttl,
	}, nil
}

func (g JWTGenerator) GenerateToken(userID string) (string, error) {
	now := time.Now().UTC()
	claims := &jwt.RegisteredClaims{
		Issuer:    g.issuer,
		Subject:   userID,
		Audience:  jwt.ClaimStrings{userID},
		ExpiresAt: jwt.NewNumericDate(now.Add(g.ttl)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return t.SignedString(g.privateKey)
}

func (g JWTGenerator) GenerateTokenWthRoles(userID string, roles []string) (string, error) {
	now := time.Now().UTC()
	claims := &ClaimsWithRoles{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    g.issuer,
			Subject:   userID,
			Audience:  jwt.ClaimStrings{userID},
			ExpiresAt: jwt.NewNumericDate(now.Add(g.ttl)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Roles: roles,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return t.SignedString(g.privateKey)
}
