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
	"github.com/rs/cors"
)

func StartService() {

	// initialize grpc connection
	service.GetGrpcClientConnection()
	defer service.CloseGrpcConnection()

	r := SetupRouter()

	// Apply CORS
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // TODO: fix the origin for the app
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: false, // if sending cookies/auth headers
	}).Handler(r)

	certFile := "/etc/ssl/cloudflare/mist-project.crt"
	keyFile := "/etc/ssl/cloudflare/mist-project.key"

	addr := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	// TODO: use better logging solution
	log.Printf("Server running at %s\n", addr)
	http.ListenAndServeTLS(addr, certFile, keyFile, handler)
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

		r.Mount("/v1/appservers", appserverRouter())
		r.Mount("/v1/appserver-roles", appserverRoleRouter())
		r.Mount("/v1/appserver-role-subs", appserverRoleSubRouter())
		r.Mount("/v1/appserver-subs", appserverSubRouter())
		r.Mount("/v1/channels", channelRouter())
		r.Mount("/v1/channel-roles", channelRoleRouter())
	})

	// TODO: change the localhost domain
	// r.Get("/swagger/*", httpSwagger.Handler(
	// 	httpSwagger.URL(fmt.Sprintf("https://192.168.0.21:%s/swagger/doc.json", os.Getenv("APP_PORT")))))

	return r
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
