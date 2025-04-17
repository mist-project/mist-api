package auth_test

import (
	"mistapi/src/auth"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

type DummyRequest struct{}

func TestAuthenticateMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name: "valid_token_is_successful",
			authHeader: bearerToken(createJwtToken(t, &CreateTokenParams{
				iss:       os.Getenv("MIST_PY_API_JWT_ISSUER"),
				aud:       []string{os.Getenv("MIST_PY_API_JWT_AUDIENCE")},
				secretKey: os.Getenv("MIST_PY_API_JWT_SECRET_KEY"),
			})),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid_token_returns_unauthorized",
			authHeader:     "Bearer invalid.jwt.token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", tt.authHeader)

			rr := httptest.NewRecorder()

			// dummy next handler to verify if middleware passes control
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			})

			// ACT
			handler := auth.AuthenticateMiddleware(next)
			handler.ServeHTTP(rr, req)

			// ASSERT
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestAuthorizeToken(t *testing.T) {
	t.Run("valid_token_is_successful", func(t *testing.T) {
		// ARRANGE
		token := createJwtToken(t, &CreateTokenParams{
			iss:       os.Getenv("MIST_PY_API_JWT_ISSUER"),
			aud:       []string{os.Getenv("MIST_PY_API_JWT_AUDIENCE")},
			secretKey: os.Getenv("MIST_PY_API_JWT_SECRET_KEY"),
		})

		// ACT
		tac, err := auth.AuthorizeToken(bearerToken(token))

		// ASSERT
		assert.Nil(t, err)
		assert.Equal(t, tac.Token, token)
	})

	t.Run("token_with_invalid_audience_errors", func(t *testing.T) {
		// ARRANGE
		token := createJwtToken(t, &CreateTokenParams{
			iss:       os.Getenv("MIST_PY_API_JWT_ISSUER"),
			aud:       []string{"invalid-audience"},
			secretKey: os.Getenv("MIST_PY_API_JWT_SECRET_KEY"),
		})

		// ACT
		tac, err := auth.AuthorizeToken(bearerToken(token))

		// ASSERT
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid audience claim")
		assert.Nil(t, tac)
	})

	t.Run("token_with_invalid_issuer_errors", func(t *testing.T) {
		// ARRANGE
		token := createJwtToken(t, &CreateTokenParams{
			aud:       []string{os.Getenv("MIST_PY_API_JWT_AUDIENCE")},
			secretKey: os.Getenv("MIST_PY_API_JWT_SECRET_KEY"),
		})

		// ACT
		tac, err := auth.AuthorizeToken(bearerToken(token))

		// ASSERT
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid issuer claim")
		assert.Nil(t, tac)
	})

	t.Run("token_with_invalid_secret_key_errors", func(t *testing.T) {
		// ARRANGE
		token := createJwtToken(t, &CreateTokenParams{
			iss:       os.Getenv("MIST_PY_API_JWT_ISSUER"),
			aud:       []string{os.Getenv("MIST_PY_API_JWT_AUDIENCE")},
			secretKey: "wrong-secret-key",
		})

		// ACT
		tac, err := auth.AuthorizeToken(bearerToken(token))

		// ASSERT
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "signature is invalid")
		assert.Nil(t, tac)
	})

	t.Run("token_with_invalid_format_errors", func(t *testing.T) {
		// ARRANGE
		badToken := "bad_token"

		// ACT
		tac, err := auth.AuthorizeToken(badToken)

		// ASSERT
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid token format")
		assert.Nil(t, tac)
	})

	t.Run("token_with_missing_authorization_header_errors", func(t *testing.T) {
		// ARRANGE
		missingToken := ""

		// ACT
		tac, err := auth.AuthorizeToken(missingToken)

		// ASSERT
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid token")
		assert.Nil(t, tac)
	})

	t.Run("token_with_invalid_authorization_bearer_header_errors", func(t *testing.T) {
		// ARRANGE
		token := "token_invalid"

		// ACT
		tac, err := auth.AuthorizeToken(bearerToken(token))

		// ASSERT
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "token is malformed: token contains an invalid number of segments")
		assert.Nil(t, tac)
	})

	t.Run("empty_authorization_bearer_header_errors", func(t *testing.T) {
		// ARRANGE
		token := ""

		// ACT
		tac, err := auth.AuthorizeToken(bearerToken(token))

		// ASSERT
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "token is malformed: token contains an invalid number of segments")
		assert.Nil(t, tac)
	})

	t.Run("token_with_invalid_claims_format_for_audience_errors", func(t *testing.T) {
		// ARRANGE
		claims := &jwt.RegisteredClaims{
			Issuer:    os.Getenv("MIST_PY_API_JWT_ISSUER"),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte(os.Getenv("MIST_PY_API_JWT_SECRET_KEY")))
		assert.Nil(t, err)

		// ACT
		tac, err := auth.AuthorizeToken(bearerToken(tokenStr))

		// ASSERT
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid audience claim")
		assert.Nil(t, tac)
	})
}
