package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

func (cfg *ApiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
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

func (cfg *ApiConfig) handleGetUsers(w http.ResponseWriter, r *http.Request, user database.User) {

	respondWithJSON(w, 200, dbUserToUser(user))
}

func (cfg *ApiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	type response struct {
		Feed       Feed       `json:"feed"`
		FeedFollow FeedFollow `json:"feed_follow"`
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

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feed.ID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		respondWithError(w, 400, "Error creating feed follow")
		return
	}

	respondWithJSON(w, 201, response{
		Feed:       dbFeedToFeed(feed),
		FeedFollow: dbFollowToFollow(feedFollow),
	})
}

func (cfg *ApiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {

	dbFeeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 500, "Error getting feeds")
		return
	}

	respondWithJSON(w, 200, dbFeedsToFeeds(dbFeeds))
}

func (cfg *ApiConfig) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "")
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    params.FeedId,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		respondWithError(w, 400, "Error following feed")
		return
	}

	respondWithJSON(w, 201, dbFollowToFollow(feedFollow))
}

func (cfg *ApiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	follow_idStr := r.PathValue("feedFollowID")
	if err := uuid.Validate(follow_idStr); err != nil {
		respondWithError(w, 403, "Invalid follow ID")
	}

	follow_id, err := uuid.Parse(follow_idStr)
	if err != nil {
		respondWithError(w, 500, "Error parsing id")
		return
	}

	if err = cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     follow_id,
		UserID: user.ID,
	}); err != nil {
		respondWithError(w, 400, err.Error())
	}

	w.WriteHeader(http.StatusOK)

}

func (cfg *ApiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	dbFeedFollows, err := cfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 500, "Error getting feed follows")
		return
	}

	respondWithJSON(w, 200, dbFollowsToFollows(dbFeedFollows))
}

func (cfg *ApiConfig) handleGetPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Printf("Error converting limit number: %s", err)
	}
	if limit == 0 {
		limit = 10
	}

	dbPosts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		Limit:  int32(limit),
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 500, "Error getting feed follows")
		return
	}

	respondWithJSON(w, 200, dbPostsToPosts(dbPosts))
}
