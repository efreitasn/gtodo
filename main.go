package main

import (
	"net/http"

	"github.com/dimfeld/httptreemux"
)

func main() {
	mux := httptreemux.NewContextMux()

	mux.GET("/foobar", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("foobar"))
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
