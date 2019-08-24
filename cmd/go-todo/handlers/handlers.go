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

	mux.GET("/", todo.List)
	mux.GET("/todos/delete", todo.DeleteList)
	mux.POST("/todos/add", todo.Add)
	mux.POST("/todos/update", todo.Update)
	mux.GET("/static/*", static)

	return mux
}
