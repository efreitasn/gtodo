package handlers

import (
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/efreitasn/go-todo/cmd/go-todo/internal/handlers/auth"
	"github.com/efreitasn/go-todo/cmd/go-todo/internal/handlers/middlewares"
	"github.com/efreitasn/go-todo/cmd/go-todo/internal/handlers/static"
	"github.com/efreitasn/go-todo/cmd/go-todo/internal/handlers/todo"
	"github.com/efreitasn/go-todo/pkg/mids"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewMux creates a mux handler.
func NewMux(db *mongo.Database) http.Handler {
	mux := httptreemux.NewContextMux()
	todo := todo.New(db.Collection("todos"))
	auth := auth.New(db.Collection("user"))

	// Root
	mux.GET("/", http.RedirectHandler("/list", 301).ServeHTTP)

	// Static
	mux.GET("/static/*", static.Static)

	// Login
	mux.GET("/login", mids.New(auth.LoginGET)(
		middlewares.SetUpTemplateData,
	))
	mux.POST("/login", auth.LoginPOST)

	// Signup
	mux.GET("/signup", mids.New(auth.SignupGET)(
		middlewares.SetUpTemplateData,
	))
	mux.POST("/signup", auth.SignupPOST)

	// List
	mux.GET("/list", mids.New(todo.ListGET)(
		middlewares.HasToBeAuth,
		middlewares.SetUpTemplateData,
	))

	// Add
	mux.GET("/add", mids.New(todo.AddGET)(
		middlewares.HasToBeAuth,
		middlewares.SetUpTemplateData,
	))
	mux.POST("/add", mids.New(todo.AddPOST)(
		middlewares.HasToBeAuth,
	))

	// Update
	mux.GET("/update", mids.New(todo.UpdateGET)(
		middlewares.HasToBeAuth,
		middlewares.SetUpTemplateData,
	))
	mux.POST("/update", mids.New(todo.UpdatePOST)(
		middlewares.HasToBeAuth,
	))

	// Delete
	mux.GET("/delete", mids.New(todo.DeleteGET)(
		middlewares.HasToBeAuth,
		middlewares.SetUpTemplateData,
	))
	mux.POST("/delete", mids.New(todo.DeletePOST)(
		middlewares.HasToBeAuth,
	))

	return *mux
}
