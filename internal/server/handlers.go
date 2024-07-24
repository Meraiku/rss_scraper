package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
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

func (cfg *apiConfig) handleGetUsers(w http.ResponseWriter, r *http.Request, user database.User) {

	respondWithJSON(w, 200, dbUserToUser(user))
}

func (cfg *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "")
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, 400, "Error creating feed")
		return
	}

	respondWithJSON(w, 200, dbFeedToFeed(feed))
}

func (cfg *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {

	dbFeeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 500, "Error getting feeds")
		return
	}
	feeds := []Feed{}

	for _, dbFeed := range dbFeeds {
		feeds = append(feeds, dbFeedToFeed(dbFeed))
	}

	respondWithJSON(w, http.StatusCreated, feeds)
}
