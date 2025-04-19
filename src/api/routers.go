package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "mistapi/docs"
	"mistapi/src/auth"
	"mistapi/src/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func StartService() {

	// initialize grpc connection
	service.GetGrpcClientConnection()
	defer service.CloseGrpcConnection()

	r := SetupRouter()

	addr := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	// TODO: use better logging solution
	log.Printf("Server running at %s\n", addr)
	http.ListenAndServe(addr, r)
}

func SetupRouter() *chi.Mux {
	r := chi.NewRouter()

	// SETUP MIDDDLEWARES
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	// Mount the user router
	r.Get("/health", HealthHandler)

	r.Route("/api/", func(r chi.Router) {
		r.Use(auth.AuthenticateMiddleware)

		r.Mount("/v1/appserver", appserverRouter())
		r.Mount("/v1/appserver-role", appserverRoleRouter())
		r.Mount("/v1/appserver-sub", appserverSubRouter())
	})

	// TODO: change the localhost domain
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", os.Getenv("APP_PORT")))))

	return r
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
