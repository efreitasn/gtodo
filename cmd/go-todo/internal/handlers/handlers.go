package handlers

import (
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/efreitasn/go-todo/cmd/go-todo/internal/handlers/auth"
	authMiddlewares "github.com/efreitasn/go-todo/cmd/go-todo/internal/handlers/middlewares/auth"
	templateMiddlewares "github.com/efreitasn/go-todo/cmd/go-todo/internal/handlers/middlewares/template"
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
	authMids := authMiddlewares.New(db.Collection("user"))
	templateMids := templateMiddlewares.New()

	// Root
	mux.GET("/", http.RedirectHandler("/list", 301).ServeHTTP)

	// Static
	mux.GET("/static/*", static.Static)

	// Login
	mux.GET("/login", mids.New(auth.LoginGET)(
		authMids.HasToBeUnauth,
		templateMids.SetUpTemplateData,
	))
	mux.POST("/login", mids.New(auth.LoginPOST)(
		authMids.HasToBeUnauth,
	))

	// Signup
	mux.GET("/signup", mids.New(auth.SignupGET)(
		authMids.HasToBeUnauth,
		templateMids.SetUpTemplateData,
	))
	mux.POST("/signup", mids.New(auth.SignupPOST)(
		authMids.HasToBeUnauth,
	))

	// List
	mux.GET("/list", mids.New(todo.ListGET)(
		authMids.HasToBeAuth,
		templateMids.SetUpTemplateData,
	))

	// Add
	mux.GET("/add", mids.New(todo.AddGET)(
		authMids.HasToBeAuth,
		templateMids.SetUpTemplateData,
	))
	mux.POST("/add", mids.New(todo.AddPOST)(
		authMids.HasToBeAuth,
	))

	// Update
	mux.GET("/update", mids.New(todo.UpdateGET)(
		authMids.HasToBeAuth,
		templateMids.SetUpTemplateData,
	))
	mux.POST("/update", mids.New(todo.UpdatePOST)(
		authMids.HasToBeAuth,
	))

	// Delete
	mux.GET("/delete", mids.New(todo.DeleteGET)(
		authMids.HasToBeAuth,
		templateMids.SetUpTemplateData,
	))
	mux.POST("/delete", mids.New(todo.DeletePOST)(
		authMids.HasToBeAuth,
	))

	return *mux
}
