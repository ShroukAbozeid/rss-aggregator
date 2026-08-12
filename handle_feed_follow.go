package main

import (
	"encoding/json"
	"net/http"

	"github.com/ShroukAbozeid/rss-aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	follow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:     uuid.New(),
		FeedID: params.FeedId,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to follow feed")
		return
	}
	respondWithJSON(w, http.StatusOK, databaseFeedFollowToAPIFeedFollow(follow))
}

func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	follows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get feed follows")
		return
	}
	respondWithJSON(w, http.StatusOK, databaseFeedFollowsToAPIFeedFollows(follows))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowID, err := uuid.Parse(chi.URLParam(r, "feedFollowID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feed_id")
		return
	}
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to unfollow feed")
		return
	}
	respondWithJSON(w, http.StatusOK, struct{}{})
}
