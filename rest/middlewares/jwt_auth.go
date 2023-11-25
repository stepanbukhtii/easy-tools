package middlewares

import (
	"crypto/ed25519"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"

	"github.com/stepanbukhtii/easy-tools/easycontext"
	"github.com/stepanbukhtii/easy-tools/easyjwt"
	"github.com/stepanbukhtii/easy-tools/rest/api"
)

var (
	NoAuthorizationHeader = errors.New("no authorization header")
	FailedToParseToken    = errors.New("failed to parse token")
	InvalidToken          = errors.New("invalid token")
)

type Middleware struct {
	publicKey      ed25519.PublicKey
	parser         *jwt.Parser
	skipValidation bool
}

func NewMiddleware(publicKey ed25519.PublicKey, issuer, audience string, skipValidation bool) Middleware {
	options := []jwt.ParserOption{
		jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}),
		jwt.WithIssuer(issuer),
		jwt.WithAudience(audience),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
	}

	return Middleware{
		publicKey:      publicKey,
		parser:         jwt.NewParser(options...),
		skipValidation: skipValidation,
	}
}

// Auth parse and verifying JWT token
func (m Middleware) Auth(c *gin.Context) {
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
			c.Request = c.Request.WithContext(easycontext.SetSubject(c.Request.Context(), claims.Subject))
		}

		c.Next()
		return
	}

	token, err := m.parser.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return m.publicKey, nil
	})
	if err != nil {
		api.RespondUnauthorized(c, FailedToParseToken)
		return
	}

	if !token.Valid || claims.Subject == "" {
		api.RespondUnauthorized(c, InvalidToken)
		return
	}

	c.Request = c.Request.WithContext(easycontext.SetSubject(c.Request.Context(), claims.Subject))

	c.Next()
}

// AuthRole parse and verifying JWT token and check role
func (m Middleware) AuthRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := request.AuthorizationHeaderExtractor.ExtractToken(c.Request)
		if err != nil {
			api.RespondUnauthorized(c, NoAuthorizationHeader)
			return
		}

		var claims easyjwt.ClaimsWithRoles

		if m.skipValidation {
			if _, _, err := m.parser.ParseUnverified(tokenString, &claims); err != nil {
				fmt.Println(err)
				api.RespondUnauthorized(c, FailedToParseToken)
				return
			}

			if claims.Subject != "" {
				c.Request = c.Request.WithContext(easycontext.SetSubject(c.Request.Context(), claims.Subject))
			}

			if claims.Roles != nil {
				c.Request = c.Request.WithContext(easycontext.SetRoles(c.Request.Context(), claims.Roles))
			}

			c.Next()
			return
		}

		token, err := m.parser.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return m.publicKey, nil
		})
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
		ctx = easycontext.SetSubject(ctx, claims.Subject)
		ctx = easycontext.SetRoles(ctx, claims.Roles)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
