package handlers

import (
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewMux creates a mux handler.
func NewMux(db *mongo.Database) http.Handler {
	mux := httptreemux.NewContextMux()
	todo := Todo{db.Collection("todos")}

	mux.GET("/todos", todo.List)
	mux.POST("/todos/add", todo.Add)

	return mux
}
