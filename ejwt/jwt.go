package ejwt

import (
	"crypto/ed25519"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ClaimsWithRoles struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

type JWTGenerator struct {
	privateKey ed25519.PrivateKey
	issuer     string
	audience   string
	ttl        time.Duration
}

func NewJWTGenerator(privateKey ed25519.PrivateKey, issuer, audience string, ttl time.Duration) JWTGenerator {
	return JWTGenerator{
		privateKey: privateKey,
		issuer:     issuer,
		audience:   audience,
		ttl:        ttl,
	}
}

func (g *JWTGenerator) GenerateToken(subject string) (string, error) {
	now := time.Now().UTC()
	claims := &jwt.RegisteredClaims{
		Issuer:    g.issuer,
		Subject:   subject,
		Audience:  jwt.ClaimStrings{g.audience},
		ExpiresAt: jwt.NewNumericDate(now.Add(g.ttl)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	return t.SignedString(g.privateKey)
}

func (g *JWTGenerator) GenerateTokenWthRoles(userID string, roles []string) (string, error) {
	now := time.Now().UTC()
	claims := &ClaimsWithRoles{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    g.issuer,
			Subject:   userID,
			Audience:  jwt.ClaimStrings{g.audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(g.ttl)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Roles: roles,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	return t.SignedString(g.privateKey)
}
