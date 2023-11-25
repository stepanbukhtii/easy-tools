package easyjwt

import (
	"crypto/ed25519"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	issuer := "issuer"
	audience := "audience"
	ttl := time.Minute
	userID := "123456789"

	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	assert.NoError(t, err)

	jwtGenerator := NewJWTGenerator(privateKey, issuer, audience, ttl)

	token, err := jwtGenerator.GenerateToken(userID)
	assert.NoError(t, err)

	parser := jwt.NewParser()

	var claims jwt.RegisteredClaims
	parsedToken, err := parser.ParseWithClaims(token, &claims, func(_ *jwt.Token) (any, error) {
		return publicKey, nil
	})

	assert.NotNil(t, parsedToken)
	assert.True(t, parsedToken.Valid)
	assert.Equal(t, issuer, claims.Issuer)
	assert.Equal(t, userID, claims.Subject)
	assert.Equal(t, jwt.ClaimStrings{audience}, claims.Audience)
	assert.Equal(t, time.Now().Add(ttl).Unix(), claims.ExpiresAt.Unix())
	assert.Equal(t, time.Now().Unix(), claims.NotBefore.Unix())
	assert.Equal(t, time.Now().Unix(), claims.IssuedAt.Unix())

	roles := []string{"user", "admin"}
	tokenWithRoles, err := jwtGenerator.GenerateTokenWthRoles(userID, roles)
	assert.NoError(t, err)

	var claimsWithRoles ClaimsWithRoles
	parsedToken, err = parser.ParseWithClaims(tokenWithRoles, &claimsWithRoles, func(_ *jwt.Token) (any, error) {
		return publicKey, nil
	})

	assert.NotNil(t, parsedToken)
	assert.True(t, parsedToken.Valid)
	assert.Equal(t, roles, claimsWithRoles.Roles)
}
