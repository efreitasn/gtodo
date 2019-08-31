package handlers

import (
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/efreitasn/gtodo/cmd/gtodo/internal/handlers/about"
	"github.com/efreitasn/gtodo/cmd/gtodo/internal/handlers/auth"
	authMiddlewares "github.com/efreitasn/gtodo/cmd/gtodo/internal/handlers/middlewares/auth"
	templateMiddlewares "github.com/efreitasn/gtodo/cmd/gtodo/internal/handlers/middlewares/template"
	"github.com/efreitasn/gtodo/cmd/gtodo/internal/handlers/static"
	"github.com/efreitasn/gtodo/cmd/gtodo/internal/handlers/todo"
	"github.com/efreitasn/gtodo/pkg/mids"
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

	// About
	mux.GET("/about", mids.New(about.About)(
		authMids.PerformAuth,
		templateMids.SetUpTemplateData,
	))

	// Login
	mux.GET("/login", mids.New(auth.LoginGET)(
		authMids.PerformAuth,
		authMids.HasToBeUnauth,
		templateMids.SetUpTemplateData,
		templateMids.PushAssets,
	))
	mux.POST("/login", mids.New(auth.LoginPOST)(
		authMids.PerformAuth,
		authMids.HasToBeUnauth,
	))

	// Signup
	mux.GET("/signup", mids.New(auth.SignupGET)(
		authMids.PerformAuth,
		authMids.HasToBeUnauth,
		templateMids.SetUpTemplateData,
		templateMids.PushAssets,
	))
	mux.POST("/signup", mids.New(auth.SignupPOST)(
		authMids.PerformAuth,
		authMids.HasToBeUnauth,
	))

	// List
	mux.GET("/list", mids.New(todo.ListGET)(
		authMids.PerformAuth,
		authMids.HasToBeAuth,
		templateMids.SetUpTemplateData,
		templateMids.PushAssets,
	))

	// Add
	mux.GET("/add", mids.New(todo.AddGET)(
		authMids.PerformAuth,
		authMids.HasToBeAuth,
		templateMids.SetUpTemplateData,
		templateMids.PushAssets,
	))
	mux.POST("/add", mids.New(todo.AddPOST)(
		authMids.PerformAuth,
		authMids.HasToBeAuth,
	))

	// Update
	mux.GET("/update", mids.New(todo.UpdateGET)(
		authMids.PerformAuth,
		authMids.HasToBeAuth,
		templateMids.SetUpTemplateData,
		templateMids.PushAssets,
	))
	mux.POST("/update", mids.New(todo.UpdatePOST)(
		authMids.PerformAuth,
		authMids.HasToBeAuth,
	))

	// Delete
	mux.GET("/delete", mids.New(todo.DeleteGET)(
		authMids.PerformAuth,
		authMids.HasToBeAuth,
		templateMids.SetUpTemplateData,
		templateMids.PushAssets,
	))
	mux.POST("/delete", mids.New(todo.DeletePOST)(
		authMids.PerformAuth,
		authMids.HasToBeAuth,
	))

	return *mux
}
