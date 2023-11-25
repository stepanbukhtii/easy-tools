package api

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"

	"github.com/stepanbukhtii/easy-tools/api"
	"github.com/stepanbukhtii/easy-tools/crypto"
	"github.com/stepanbukhtii/easy-tools/easycontext"
)

var (
	NoAuthorizationHeader = errors.New("no authorization header")
	InvalidTokenSignature = errors.New("invalid token signature")
	FailedToParseToken    = errors.New("failed to parse token")
	InvalidToken          = errors.New("invalid token")
)

type JWTMiddleware struct {
	publicKey *rsa.PublicKey
	parser    *jwt.Parser
}

func NewMiddleware(publicKeyStr string, enabled bool) (*JWTMiddleware, error) {
	publicKeyPEM := fmt.Sprintf("-----BEGIN RSA PUBLIC KEY-----\n%s\n-----END RSA PUBLIC KEY-----", publicKeyStr)
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

// Auth parse and verifying JWT token
func (m JWTMiddleware) Auth(c *gin.Context) {
	t, err := request.AuthorizationHeaderExtractor.ExtractToken(c.Request)
	if err != nil {
		api.RespondUnauthorized(c, NoAuthorizationHeader)
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
			api.RespondUnauthorized(c, InvalidTokenSignature)
			return
		}
		api.RespondUnauthorized(c, FailedToParseToken)
		return
	}

	if !token.Valid {
		api.RespondUnauthorized(c, InvalidToken)
		return
	}

	c.Request = c.Request.WithContext(easycontext.SetSubject(c.Request.Context(), claims.Subject))

	c.Next()
}

// AuthRole parse and verifying JWT token and check role
func (m JWTMiddleware) AuthRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		t, err := request.AuthorizationHeaderExtractor.ExtractToken(c.Request)
		if err != nil {
			api.RespondUnauthorized(c, NoAuthorizationHeader)
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
				api.RespondUnauthorized(c, InvalidTokenSignature)
				return
			}
			api.RespondUnauthorized(c, FailedToParseToken)
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
			api.RespondUnauthorized(c, InvalidToken)
			return
		}

		ctx := c.Request.Context()
		ctx = easycontext.SetSubject(ctx, claims.Subject)
		ctx = easycontext.SetRoles(ctx, claims.Roles)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
