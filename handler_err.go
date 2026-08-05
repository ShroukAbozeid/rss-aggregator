package main

import "net/http"

func handlerError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusBadRequest, "some error occurred")
}
