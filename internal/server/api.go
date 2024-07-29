package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/meraiku/rss_scraper/internal/database"
	"github.com/meraiku/rss_scraper/internal/scraper"
)

type ApiConfig struct {
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

	cfg := ApiConfig{
		DB: database.New(db),
	}

	go scraper.StartScraper(cfg.DB, 10, time.Minute)

	router := chi.NewRouter()
	router.Use(cors.AllowAll().Handler)

	routerV1 := chi.NewRouter()

	routerV1.Get("/healthz", handleHealthz)
	routerV1.Get("/err", handleError)

	routerV1.Get("/users", cfg.middlewareAuth(cfg.handleGetUsers))
	routerV1.Post("/users", cfg.handleCreateUser)

	routerV1.Get("/feeds", cfg.handleGetFeeds)
	routerV1.Post("/feeds", cfg.middlewareAuth(cfg.handleCreateFeed))

	routerV1.Get("/feed_follows", cfg.middlewareAuth(cfg.handleGetFeedFollows))
	routerV1.Post("/feed_follows", cfg.middlewareAuth(cfg.handleCreateFeedFollow))
	routerV1.Delete("/feed_follows/{feedFollowID}", cfg.middlewareAuth(cfg.handleDeleteFeedFollow))

	routerV1.Get("/posts", cfg.middlewareAuth(cfg.handleGetPostsByUser))

	router.Mount("/v1", routerV1)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	fmt.Printf("Server starting at: %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
