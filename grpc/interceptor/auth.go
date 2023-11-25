package interceptor

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stepanbukhtii/easy-tools/econtext"
	"github.com/stepanbukhtii/easy-tools/ejwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	AuthScheme          = "bearer"
	HeaderAuthorization = "authorization"
)

type JWTAuth struct {
	publicKey            ed25519.PublicKey
	parser               *jwt.Parser
	skipValidation       bool
	expirationValidation bool
}

func NewJWTAuth(publicKey ed25519.PublicKey, issuer, audience string, skipValidation, expirationValidate bool) *JWTAuth {
	options := []jwt.ParserOption{
		jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}),
		jwt.WithIssuer(issuer),
		jwt.WithAudience(audience),
		jwt.WithIssuedAt(),
	}

	if expirationValidate {
		options = append(options, jwt.WithExpirationRequired())
	}

	return &JWTAuth{
		publicKey:      publicKey,
		parser:         jwt.NewParser(options...),
		skipValidation: skipValidation,
	}
}

// Auth parse and verifying JWT token
func (m JWTAuth) Auth(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	tokenString, err := extractAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	if tokenString == "" {
		return ctx, nil
	}

	var claims jwt.RegisteredClaims

	if m.skipValidation {
		if _, _, err := m.parser.ParseUnverified(tokenString, &claims); err != nil {
			return nil, status.Error(codes.Unauthenticated, "parse token failed")
		}

		if claims.Subject != "" {
			ctx = econtext.SetSubject(ctx, claims.Subject)
		}

		return handler(ctx, req)
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) { return m.publicKey, nil }
	token, err := m.parser.ParseWithClaims(tokenString, &claims, keyFunc)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "parse token failed")
	}

	if !token.Valid || claims.Subject == "" {
		return nil, status.Error(codes.Unauthenticated, "token validating failed")
	}

	ctx = econtext.SetSubject(ctx, claims.Subject)

	return handler(ctx, req)
}

// AuthRole parse and verifying JWT token and check role
func (m JWTAuth) AuthRole(role string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		tokenString, err := extractAuthToken(ctx)
		if err != nil {
			return nil, err
		}
		if tokenString == "" {
			return ctx, nil
		}

		var claims ejwt.ClaimsWithRoles

		if m.skipValidation {
			if _, _, err := m.parser.ParseUnverified(tokenString, &claims); err != nil {
				return nil, status.Error(codes.Unauthenticated, "parse token failed")
			}

			if claims.Subject != "" {
				ctx = econtext.SetSubject(ctx, claims.Subject)
			}

			if claims.Roles != nil {
				ctx = econtext.SetRoles(ctx, claims.Roles)
			}

			return handler(ctx, req)
		}

		keyFunc := func(token *jwt.Token) (interface{}, error) { return m.publicKey, nil }
		token, err := m.parser.ParseWithClaims(tokenString, &claims, keyFunc)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "parse token failed")
		}

		var hasRole bool
		for _, claimRole := range claims.Roles {
			if claimRole == role {
				hasRole = true
				break
			}
		}

		if !token.Valid || claims.Subject == "" || !hasRole {
			return nil, status.Error(codes.Unauthenticated, "token validating failed")
		}

		ctx = econtext.SetSubject(ctx, claims.Subject)
		ctx = econtext.SetRoles(ctx, claims.Roles)

		return handler(ctx, req)
	}
}

func ClientAuth(
	ctx context.Context,
	method string,
	req, reply any,
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	token := econtext.AuthToken(ctx)
	if token == "" {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	token = fmt.Sprintf("%s %s", AuthScheme, token)

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.Pairs(HeaderAuthorization, token)
	} else {
		md = md.Copy()
		md.Set(HeaderAuthorization, token)
	}

	ctx = metadata.NewOutgoingContext(ctx, md)

	return invoker(ctx, method, req, reply, cc, opts...)
}

func extractAuthToken(ctx context.Context) (string, error) {
	vals := metadata.ValueFromIncomingContext(ctx, HeaderAuthorization)
	if len(vals) == 0 {
		return "", nil
	}

	scheme, tokenString, found := strings.Cut(vals[0], " ")
	if !found {
		return "", status.Error(codes.Unauthenticated, "bad authorization string")
	}

	if !strings.EqualFold(scheme, AuthScheme) {
		return "", status.Errorf(codes.Unauthenticated, "request unauthenticated with %s", AuthScheme)
	}

	return tokenString, nil
}
