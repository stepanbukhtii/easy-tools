package api

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/rs/zerolog/log"
	"github.com/stepanbukhtii/easy-tools/config"
	"github.com/stepanbukhtii/easy-tools/crypto"
	"time"
)

const (
	KeySubjectJWT = "subject"
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

func NewJWTGenerator(issuer, privateKeyPEM string, ttl time.Duration) (*JWTGenerator, error) {
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

func NewJWTMiddleware(publicKeyPEM string, c config.JWT) (*JWTMiddleware, error) {
	publicKey, err := crypto.DecodeRSAPublicKey(publicKeyPEM)
	if err != nil {
		return nil, err
	}

	var options []jwt.ParserOption
	if c.Enabled {
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
		log.Error().
			Err(err).
			Str(HeaderTraceID, c.Request.Context().Value(HeaderTraceID).(string)).
			Str("ip", c.ClientIP()).
			Msg("Authorization header is empty or malformed")
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
			log.Error().
				Err(err).
				Str(HeaderTraceID, c.Request.Context().Value(HeaderTraceID).(string)).
				Str("ip", c.ClientIP()).
				Msg("Token signature is not valid")
			RespondUnauthorized(c, "", InvalidTokenSignature)
			return
		}
		log.Error().
			Err(err).
			Str(HeaderTraceID, c.Request.Context().Value(HeaderTraceID).(string)).
			Str("ip", c.ClientIP()).
			Msg("Failed to parse jwt token")
		RespondUnauthorized(c, "", FailedToParseToken)
		return
	}

	if !token.Valid {
		log.Error().
			Err(err).
			Str(HeaderTraceID, c.Request.Context().Value(HeaderTraceID).(string)).
			Str("ip", c.ClientIP()).
			Msg("Invalid token")
		RespondUnauthorized(c, "", InvalidToken)
		return
	}

	c.Set(KeySubjectJWT, claims.Subject)

	c.Next()
}

// JWTAuthRole parse and verifying JWT token and check role
func (m JWTMiddleware) JWTAuthRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		t, err := request.AuthorizationHeaderExtractor.ExtractToken(c.Request)
		if err != nil {
			log.Error().
				Err(err).
				Str(HeaderTraceID, c.Request.Context().Value(HeaderTraceID).(string)).
				Str("ip", c.ClientIP()).
				Msg("Authorization header is empty or malformed")
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
				log.Error().
					Err(err).
					Str(HeaderTraceID, c.Request.Context().Value(HeaderTraceID).(string)).
					Str("ip", c.ClientIP()).
					Msg("Token signature is not valid")
				RespondUnauthorized(c, "", InvalidTokenSignature)
				return
			}
			log.Error().
				Err(err).
				Str(HeaderTraceID, c.Request.Context().Value(HeaderTraceID).(string)).
				Str("ip", c.ClientIP()).
				Msg("Failed to parse jwt token")
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
			log.Error().
				Err(err).
				Str(HeaderTraceID, c.Request.Context().Value(HeaderTraceID).(string)).
				Str("ip", c.ClientIP()).
				Msg("Invalid token")
			RespondUnauthorized(c, "", InvalidToken)
			return
		}

		c.Set(KeySubjectJWT, claims.Subject)

		c.Next()
	}
}
