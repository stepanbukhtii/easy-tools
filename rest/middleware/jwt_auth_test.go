package middleware

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stepanbukhtii/easy-tools/rest/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stepanbukhtii/easy-tools/econtext"
	"github.com/stepanbukhtii/easy-tools/ejwt"
)

const (
	testIssuer   = "issuer"
	testAudience = "audience"
	testSubject  = "subject"
	testRole     = "role"
)

func TestMiddleware_Auth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	assert.NoError(t, err)

	_, wrongPrivateKey, err := ed25519.GenerateKey(nil)
	assert.NoError(t, err)

	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	now := time.Now().UTC()

	tests := []struct {
		name             string
		signingMethod    jwt.SigningMethod
		registeredClaims *jwt.RegisteredClaims
		signPrivateKey   any
		skipValidation   bool
		expectedStatus   int
	}{
		{
			name:          "success",
			signingMethod: jwt.SigningMethodEdDSA,
			registeredClaims: &jwt.RegisteredClaims{
				Issuer:    testIssuer,
				Subject:   testSubject,
				Audience:  jwt.ClaimStrings{testAudience},
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
				NotBefore: jwt.NewNumericDate(now),
				IssuedAt:  jwt.NewNumericDate(now),
			},
			signPrivateKey: privateKey,
			expectedStatus: http.StatusOK,
		}, {
			name:          "invalid issuer",
			signingMethod: jwt.SigningMethodEdDSA,
			registeredClaims: &jwt.RegisteredClaims{
				Issuer:    "wrong_issuer",
				Subject:   testSubject,
				Audience:  jwt.ClaimStrings{testAudience},
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
				NotBefore: jwt.NewNumericDate(now),
				IssuedAt:  jwt.NewNumericDate(now),
			},
			signPrivateKey: privateKey,
			expectedStatus: http.StatusUnauthorized,
		}, {
			name:          "invalid audience",
			signingMethod: jwt.SigningMethodEdDSA,
			registeredClaims: &jwt.RegisteredClaims{
				Issuer:    testIssuer,
				Subject:   testSubject,
				Audience:  jwt.ClaimStrings{"wrong_audience"},
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
				NotBefore: jwt.NewNumericDate(now),
				IssuedAt:  jwt.NewNumericDate(now),
			},
			signPrivateKey: privateKey,
			expectedStatus: http.StatusUnauthorized,
		}, {
			name:          "invalid expires at",
			signingMethod: jwt.SigningMethodEdDSA,
			registeredClaims: &jwt.RegisteredClaims{
				Issuer:    testIssuer,
				Subject:   testSubject,
				Audience:  jwt.ClaimStrings{testAudience},
				ExpiresAt: jwt.NewNumericDate(now.Add(-time.Minute)),
				NotBefore: jwt.NewNumericDate(now),
				IssuedAt:  jwt.NewNumericDate(now),
			},
			signPrivateKey: privateKey,
			expectedStatus: http.StatusUnauthorized,
		}, {
			name:          "invalid not before at",
			signingMethod: jwt.SigningMethodEdDSA,
			registeredClaims: &jwt.RegisteredClaims{
				Issuer:    testIssuer,
				Subject:   testSubject,
				Audience:  jwt.ClaimStrings{testAudience},
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
				NotBefore: jwt.NewNumericDate(now.Add(30 * time.Second)),
				IssuedAt:  jwt.NewNumericDate(now),
			},
			signPrivateKey: privateKey,
			expectedStatus: http.StatusUnauthorized,
		}, {
			name:          "invalid issued at",
			signingMethod: jwt.SigningMethodEdDSA,
			registeredClaims: &jwt.RegisteredClaims{
				Issuer:    testIssuer,
				Subject:   testSubject,
				Audience:  jwt.ClaimStrings{testAudience},
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
				NotBefore: jwt.NewNumericDate(now),
				IssuedAt:  jwt.NewNumericDate(now.Add(time.Second)),
			},
			signPrivateKey: privateKey,
			expectedStatus: http.StatusUnauthorized,
		}, {
			name:          "invalid signing method",
			signingMethod: jwt.SigningMethodRS256,
			registeredClaims: &jwt.RegisteredClaims{
				Issuer:    testIssuer,
				Subject:   testSubject,
				Audience:  jwt.ClaimStrings{testAudience},
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
				NotBefore: jwt.NewNumericDate(now),
				IssuedAt:  jwt.NewNumericDate(now),
			},
			signPrivateKey: rsaPrivateKey,
			expectedStatus: http.StatusUnauthorized,
		}, {
			name:          "invalid private key",
			signingMethod: jwt.SigningMethodEdDSA,
			registeredClaims: &jwt.RegisteredClaims{
				Issuer:    testIssuer,
				Subject:   testSubject,
				Audience:  jwt.ClaimStrings{testAudience},
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
				NotBefore: jwt.NewNumericDate(now),
				IssuedAt:  jwt.NewNumericDate(now),
			},
			signPrivateKey: wrongPrivateKey,
			expectedStatus: http.StatusUnauthorized,
		}, {
			name:          "skip validation claims",
			signingMethod: jwt.SigningMethodEdDSA,
			registeredClaims: &jwt.RegisteredClaims{
				Issuer:    "wrong_issuer",
				Subject:   "wrong_subject",
				Audience:  jwt.ClaimStrings{"wrong_audience"},
				ExpiresAt: jwt.NewNumericDate(now.Add(-time.Minute)),
				NotBefore: jwt.NewNumericDate(now.Add(time.Minute)),
				IssuedAt:  jwt.NewNumericDate(now.Add(2 * time.Minute)),
			},
			signPrivateKey: privateKey,
			skipValidation: true,
			expectedStatus: http.StatusOK,
		}, {
			name:          "skip validation method and signature",
			signingMethod: jwt.SigningMethodRS256,
			registeredClaims: &jwt.RegisteredClaims{
				Issuer:    testIssuer,
				Subject:   testSubject,
				Audience:  jwt.ClaimStrings{testAudience},
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
				NotBefore: jwt.NewNumericDate(now),
				IssuedAt:  jwt.NewNumericDate(now),
			},
			signPrivateKey: rsaPrivateKey,
			skipValidation: true,
			expectedStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token, err := jwt.NewWithClaims(test.signingMethod, test.registeredClaims).SignedString(test.signPrivateKey)
			assert.NoError(t, err)

			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
			c.Request.Header.Set(api.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

			middleware := NewJWTAuth(publicKey, testIssuer, testAudience, test.skipValidation)
			middleware.Auth(c)

			require.Equal(t, test.expectedStatus, c.Writer.Status())
			if c.Writer.Status() == http.StatusOK {
				require.Equal(t, test.registeredClaims.Subject, econtext.ClientInfo(c.Request.Context()).Subject)
			}
		})
	}
}

