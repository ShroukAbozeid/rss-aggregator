package main

import (
	"fmt"
	"net/http"

	"github.com/ShroukAbozeid/rss-aggregator/internal/auth"
	"github.com/ShroukAbozeid/rss-aggregator/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

// this is to be able to use it with chi router handler signuatre
func (apiCfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Auth error: %v", err))
			return
		}
		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Auth error: %v", err))
			return
		}
		handler(w, r, user)
	}
}
