package api

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	issuer := "issuer"
	ttl := time.Minute
	userID := "123456789"

	key, err := rsa.GenerateKey(rand.Reader, 4096)
	assert.NoError(t, err)

	jwtGenerator := NewJWTGenerator(issuer, key, ttl)
	token, err := jwtGenerator.GenerateToken(userID)
	assert.NoError(t, err)

	parser := jwt.NewParser()

	var claims jwt.RegisteredClaims
	jwtToken, err := parser.ParseWithClaims(token, &claims, func(_ *jwt.Token) (interface{}, error) {
		return &key.PublicKey, nil
	})

	assert.True(t, jwtToken.Valid)
	assert.Equal(t, issuer, claims.Issuer)
	assert.Equal(t, userID, claims.Subject)
	assert.Equal(t, jwt.ClaimStrings{userID}, claims.Audience)
	assert.Equal(t, time.Now().Add(ttl).Unix(), claims.ExpiresAt.Unix())
	assert.Equal(t, time.Now().Unix(), claims.NotBefore.Unix())
	assert.Equal(t, time.Now().Unix(), claims.IssuedAt.Unix())

	roles := []string{"user", "admin"}
	tokenWithRoles, err := jwtGenerator.GenerateTokenWthRoles(userID, roles)
	assert.NoError(t, err)

	var claimsWithRoles ClaimsWithRoles
	jwtToken, err = parser.ParseWithClaims(tokenWithRoles, &claimsWithRoles, func(_ *jwt.Token) (interface{}, error) {
		return &key.PublicKey, nil
	})

	assert.True(t, jwtToken.Valid)
	assert.Equal(t, roles, claimsWithRoles.Roles)
}
