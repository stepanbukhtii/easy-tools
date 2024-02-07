package api

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/stepanbukhtii/easy-tools/crypto"
	"time"
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

type JWTMiddleware struct {
	publicKey *rsa.PublicKey
	parser    *jwt.Parser
}

func NewJWTMiddleware(publicKeyString string, enabled bool) (*JWTMiddleware, error) {
	publicKeyPEM := fmt.Sprintf("-----BEGIN RSA PUBLIC KEY-----\n%s\n-----END RSA PUBLIC KEY-----", publicKeyString)
	publicKey, err := crypto.DecodeRSAPublicKey(publicKeyPEM)
	if err != nil {
		return nil, err
	}

	var options []jwt.ParserOption
	if enabled {
		options = append(options, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}))
	}

	return &JWTMiddleware{
		publicKey: publicKey,
		parser:    jwt.NewParser(options...),
	}, nil
}

// JWTAuth parse and verifying JWT token
func (m JWTMiddleware) JWTAuth(c *gin.Context) {
	t, err := request.AuthorizationHeaderExtractor.ExtractToken(c.Request)
	if err != nil {
		RespondUnauthorized(c, "", NoAuthorizationHeader)
		return
	}

	var claims jwt.RegisteredClaims
	token, err := m.parser.ParseWithClaims(t, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.publicKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			RespondUnauthorized(c, "", InvalidTokenSignature)
			return
		}
		RespondUnauthorized(c, "", FailedToParseToken)
		return
	}

	if !token.Valid {
		RespondUnauthorized(c, "", InvalidToken)
		return
	}

	params := GetParams(c)
	params.Subject = claims.Subject
	c.Set(KeyParams, params)

	c.Next()
}

// JWTAuthRole parse and verifying JWT token and check role
func (m JWTMiddleware) JWTAuthRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		t, err := request.AuthorizationHeaderExtractor.ExtractToken(c.Request)
		if err != nil {
			RespondUnauthorized(c, "", NoAuthorizationHeader)
			return
		}

		// ParseWithClaims
		var claims ClaimsWithRoles
		token, err := m.parser.ParseWithClaims(t, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return m.publicKey, nil
		})
		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				RespondUnauthorized(c, "", InvalidTokenSignature)
				return
			}
			RespondUnauthorized(c, "", FailedToParseToken)
			return
		}

		var hasRole bool
		for _, claimRole := range claims.Roles {
			if claimRole == role {
				hasRole = true
				break
			}
		}

		if !token.Valid || !hasRole {
			RespondUnauthorized(c, "", InvalidToken)
			return
		}

		params := GetParams(c)
		params.Subject = claims.Subject
		params.Roles = claims.Roles
		c.Set(KeyParams, claims.Subject)

		c.Next()
	}
}
