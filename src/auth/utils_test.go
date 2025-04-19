package auth_test

import (
	"fmt"
	"mistapi/src/auth"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CreateTokenParams struct {
	iss       string
	aud       []string
	secretKey string
	userId    string
}

func CreateTokenClaims(params *CreateTokenParams) *auth.CustomJWTClaims {
	return &auth.CustomJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   params.iss,
			Audience: params.aud,

			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: params.userId,
	}
}

func createJwtToken(t *testing.T, params *CreateTokenParams) string {
	// Define secret key for signing the token

	// Define JWT claims
	claims := CreateTokenClaims(params)

	// Create a new token with specified claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using the secret key
	tokenString, err := token.SignedString([]byte(params.secretKey))
	if err != nil {
		t.Fatalf("error signing the token %v", err)
	}
	return tokenString
}

func bearerToken(token string) string {
	return fmt.Sprintf("Bearer %s", token)
}
