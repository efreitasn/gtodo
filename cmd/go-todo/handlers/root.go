package handlers

import "net/http"

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "/list")
	w.WriteHeader(301)
}