func TestMiddleware_AuthRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	assert.NoError(t, err)

	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	now := time.Now().UTC()

	tests := []struct {
		name             string
		signingMethod    jwt.SigningMethod
		registeredClaims *ejwt.ClaimsWithRoles
		signPrivateKey   any
		skipValidation   bool
		expectedStatus   int
	}{
		{
			name:          "success",
			signingMethod: jwt.SigningMethodEdDSA,
			registeredClaims: &ejwt.ClaimsWithRoles{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    testIssuer,
					Subject:   testSubject,
					Audience:  jwt.ClaimStrings{testAudience},
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
					NotBefore: jwt.NewNumericDate(now),
					IssuedAt:  jwt.NewNumericDate(now),
				},
				Roles: []string{testRole},
			},
			signPrivateKey: privateKey,
			expectedStatus: http.StatusOK,
		}, {
			name:          "invalid role",
			signingMethod: jwt.SigningMethodEdDSA,
			registeredClaims: &ejwt.ClaimsWithRoles{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    testIssuer,
					Subject:   testSubject,
					Audience:  jwt.ClaimStrings{testAudience},
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
					NotBefore: jwt.NewNumericDate(now),
					IssuedAt:  jwt.NewNumericDate(now),
				},
				Roles: []string{"invalid_role"},
			},
			signPrivateKey: privateKey,
			expectedStatus: http.StatusUnauthorized,
		}, {
			name:          "skip validation role",
			signingMethod: jwt.SigningMethodEdDSA,
			registeredClaims: &ejwt.ClaimsWithRoles{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    testIssuer,
					Subject:   testSubject,
					Audience:  jwt.ClaimStrings{testAudience},
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
					NotBefore: jwt.NewNumericDate(now),
					IssuedAt:  jwt.NewNumericDate(now),
				},
				Roles: []string{"invalid_role"},
			},
			signPrivateKey: privateKey,
			skipValidation: true,
			expectedStatus: http.StatusOK,
		}, {
			name:          "skip validation method and signature",
			signingMethod: jwt.SigningMethodRS256,
			registeredClaims: &ejwt.ClaimsWithRoles{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    testIssuer,
					Subject:   testSubject,
					Audience:  jwt.ClaimStrings{testAudience},
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
					NotBefore: jwt.NewNumericDate(now),
					IssuedAt:  jwt.NewNumericDate(now),
				},
				Roles: []string{testRole},
			},
			signPrivateKey: rsaPrivateKey,
			skipValidation: true,
			expectedStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token, err := jwt.NewWithClaims(test.signingMethod, test.registeredClaims).SignedString(test.signPrivateKey)
			assert.NoError(t, err)

			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
			c.Request.Header.Set(api.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

			fmt.Println("token", token)

			middleware := NewJWTAuth(publicKey, testIssuer, testAudience, test.skipValidation)
			middleware.AuthRole(testRole)(c)

			require.Equal(t, test.expectedStatus, c.Writer.Status())
			if c.Writer.Status() == http.StatusOK {
				require.Equal(t, test.registeredClaims.Subject, econtext.ClientInfo(c.Request.Context()).Subject)
				require.Equal(t, test.registeredClaims.Roles, econtext.ClientInfo(c.Request.Context()).Roles)
			}
		})
	}
}
