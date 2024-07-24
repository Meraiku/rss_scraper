package server

import (
	"encoding/json"
	"net/http"
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

	cfg.DB.CreateUser(r.Context())
}
