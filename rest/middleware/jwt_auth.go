package middleware

import (
	"crypto/ed25519"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"

	"github.com/stepanbukhtii/easy-tools/econtext"
	"github.com/stepanbukhtii/easy-tools/ejwt"
	"github.com/stepanbukhtii/easy-tools/rest/api"
)

var (
	NoAuthorizationHeader = errors.New("no authorization header")
	FailedToParseToken    = errors.New("failed to parse token")
	InvalidToken          = errors.New("invalid token")
)

type JWTAuth struct {
	publicKey      ed25519.PublicKey
	parser         *jwt.Parser
	skipValidation bool
}

func NewJWTAuth(publicKey ed25519.PublicKey, issuer, audience string, skipValidation bool) JWTAuth {
	options := []jwt.ParserOption{
		jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}),
		jwt.WithIssuer(issuer),
		jwt.WithAudience(audience),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
	}

	return JWTAuth{
		publicKey:      publicKey,
		parser:         jwt.NewParser(options...),
		skipValidation: skipValidation,
	}
}

// Auth parse and verifying JWT token
func (m JWTAuth) Auth(c *gin.Context) {
	tokenString, err := request.AuthorizationHeaderExtractor.ExtractToken(c.Request)
	if err != nil {
		api.RespondUnauthorized(c, NoAuthorizationHeader)
		return
	}

	var claims jwt.RegisteredClaims

	if m.skipValidation {
		if _, _, err := m.parser.ParseUnverified(tokenString, &claims); err != nil {
			api.RespondUnauthorized(c, FailedToParseToken)
			return
		}

		if claims.Subject != "" {
			c.Request = c.Request.WithContext(econtext.SetSubject(c.Request.Context(), claims.Subject))
		}

		return
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) { return m.publicKey, nil }
	token, err := m.parser.ParseWithClaims(tokenString, &claims, keyFunc)
	if err != nil {
		api.RespondUnauthorized(c, FailedToParseToken)
		return
	}

	if !token.Valid || claims.Subject == "" {
		api.RespondUnauthorized(c, InvalidToken)
		return
	}

	c.Request = c.Request.WithContext(econtext.SetSubject(c.Request.Context(), claims.Subject))
}

// AuthRole parse and verifying JWT token and check role
func (m JWTAuth) AuthRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := request.AuthorizationHeaderExtractor.ExtractToken(c.Request)
		if err != nil {
			api.RespondUnauthorized(c, NoAuthorizationHeader)
			return
		}

		var claims ejwt.ClaimsWithRoles

		if m.skipValidation {
			if _, _, err := m.parser.ParseUnverified(tokenString, &claims); err != nil {
				fmt.Println(err)
				api.RespondUnauthorized(c, FailedToParseToken)
				return
			}

			if claims.Subject != "" {
				c.Request = c.Request.WithContext(econtext.SetSubject(c.Request.Context(), claims.Subject))
			}

			if claims.Roles != nil {
				c.Request = c.Request.WithContext(econtext.SetRoles(c.Request.Context(), claims.Roles))
			}

			return
		}

		keyFunc := func(token *jwt.Token) (interface{}, error) { return m.publicKey, nil }
		token, err := m.parser.ParseWithClaims(tokenString, &claims, keyFunc)
		if err != nil {
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

		if !token.Valid || claims.Subject == "" || !hasRole {
			api.RespondUnauthorized(c, InvalidToken)
			return
		}

		ctx := c.Request.Context()
		ctx = econtext.SetSubject(ctx, claims.Subject)
		ctx = econtext.SetRoles(ctx, claims.Roles)

		c.Request = c.Request.WithContext(ctx)
	}
}
