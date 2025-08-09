package auth_test

import (
	"context"
	"log"
	"mistapi/src/auth"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type DummyRequest struct{}

func TestAuthenticateMiddleware(t *testing.T) {
	log.SetOutput(new(strings.Builder))
	t.Parallel()

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name: "Success:valid_token_is_successful",
			authHeader: bearerToken(createJwtToken(t, &CreateTokenParams{
				iss:       os.Getenv("MIST_API_JWT_ISSUER"),
				aud:       []string{os.Getenv("MIST_API_JWT_AUDIENCE")},
				secretKey: os.Getenv("MIST_API_JWT_SECRET_KEY"),
			})),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error:Error:invalid_token_returns_unauthorized",
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
	log.SetOutput(new(strings.Builder))
	t.Parallel()

	t.Run("Success:valid_token_is_successful", func(t *testing.T) {
		// ARRANGE
		token := createJwtToken(t, &CreateTokenParams{
			iss:       os.Getenv("MIST_API_JWT_ISSUER"),
			aud:       []string{os.Getenv("MIST_API_JWT_AUDIENCE")},
			secretKey: os.Getenv("MIST_API_JWT_SECRET_KEY"),
		})

		// ACT
		tac, err := auth.AuthorizeToken(bearerToken(token))

		// ASSERT
		assert.Nil(t, err)
		assert.Equal(t, tac.Token, token)
	})

	t.Run("Error:token_with_invalid_audience_errors", func(t *testing.T) {
		// ARRANGE
		token := createJwtToken(t, &CreateTokenParams{
			iss:       os.Getenv("MIST_API_JWT_ISSUER"),
			aud:       []string{"invalid-audience"},
			secretKey: os.Getenv("MIST_API_JWT_SECRET_KEY"),
		})

		// ACT
		tac, err := auth.AuthorizeToken(bearerToken(token))

		// ASSERT
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid audience claim")
		assert.Nil(t, tac)
	})

	t.Run("Error:token_with_invalid_issuer_errors", func(t *testing.T) {
		// ARRANGE
		token := createJwtToken(t, &CreateTokenParams{
			aud:       []string{os.Getenv("MIST_API_JWT_AUDIENCE")},
			secretKey: os.Getenv("MIST_API_JWT_SECRET_KEY"),
		})

		// ACT
		tac, err := auth.AuthorizeToken(bearerToken(token))

		// ASSERT
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid issuer claim")
		assert.Nil(t, tac)
	})

	t.Run("Error:token_with_invalid_secret_key_errors", func(t *testing.T) {
		// ARRANGE
		token := createJwtToken(t, &CreateTokenParams{
			iss:       os.Getenv("MIST_API_JWT_ISSUER"),
			aud:       []string{os.Getenv("MIST_API_JWT_AUDIENCE")},
			secretKey: "wrong-secret-key",
		})

		// ACT
		tac, err := auth.AuthorizeToken(bearerToken(token))

		// ASSERT
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "signature is invalid")
		assert.Nil(t, tac)
	})

	t.Run("Error:token_with_invalid_format_errors", func(t *testing.T) {
		// ARRANGE
		badToken := "bad_token"

		// ACT
		tac, err := auth.AuthorizeToken(badToken)

		// ASSERT
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid token format")
		assert.Nil(t, tac)
	})

	t.Run("Error:token_with_missing_authorization_header_errors", func(t *testing.T) {
		// ARRANGE
		missingToken := ""

		// ACT
		tac, err := auth.AuthorizeToken(missingToken)

		// ASSERT
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid token")
		assert.Nil(t, tac)
	})

	t.Run("Error:token_with_invalid_authorization_bearer_header_errors", func(t *testing.T) {
		// ARRANGE
		token := "token_invalid"

		// ACT
		tac, err := auth.AuthorizeToken(bearerToken(token))

		// ASSERT
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "token is malformed: token contains an invalid number of segments")
		assert.Nil(t, tac)
	})

	t.Run("Error:empty_authorization_bearer_header_errors", func(t *testing.T) {
		// ARRANGE
		token := ""

		// ACT
		tac, err := auth.AuthorizeToken(bearerToken(token))

		// ASSERT
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "token is malformed: token contains an invalid number of segments")
		assert.Nil(t, tac)
	})

	t.Run("Error:token_with_invalid_claims_format_for_audience_errors", func(t *testing.T) {
		// ARRANGE
		claims := &jwt.RegisteredClaims{
			Issuer:    os.Getenv("MIST_API_JWT_ISSUER"),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte(os.Getenv("MIST_API_JWT_SECRET_KEY")))
		assert.Nil(t, err)

		// ACT
		tac, err := auth.AuthorizeToken(bearerToken(tokenStr))

		// ASSERT
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid audience claim")
		assert.Nil(t, tac)
	})
}

func TestGetAuthotizationToken(t *testing.T) {
	log.SetOutput(new(strings.Builder))
	t.Parallel()

	t.Run("when_token_is_added_in_context_it_is_returned", func(t *testing.T) {
		// ARRANGE
		// Prepare the data we want to store in the context
		expectedClaims := &auth.CustomJWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:  "exampleIssuer",
				Subject: "exampleSubject",
			},
			UserID: "testUserID",
		}
		expectedToken := &auth.TokenAndClaims{
			Claims: expectedClaims,
			Token:  "mockToken123",
		}

		// Create a new HTTP request using http.NewRequest
		req, _ := http.NewRequest("GET", "/api/v1/endpoint", nil)

		// Add the token and claims to the request's context
		ctx := context.WithValue(req.Context(), auth.TokenContextKey, expectedToken)
		req = req.WithContext(ctx)

		//ACT
		actualToken, _ := auth.GetAuthotizationToken(req)

		// ASSERT
		assert.NotNil(t, actualToken, "Expected a non-nil token")
		assert.Equal(t, expectedToken.Token, actualToken.Token, "Expected token values to match")
		assert.Equal(t, expectedToken.Claims, actualToken.Claims, "Expected claims to match")

	})

	t.Run("when_token_not_in_context_it_returns_nil", func(t *testing.T) {
		// ARRANGE
		// Create a new HTTP request using http.NewRequest
		req, _ := http.NewRequest("GET", "/api/v1/endpoint", nil)

		// ACT
		// Call GetAuthotizationToken with the mock request
		token, _ := auth.GetAuthotizationToken(req)

		// ASSERT
		// Check that the returned token is nil
		require.Nil(t, token, "Expected token to be nil when no token is set in context")
	})
}
