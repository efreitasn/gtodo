package handlers

import (
	"net/http"

	"github.com/dimfeld/httptreemux"
	"go.mongodb.org/mongo-driver/mongo"
)

// Setup generates a mux handler.
func Setup(db *mongo.Database) http.Handler {
	mux := httptreemux.NewContextMux()
	todo := Todo{db}

	mux.GET("/todos", todo.List)

	return mux
}
