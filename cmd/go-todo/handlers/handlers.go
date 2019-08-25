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
	auth := Auth{db.Collection("user")}

	// Root
	mux.GET("/", http.RedirectHandler("/list", 301).ServeHTTP)

	// Static
	mux.GET("/static/*", static)

	// Login
	mux.GET("/login", SetUpTemplateData(auth.LoginGET))
	mux.POST("/login", auth.LoginPOST)

	// Signup
	mux.GET("/signup", SetUpTemplateData(auth.SignupGET))
	mux.POST("/signup", auth.SignupPOST)

	// List
	mux.GET("/list", SetUpAuth(SetUpTemplateData(todo.ListGET)))

	// Add
	mux.GET("/add", SetUpAuth(SetUpTemplateData(todo.AddGET)))
	mux.POST("/add", SetUpAuth(todo.AddPOST))

	// Update
	mux.GET("/update", SetUpAuth(SetUpTemplateData(todo.UpdateGET)))
	mux.POST("/update", SetUpAuth(todo.UpdatePOST))

	// Delete
	mux.GET("/delete", SetUpAuth(SetUpTemplateData(todo.DeleteGET)))
	mux.POST("/delete", SetUpAuth(todo.DeletePOST))

	return *mux
}
