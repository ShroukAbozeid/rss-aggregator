package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		jsonErrorResponse, _ := json.Marshal(errorResponse{Error: "Internal Server Error"})
		w.Write(jsonErrorResponse)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Printf("Server error: %s", message)
	}

	respondWithJSON(w, code, errorResponse{Error: message})
}
