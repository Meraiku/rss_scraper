package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func StartServer() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	router := chi.NewRouter()
	router.Use(cors.AllowAll().Handler)

	routerV1 := chi.NewRouter()

	routerV1.Get("/healthz", handleHealthz)
	routerV1.Get("/err", handleError)

	router.Mount("/v1", routerV1)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	fmt.Printf("Server starting at: %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
