package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "mistapi/docs"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcConnectionKey string

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func StartService() {

	clientConn, err := grpc.NewClient(
		os.Getenv("MIST_BACKEND_APP_URL"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	defer clientConn.Close()

	if err != nil {
		log.Panicf("Error communicating with backend service: %v", err)
	}

	r := chi.NewRouter()

	// SETUP MIDDDLEWARES
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(setGRPCConnection(clientConn))

	// Mount the user router
	r.Mount("/api/v1/appserver", appserverRouter())

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", os.Getenv("APP_PORT")))))

	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), r)

}

func setGRPCConnection(clientConn *grpc.ClientConn) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, grpcConnectionKey("grpc_conn"), clientConn)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}

}

func GetGRPCConnFromContext(r *http.Request) *grpc.ClientConn {
	return r.Context().Value(grpcConnectionKey("grpc_conn")).(*grpc.ClientConn)
}
