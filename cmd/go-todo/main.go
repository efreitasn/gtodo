package main

import (
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/efreitasn/go-todo/cmd/go-todo/handlers"
)

func main() {
	mux := httptreemux.NewContextMux()

	mux.GET("/todos", handlers.ListTodos)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
