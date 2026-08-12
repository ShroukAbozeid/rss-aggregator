package main

import (
	"encoding/json"
	"net/http"

	"github.com/ShroukAbozeid/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   params.Name,
		Url:    params.Url,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create feed")
		return
	}
	respondWithJSON(w, http.StatusOK, databaseFeedToAPIFeed(feed))
}

func (apiCfg *apiConfig) handlerGetFeed(w http.ResponseWriter, r *http.Request, feedUUID uuid.UUID) {

	feed, err := apiCfg.DB.GetFeed(r.Context(), feedUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get feed")
		return
	}
	respondWithJSON(w, http.StatusOK, databaseFeedToAPIFeed(feed))
}

func (apiCfg *apiConfig) handlerGetUserFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := apiCfg.DB.GetUserFeeds(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get feeds")
		return
	}
	respondWithJSON(w, http.StatusOK, databseFeedsToAPIFeeds(feeds))
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get feeds")
		return
	}

	respondWithJSON(w, http.StatusOK, databseFeedsToAPIFeeds(feeds))
}
