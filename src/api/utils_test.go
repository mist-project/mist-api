package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"mistapi/src/api"
	"mistapi/src/auth"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateTokenParams struct {
	iss       string
	aud       []string
	secretKey string
	userId    string
}

func addContextHeaders(req *http.Request) *http.Request {
	claims := &auth.CustomJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   os.Getenv("MIST_PY_API_JWT_ISSUER"),
			Audience: []string{os.Getenv("MIST_PY_API_JWT_AUDIENCE")},

			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: "123",
	}

	token := &auth.TokenAndClaims{
		Claims: claims,
		Token:  "mockToken123",
	}

	ctx := context.WithValue(req.Context(), auth.TokenContextKey, token)
	return req.WithContext(ctx)
}

func withURLParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	rctx.URLParams.Add(key, value)
	return r
}

func marshallPayload(t *testing.T, data interface{}) *bytes.Buffer {
	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("error marshalling body: %v", err)
	}
	return bytes.NewBuffer(body)
}

func marshallResponse(t *testing.T, data interface{}) string {
	expected, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("error marshalling response: %v", err)
	}

	return string(expected)
}

func TestHandleGrpcError(t *testing.T) {
	t.Parallel()
	log.SetOutput(new(strings.Builder))

	tests := []struct {
		name           string
		grpcCode       codes.Code
		expectedStatus int
		expectedDetail string
	}{
		{"Unavailable", codes.Unavailable, http.StatusBadGateway, "Server is unresponsive."},
		{"DeadlineExceeded", codes.DeadlineExceeded, http.StatusBadGateway, "Server timed out."},
		{"Canceled", codes.Canceled, http.StatusBadGateway, "Server error."},
		{"Unauthenticated", codes.Unauthenticated, http.StatusUnauthorized, "Unauthorized request."},
		{"NotFound", codes.NotFound, http.StatusNotFound, "Not found."},
		{"AlreadyExists", codes.AlreadyExists, http.StatusConflict, "Resource already exists."},
		{"InvalidArgument", codes.InvalidArgument, http.StatusBadRequest, "simulated error"},
		{"UnhandledCode", codes.DataLoss, http.StatusInternalServerError, "Internal Server Error."},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			// Silence logs for cleaner test output

			// ARRANGE
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/some-endpoint", nil)
			err := status.Error(tt.grpcCode, "simulated error")

			// ACT
			api.HandleGrpcError(w, r, err)

			// ASSERT
			expected, _ := json.Marshal(api.CreateErrorResponse(tt.expectedDetail))

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t,
				string(expected),
				w.Body.String(),
			)
		})
	}
}
