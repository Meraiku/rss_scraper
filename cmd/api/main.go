package main

import (
	"github.com/joho/godotenv"
	"github.com/meraiku/rss_scraper/internal/server"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load("./.env")

	server.StartServer()
}
