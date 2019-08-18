package handlers

import (
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewMux creates a mux handler.
func NewMux(db *mongo.Database) http.Handler {
	mux := httptreemux.NewContextMux()
	todo := Todo{db}

	mux.GET("/todos", todo.List)

	return mux
}
