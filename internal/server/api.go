package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/meraiku/rss_scraper/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func StartServer() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbURL := os.Getenv("DB_URL")
	if port == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting database: %s", err)
	}

	cfg := apiConfig{
		DB: database.New(db),
	}

	router := chi.NewRouter()
	router.Use(cors.AllowAll().Handler)

	routerV1 := chi.NewRouter()

	routerV1.Get("/healthz", handleHealthz)
	routerV1.Get("/err", handleError)

	routerV1.Get("/users", cfg.handleGetUserByAPIKey)
	routerV1.Post("/users", cfg.handleCreateUser)

	router.Mount("/v1", routerV1)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	fmt.Printf("Server starting at: %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
