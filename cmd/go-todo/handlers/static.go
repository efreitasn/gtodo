package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

func static(w http.ResponseWriter, r *http.Request) {
	path, _ := os.Getwd()
	path = filepath.Join(path, "web/static")

	http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir(path)),
	).ServeHTTP(w, r)
}
