package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func StartServer() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", handleHealthz)
	mux.HandleFunc("GET /v1/err", handleError)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Server starting at: %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
