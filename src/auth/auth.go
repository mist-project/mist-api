package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type CustomJWTClaims struct {
	jwt.RegisteredClaims        // Embed the standard registered claims
	UserID               string `json:"user_id"`
}

type TokenAndClaims struct {
	Claims *CustomJWTClaims
	Token  string
}

type contextKey string

const tokenContextKey = contextKey("auth_token")

func AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		tac, err := AuthorizeToken(authorization)

		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// TODO: add authorization in the future

		// Add to context
		ctx := context.WithValue(r.Context(), tokenContextKey, tac)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthorizeToken(authorization string) (*TokenAndClaims, error) {
	parts := strings.Split(authorization, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("invalid token format")
	}

	claims, err := verifyJWT(parts[1])
	if err != nil {
		return nil, err
	}

	return &TokenAndClaims{
		Token:  parts[1],
		Claims: claims,
	}, nil
}

func GetAuthotizationToken(r *http.Request) *TokenAndClaims {
	return r.Context().Value(contextKey(tokenContextKey)).(*TokenAndClaims)
}

func verifyJWT(tokenStr string) (*CustomJWTClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenStr, &CustomJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// TODO: we will need this in the future, for now skip
		// if token.Method != jwt.SigningMethodHS256 {
		// 	return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		// }
		// Return the secret key to validate the token's signature
		return []byte(os.Getenv("MIST_PY_API_JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	// Now validate the token's claims
	claims, err := verifyJWTTokenClaims(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func verifyJWTTokenClaims(token *jwt.Token) (*CustomJWTClaims, error) {
	// Now validate the token's claims
	claims, _ := token.Claims.(*CustomJWTClaims)

	// Validate aud
	validAudience := false
	auds := claims.Audience

	// If "aud" is an array of strings, cast each element to string
	for _, aud := range auds {
		if aud == os.Getenv("MIST_PY_API_JWT_AUDIENCE") {
			validAudience = true
			break
		}
	}

	if !validAudience {
		return nil, fmt.Errorf("invalid audience claim")
	}

	// Validate the issuer (iss) claim
	if claims.Issuer != os.Getenv("MIST_PY_API_JWT_ISSUER") {
		return nil, fmt.Errorf("invalid issuer claim")
	}

	// AuthJWTClaims
	return claims, nil
}
