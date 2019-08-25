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

	// Root
	mux.GET("/", root)

	// Static
	mux.GET("/static/*", static)

	// List
	mux.GET("/list", SetUpTemplateData(todo.ListGET))

	// Add
	mux.GET("/add", SetUpTemplateData(todo.AddGET))
	mux.POST("/add", todo.AddPOST)

	// Update
	mux.GET("/update", SetUpTemplateData(todo.UpdateGET))
	mux.POST("/update", todo.UpdatePOST)

	// Delete
	mux.GET("/delete", SetUpTemplateData(todo.DeleteGET))
	mux.POST("/delete", todo.DeletePOST)

	return *mux
}
