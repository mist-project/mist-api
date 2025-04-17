package api

import (
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
		// TODO: change the localhost domain
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", os.Getenv("APP_PORT")))))

	addr := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	fmt.Printf("Server running at %s\n", addr)
	http.ListenAndServe(addr, r)

}
