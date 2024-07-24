package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/meraiku/rss_scraper/internal/auth"
	"github.com/meraiku/rss_scraper/internal/database"
)

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Name string `json:"name"`
	}
	b := body{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&b)
	if err != nil {
		respondWithError(w, 400, "")
		return
	}

	dbUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.Must(uuid.NewRandom()),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      b.Name,
	})
	if err != nil {
		respondWithError(w, 500, "")
		return
	}
	respondWithJSON(w, http.StatusCreated, dbUserToUser(dbUser))
}

func (cfg *apiConfig) handleGetUserByAPIKey(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 403, fmt.Sprintf("Auth error: %s", err))
		return
	}

	user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, 400, "User with API key not found")
		return
	}

	respondWithJSON(w, 200, dbUserToUser(user))
}
